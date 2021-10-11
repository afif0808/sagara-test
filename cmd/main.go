package main

import (
	"fmt"
	"log"
	"os"

	"github.com/afif0808/sagara-test/internal/domain"
	authmodule "github.com/afif0808/sagara-test/internal/modules/auth"
	filestoragemodule "github.com/afif0808/sagara-test/internal/modules/filestorage"
	productModule "github.com/afif0808/sagara-test/internal/modules/product"
	usermodule "github.com/afif0808/sagara-test/internal/modules/user"
	"github.com/labstack/echo"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func migrateDatabase() {
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	db := os.Getenv("MYSQL_DATABASE")
	port := os.Getenv("MYSQL_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, db)
	conn, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Panic(err)
	}

	err = conn.AutoMigrate(&domain.Product{}, &domain.User{})
	if err != nil {
		log.Panic(err)
	}

}

func initMySQLDB() (readDB, writeDB *sqlx.DB) {
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	db := os.Getenv("MYSQL_DATABASE")
	port := os.Getenv("MYSQL_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, db)
	readDB, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Panic(err)
	}
	writeDB, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Panic(err)
	}

	return
}

func main() {
	godotenv.Load(".env")
	readDB, writeDB := initMySQLDB()
	defer func() {
		readDB.Close()
		writeDB.Close()
	}()
	migrateDatabase()

	e := echo.New()

	usermodule.InjectUserModule(e, readDB, writeDB)
	authmodule.InjectAuthModule(e, readDB, writeDB)
	productModule.InjectProductModule(e, readDB, writeDB)
	filestoragemodule.InjectFileStorageModule(e, readDB, writeDB)
	e.Start(":" + os.Getenv("APP_PORT"))
}
