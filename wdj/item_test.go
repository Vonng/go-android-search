package wdj

import (
	"testing"
	"github.com/go-pg/pg"
	"fmt"
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
	//Pg.Exec(`TRUNCATE wdj;`)

	for _, id := range []string{
		"com.autonavi.minimap",
		"com.tencent.tmgp.sgame",
		"com.tencent.mm",
	} {
		app, err := Parse(id)
		if err != nil {
			t.Error(err)
		}
		app.Print()
		err = app.Save(Pg)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestSearch(t *testing.T) {
	keywords := []string{
		"王者荣耀",
		"守望先锋",
	}
	for _, kw := range keywords {
		if result, err := Search(kw); err != nil {
			t.Error(err)
		} else {
			fmt.Println(result)
		}
	}
}
