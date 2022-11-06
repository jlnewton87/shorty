package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

var bucket = []byte("shorty")

type Store struct {
}

func (s *Store) init() {
	db, err := bolt.Open("./store.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func (s *Store) addShorty(key string, shortyType string, shortyTarget string) {
	value := shortyType + shortyTarget
	db, err := bolt.Open("./store.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Store) getShorty(key string) shortyReq {
	var output shortyReq
	db, err := bolt.Open("./store.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}
		val := bucket.Get([]byte(key))
		if len(val) > 0 {
			output = shortyReq{sType: shortyType(val[0:1]), target: string(val[1:])}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return output
}
