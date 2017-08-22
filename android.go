package go_android_search

import "github.com/go-pg/pg"

type Android interface {
	Print()
	Save(db *pg.DB)
}
