package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jasonlvhit/gocron"
)

var ROOT string

var USERNAME string

func save(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)

	}
	defer file.Close()
	png.Encode(file, img)

}

func Signup(username string) {
	response, err := http.Post("http://"+ROOT+"/api/users?id="+username, "text", nil)
	if err != nil {
		log.Println(err)
	} else {
		defer response.Body.Close()
		_, err := io.Copy(os.Stdout, response.Body)
		if err != nil {
			log.Println(err)
		}
	}
}

func Login(username string) {
	response, err := http.Post("http://"+ROOT+"/api/timelog?id="+username, "text", nil)
	if err != nil {
		log.Println(err)
	} else {
		defer response.Body.Close()
		_, err := io.Copy(os.Stdout, response.Body)
		if err != nil {
			log.Println(err)
		}
	}
}

func init() {
	USERNAME = "tony"

}

func main() {
	db, err := bolt.Open("monitnorxx.bolt", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Start a writable transaction.

	tx, err := db.Begin(true)
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback()

	// Use the transaction...
	bucket, err := tx.CreateBucketIfNotExists([]byte("monitor"))
	if err != nil {
		log.Println(err)
	}

	usernameB := bucket.Get([]byte("username"))
	hostB := bucket.Get([]byte("host"))
	lastLoginB := bucket.Get([]byte("lastLogin"))

	// Commit the transaction and check for error.
	if err = tx.Commit(); err != nil {
		log.Println(err)
	}

	log.Println(string(usernameB) == "")
	log.Println(string(hostB))

	if string(hostB) == "" {
		fmt.Print("Host IP Address: ")
		fmt.Scanln(&ROOT)
		fmt.Println(ROOT)
		//log.Println(&ROOT)
	} else {
		ROOT = string(hostB)
	}

	if string(usernameB) == "" {
		fmt.Print("Username: ")
		fmt.Scanln(&USERNAME)
		Signup(USERNAME)
	} else {
		USERNAME = string(usernameB)
	}

	log.Printf("username: %s ,  host: %s", USERNAME, ROOT)

	tx, err = db.Begin(true)
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback()

	// Use the transaction...
	bucket = tx.Bucket([]byte("monitor"))
	if err != nil {
		log.Println(err)
	}

	err = bucket.Put([]byte("username"), []byte(USERNAME))
	if err != nil {
		log.Println(err)
	}
	err = bucket.Put([]byte("host"), []byte(ROOT))
	if err != nil {
		log.Println(err)
	}

	// Commit the transaction and check for error.
	if err = tx.Commit(); err != nil {
		log.Println(err)
	}

	lastLogin, err := time.Parse(time.RFC3339, string(lastLoginB))
	if err != nil {
		log.Println(err)
	}
	if lastLogin.YearDay() < time.Now().YearDay() || lastLogin.Year() < time.Now().Year() {
		Login(USERNAME)
		tx, err = db.Begin(true)
		if err != nil {
			log.Println(err)
		}
		defer tx.Rollback()

		// Use the transaction...
		bucket = tx.Bucket([]byte("monitor"))
		if err != nil {
			log.Println(err)
		}

		err = bucket.Put([]byte("lastLogin"), []byte(time.Now().Format(time.RFC3339)))
		if err != nil {
			log.Println(err)
		}
		// Commit the transaction and check for error.
		if err := tx.Commit(); err != nil {
			log.Println(err)
		}

	}

	gocron.Every(20).Seconds().Do(takeShotAndUpload, USERNAME)
	<-gocron.Start()

}
