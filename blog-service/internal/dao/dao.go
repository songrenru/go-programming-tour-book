package dao

import (
	"github.com/jinzhu/gorm"
)

type Dao struct {
	engine *gorm.DB
}

func New(engine *gorm.DB) *Dao {
	return &Dao{engine: engine}
}

func (d *Dao) Begin() *Dao {
	d.engine = d.engine.Begin()
	return d
}

func (d *Dao) Commit() *Dao {
	d.engine = d.engine.Commit()
	return d
}

func (d *Dao) Rollback() *Dao {
	d.engine = d.engine.Rollback()
	return d
}
