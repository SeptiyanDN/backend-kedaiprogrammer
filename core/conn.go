package core

import (
	"fmt"
	"kedaiprogrammer/kedaihelpers"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConnect() kedaihelpers.DBStruct {
	username := viper.GetString("database_staging.username")
	password := viper.GetString("database_staging.password")
	database := viper.GetString("database_staging.name")
	host := viper.GetString("database_staging.host")
	port := viper.GetInt("database_staging.port")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error Connecting DB => ", err)
		os.Exit(0)
	}

	maxLifetime, _ := time.ParseDuration(viper.GetString("database.max_lifetime_connection") + "s")
	db.SetMaxIdleConns(viper.GetInt("database.max_idle_dbection"))
	db.SetConnMaxLifetime(maxLifetime)
	dbs := kedaihelpers.DBStruct{Dbx: db}

	return dbs
}

func InitGorm() (*gorm.DB, error) {
	username := viper.GetString("database_staging.username")
	password := viper.GetString("database_staging.password")
	database := viper.GetString("database_staging.name")
	host := viper.GetString("database_staging.host")
	port := viper.GetInt("database_staging.port")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Tokyo", host, port, username, password, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("🚀 Connected Successfully to the Database")
	// db.AutoMigrate(users.User{})
	// db.AutoMigrate(businesses.Business{})
	// db.AutoMigrate(categories.Category{})
	// db.AutoMigrate(services.Service{})

	return db, nil
}
