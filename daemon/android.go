package main

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"strings"
)

import (
	"github.com/go-pg/pg"
	"github.com/Vonng/go-android-search/wdj"
	"github.com/Vonng/go-android-search/sjqq"
	log "github.com/Sirupsen/logrus"
)

// ID type indicator
const (
	TypePackage  = '!'
	TypeKeywords = '#'
)

// Message hold msg type with one [optional] leading byte and following ID value.
// If no leading letter of `!@#` is provided, Bundle ID is used as default.
type Message struct {
	Type byte
	ID   string
}

// NewMessage will build message from raw string
func NewMessage(msg string) (m Message) {
	if len(msg) < 2 {
		return
	}

	m.Type, m.ID = msg[0], string(msg[1:])
	if m.Type == TypePackage || m.Type == TypeKeywords {
		return
	} else {
		m.Type = TypePackage
		m.ID = msg
	}
	return
}

// Message_Valid tells if this message is valid
func (m *Message) Valid() bool {
	return (m.Type == TypePackage || m.Type == TypeKeywords) && len(m.ID) > 0
}

// Global postgreSQL instance
var Pg = pg.Connect(&pg.Options{
	Addr:     ":5432",
	Database: "meta",
	User:     "meta",
	Password: "meta",
})

// SeenID will check whether given iTunesID is already in database
func SeenID(apk string) bool {
	var res int64
	_, err := Pg.Query(&res, `SELECT count(id) FROM wdj WHERE id =`+apk)
	if err == nil && res == 1 {
		return true
	}
	return false
}

// HandleWdj will fetch and save android application info from wandoujia by package name
func HandleWdj(apk string) error {
	if android, err := wdj.Parse(apk); err != nil {
		return err
	} else {
		return android.Save(Pg)
	}
}

// HandleSjqq will fetch and save android application info from yingyongbao by package name
func HandleSjqq(apk string) error {
	if android, err := sjqq.Parse(apk); err != nil {
		return err
	} else {
		return android.Save(Pg)
	}
}

// HandleApplesByKeyword find a series of app returned by iTunes Search API
// and put them into queue
func HandleKeyword(keyword string) error {
	apks, err := wdj.Search(keyword)
	if err != nil {
		return err
	}

	if len(apks) == 0 {
		return nil
	}

	var sql bytes.Buffer
	cnt := 0
	sql.WriteString("INSERT INTO android_queue(id) VALUES ")
	for i, apk := range apks {
		if SeenID(apk) {
			continue
		}
		cnt += 1
		if i > 0 {
			sql.WriteByte(',')
		}
		sql.WriteString(`('!`)
		sql.WriteString(apk)
		sql.WriteString(`')`)
	}
	if cnt == 0 {
		return nil
	}
	sql.WriteString("ON CONFLICT DO NOTHING;")
	res, err := Pg.Exec(sql.String())
	if err != nil {
		return err
	}
	log.Infof("[SEARCH] keyword %s found %d, add %d", keyword, len(apks), res.RowsAffected())
	return nil
}

// Producer will pull task from PostgreSQL table `android_queue`
func Producer() <-chan Message {
	log.Info("[PROD] initializing...")
	stmt, err := Pg.Prepare(`DELETE FROM android_queue WHERE id IN (SELECT id FROM android_queue LIMIT 100) RETURNING id;`)
	if err != nil {
		log.Info("[PROD] prepare job statement failed...check postgres instance")
		return nil
	}
	c := make(chan Message)
	go func(chan<- Message) {
		sleep := time.Second
		for {
			var ids []string
			_, err := stmt.Query(&ids)
			if len(ids) == 0 {
				log.Infof("[PROD] empty queue. sleep %d s", sleep/1e9)
				time.Sleep(sleep)
				if sleep < 30*time.Second {
					sleep *= 2
				}
				continue
			} else {
				// reset sleep counter to 1s
				sleep = time.Second
			}

			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			for _, id := range ids {
				if msg := NewMessage(id); msg.Valid() {
					c <- msg
				}
			}
		}
	}(c)
	log.Info("[PROD] init complete")
	return c
}

// Worker will handle incoming task
func Worker(id int, c <-chan Message) {
	log.Infof("[WORKER:%d] init", id)
	var err error
	for msg := range c {
		switch msg.Type {
		case TypePackage:
			log.Infof("[WORKER:%d] handle Package=%s @ wdj", id, msg.ID)
			if err = HandleWdj(msg.ID); err != nil {
				log.Errorf("[WORKER:%d] handle Package=%s @ wdj failed: %s", id, msg.ID, err.Error())
			}
			log.Infof("[WORKER:%d] done Package=%s @ wdj", id, msg.ID)
			log.Infof("[WORKER:%d] handle Package=%s @ sjqq", id, msg.ID)
			if err = HandleSjqq(msg.ID); err != nil {
				log.Errorf("[WORKER:%d] handle Package=%s @ sjqq failed: %s", id, msg.ID, err.Error())
			}
			log.Infof("[WORKER:%d] done Package=%s @ sjqq", id, msg.ID)
		case TypeKeywords:
			log.Infof("[WORKER:%d] handle Keyword=%s", id, msg.ID)
			if err = HandleKeyword(msg.ID); err != nil {
				log.Errorf("[WORKER:%d] handle Keyword=%s failed: %s", id, msg.ID, err.Error())
			}
			log.Infof("[WORKER:%d] done keyword=%s", id, msg.ID)
		}
	}
	log.Infof("[WORK] %d finish", id)
}

// Run will start n worker and one producer.
func Run(n int) {
	log.Infof("[RUN] init with %d worker...", n)
	c := Producer()
	for i := 1; i <= n; i++ {
		go Worker(i, c)
	}
}

func main() {
	log.SetLevel(log.InfoLevel)

	if len(os.Args) > 2 {
		action, id := os.Args[1], os.Args[2]
		action = strings.ToLower(action)
		var err error
		switch action {

		case "id", "pkg", "package", "apk":
			log.Infof("handle Package=%s @ wdj", id)
			if err = HandleWdj(id); err != nil {
				log.Errorf("handle Package=%s @ wdj failed: %s", id, err.Error())
			}
			log.Infof("done Package=%s @ wdj", id)
			log.Infof("handle Package=%s @ sjqq", id)
			if err = HandleSjqq(id); err != nil {
				log.Errorf("handle Package=%s @ sjqq failed: %s", id, err.Error())
			}
			log.Infof("done Package=%s @ sjqq", id)
		case "k", "key", "keyword", "keywords", "search":
			log.Infof("handle Keywords=%s", id)
			if err := HandleKeyword(id); err != nil {
				log.Errorf("handle Keywords=%s failed: %s", id, err.Error())
			}
			log.Infof("done Keywords=%s", id)
		}
		os.Exit(0)
	}

	Run(5)
	<-make(chan bool)
}
