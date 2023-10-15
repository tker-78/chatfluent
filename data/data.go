package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DbConnection *sql.DB

func init() {
	var err error
	DbConnection, err = sql.Open("postgres", "dbname=chatfluent sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
}

type UUID [16]byte

func createUUID() (uuid string) {
	u := new(UUID)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0x7F) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// SHA-1による暗号化処理
func Encrypt(plaintext string) string {
	cryptext := fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
