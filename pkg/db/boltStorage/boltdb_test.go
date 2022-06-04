package boltStorage

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func ReadView(id int, db *BoltdbStorage, key []byte, bucketName []byte) {
	err := db.db.View(func(tx *bolt.Tx) error {
		time.Sleep(time.Second * 2)
		b := tx.Bucket(bucketName)
		buffer := b.Get(key)
		log.Println(time.Now(), "Read: ", id, buffer)
		return nil
	})
	wg.Done()
	if err != nil {
		log.Println(err)
	}
}

func WriteView(id int, db *BoltdbStorage, key []byte, value []byte, bucketName []byte) {
	err := db.db.Update(func(tx *bolt.Tx) error {
		time.Sleep(time.Second * 2)
		b := tx.Bucket(bucketName)
		err := b.Put(key, value)
		log.Println(time.Now(), "Write: ", id)
		return err
	})
	wg.Done()
	if err != nil {
		log.Println(err)
	}
}

func TestBoltdb(t *testing.T) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := Open("my.db")
	if err != nil {
		log.Println("open error!")
		return
	}

	bucketName := []byte("myBucket")
	bucketName2 := []byte("myBucket2")
	key := []byte("1")
	value := []byte("你好世界")

	err = db.CreateBucketIfNotExist(bucketName)
	if err != nil {
		log.Println(fmt.Errorf("create Bucket Error: %s", err))
	}

	err = db.Set(key, value, bucketName)
	if err != nil {
		log.Println(fmt.Errorf("set Error: %s", err))
	}

	value = []byte("你好")
	err = db.Set(key, value, bucketName)
	if err != nil {
		log.Println(fmt.Errorf("set Error: %s", err))
	}

	buffer, ok := db.Get(key, bucketName)
	log.Println(string(buffer), ok)

	err = db.Delete(key, bucketName)
	if err != nil {
		log.Println(fmt.Errorf("delete Error: %s", err))
	}

	buffer, ok = db.Get(key, bucketName)
	log.Println(buffer, ok)

	wg.Add(20)

	for i := 0; i < 10; i++ {
		go WriteView(i, db, []byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)), bucketName)
		go WriteView(i, db, []byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)), bucketName2)
	}

	for i := 0; i < 10; i++ {
		go ReadView(i, db, []byte(strconv.Itoa(i)), bucketName)
		go ReadView(i, db, []byte(strconv.Itoa(i)), bucketName2)
	}

	wg.Wait()

	err = db.DeleteBucket(bucketName)
	if err != nil {
		log.Println(fmt.Errorf("delete bucket Error: %s", err))
	}
	err = db.DeleteBucket(bucketName2)
	if err != nil {
		log.Println(fmt.Errorf("delete bucket Error: %s", err))
	}

	err = db.Close()
	if err != nil {
		log.Println(fmt.Errorf("close Error: %s", err))
	}
}
