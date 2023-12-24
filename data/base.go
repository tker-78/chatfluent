package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DbConnection *sql.DB

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	// production environment
	if os.Getenv("environment") == "production" {
		_, err = pgxv4.RegisterDriver("cloudsql-postgres", cloudsqlconn.WithIAMAuthN())
		if err != nil {
			log.Fatalf("Error on pgxv4.RegisterDriver: %v", err)
		}

		dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", os.Getenv("INSTANCE_CONNECTION_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"))
		DbConnection, err = sql.Open("cloudsql-postgres", dsn)
		if err != nil {
			log.Fatalf("Error on sql.Open: %v", err)
		}
	} else if os.Getenv("environment") == "development" {
		// todo
		DbConnection, err = sql.Open("postgres", "dbname=chatfluent sslmode=disable")
		if err != nil {
			log.Fatalln(err)
		}
		return

	} else {
		log.Fatalln("cannot read environment")
	}
}

func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("cannot generate UUID", err)
	}
	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

func Encrypt(text string) (crypt string) {
	crypt = fmt.Sprintf("%x", sha1.Sum([]byte(text)))
	return
}
