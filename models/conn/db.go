// set up configuration to database mysql
package conn

import (
	"database/sql"
	"fmt"

	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ModelsDB *gorm.DB
var PostgresDB *gorm.DB

func ConnectSQLDb(retries int) *gorm.DB {
	if retries > 1 {
		log.Printf("Retrying connect to DB instance, Attempt %v", strconv.Itoa(retries))

		if retries > 5 {
			log.Printf("Cannot recovery situation retries > 5 attempt")
			os.Exit(1)
		}
	}

	var connString string

	connString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_NAME"))

	sqlDB, err := sql.Open("mysql", connString)

	if err != nil {
		log.Printf("error on creating conn sql database %v", err)
		ConnectSQLDb(retries + 1)
		return nil
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Second * 10)

	// gormDB, err := gorm.Open("mysql", sqlDB)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		log.Println("error on creating gorm conn ", err)
		ConnectSQLDb(retries + 1)
		return nil
	}

	if err != nil {
		log.Printf("error on creating conn database %v", err)

		ConnectSQLDb(retries + 1)
		return nil
	}

	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		})
	gormDB.Session(&gorm.Session{Logger: newLogger})

	log.Println("database conn successfully")

	return gormDB
}
func ConnectionPostgresDb(retries int) *gorm.DB {
	if retries > 1 {
		log.Printf("Retrying connect to DB instance, Attempt %v", strconv.Itoa(retries))

		if retries > 5 {
			log.Printf("Cannot recovery situation retries > 5 attempt")
			os.Exit(1)
		}
	}
	println(os.Getenv("POSTGRES_HOST"), "host")
	//dsn := "host=34.87.101.46 user=data password=dataapergu123#@! dbname=data port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_NAME"),
		os.Getenv("POSTGRES_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		ConnectionPostgresDb(retries + 1)
		return nil
	}

	fmt.Println("Connected to PostgreSQL successfully!")

	return db
}
