package database

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func (bd *BlockchainDB) Put(k, v []byte, bt BucketType) {
	err := bd.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			var err error
			bucket, err = tx.CreateBucket([]byte(bt))
			if err != nil {
				log.Panic(err)
			}
		}
		err := bucket.Put(k, v)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (bd *BlockchainDB) View(k []byte, bt BucketType) []byte {
	result := []byte{}
	err := bd.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			msg := "datebase view warnning:没有找到仓库：" + string(bt)
			return errors.New(msg)
		}
		result = bucket.Get(k)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}

func (bd *BlockchainDB) Delete(k []byte, bt BucketType) bool {
	err := bd.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			msg := "datebase delete warnning:没有找到仓库：" + string(bt)
			return errors.New(msg)
		}
		err := bucket.Delete(k)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return true
}

func (bd *BlockchainDB) DeleteBucket(bt BucketType) bool {
	err := bd.DB.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bt))
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return true
}
