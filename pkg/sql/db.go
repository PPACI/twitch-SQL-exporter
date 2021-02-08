package sql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbType int

type DbOpts struct {
	Type string
	Dsn  string
	Opts *gorm.Config
}

func NewDb(opts *DbOpts) (*gorm.DB, error) {
	var d gorm.Dialector
	switch opts.Type {
	case "postgres":
		d = postgres.Open(opts.Dsn)
	case "sqlite":
		d = sqlite.Open(opts.Dsn)
	default:
		return nil, fmt.Errorf("unsupported db type: %s", opts.Type)

	}
	db, err := gorm.Open(d, opts.Opts)
	return db, err
}
