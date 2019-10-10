package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var db *gorm.DB

func init() {
	rand.Seed(time.Now().Unix())
}
func CreateMysqlDb(verbose bool) error {
	user := viper.GetString("db.user")
	host := viper.GetString("db.host")
	password := viper.GetString("db.password")
	port := viper.GetInt("db.port")
	dbName := viper.GetString("db.dbName")

	conn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbName)
	var databaseInstance *gorm.DB
	var err error

	databaseInstance, err = gorm.Open("mysql", conn)
	if err != nil {
		return fmt.Errorf("init MySQL db failed in %s, %s", conn, err)
	}

	db = databaseInstance
	db.LogMode(verbose)

	return RunMigrate()
}

func CreateSqliteDb(verbose bool) error {
	var databaseInstance *gorm.DB
	var err error
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("can not find executable dir %s", err)
	}
	d := filepath.Dir(exePath)
	dp := filepath.Join(d, "db.sqlite3")
	log.Println("use SQLite database in path ", dp)
	databaseInstance, err = gorm.Open("sqlite3", dp)
	if err != nil {
		return fmt.Errorf("init SQLite3 db failed in %s, %s", dp, err)
	}
	db = databaseInstance
	db.LogMode(verbose)
	return RunMigrate()
}

func Close() {
	if db != nil {
		db.Close()
	}
}

type BaseModel struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
