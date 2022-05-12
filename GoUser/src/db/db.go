package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"me.user/src/config"
)

// 註冊資訊
type User struct {
	UserID   int    `json:"-" gorm:"primaryKey;autoIncrement:true;column:id"`
	Name     string `binding:"required"`
	Mail     string `binding:"required"`
	Password string `gorm:"-" binding:"required"`
	Hash     string `json:"-" gorm:"column:password"`
	Token    string `gorm:"-" binding:"required"`
}

type DB struct {
	db *gorm.DB
}

func NewDB(cfg config.Postgres) *DB {
	connect := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	db, err := gorm.Open(postgres.Open(connect), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	database, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	database.SetMaxOpenConns(3)
	database.SetMaxIdleConns(1)

	return &DB{
		db: db,
	}
}

func (d *DB) CheckMail(mail string) bool {
	var count int64
	d.db.Table("user").Where("mail=?", mail).Count(&count)
	return count >= 1
}

func (d *DB) Create(user *User) error {
	return d.db.Create(user).Error
}

func (d *DB) Update(id int) error {
	return d.db.Table("user").Where("id=?", id).
		UpdateColumn("verified", true).Error
}
