package sjqq

import (
	"testing"
	"github.com/go-pg/pg"
)

func TestApp_Parse(t *testing.T) {
	app, err := Parse("com.tencent.mm")
	if err != nil {
		t.Error(err)
	}
	app.Print()
}

func TestApp_Save(t *testing.T) {
	var Pg = pg.Connect(&pg.Options{
		Addr:     ":5432",
		Database: "haha",
		User:     "haha",
		Password: "xixihaha",
	})

	// Dangerous!!!
	Pg.Exec(`TRUNCATE sjqq;`)

	for _, id := range []string{
		"com.autonavi.minimap",
		"com.tencent.tmgp.sgame",
		"com.tencent.mm",
	} {
		app, err := Parse(id)
		if err != nil {
			t.Error(err)
			continue
		}
		app.Print()
		err = app.Save(Pg)
		if err != nil {
			t.Error(err)
		}
	}
}
