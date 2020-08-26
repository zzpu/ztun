package sqlite

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"ztun/model"
)

type DB struct {
	w *gorm.DB
}

func NewDB() (db *DB) {
	dbw := NewDBWrite("config.db")

	db = &DB{
		w: dbw,
	}
	return
}

func (db *DB) DBWrite() *gorm.DB {
	return db.w
}

func (db *DB) Close() {
	db.w.Close()
}

//初始化写的数据库
func NewDBWrite(name string) *gorm.DB {
	db, err := gorm.Open("sqlite3", name)
	if err != nil {
		panic("连接数据库失败")
	}

	// 自动迁移模式
	db.AutoMigrate(&model.Tunnel{})

	return db
}
