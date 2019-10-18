package database

import (
	"errors"
	"github.com/boltdb/bolt"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
)

func (bd *BlockchainDB) Put(k, v []byte, bt BucketType) {
	var DBFileName = "blockchain_" + bd.nodeID + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
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
	var DBFileName = "blockchain_" + bd.nodeID + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	result := []byte{}
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			msg := "datebase view warnning:没有找到仓库：" + string(bt)
			return errors.New(msg)
		}
		result = bucket.Get(k)
		return nil
	})
	if err != nil {
		log.Warn(err)
		return nil
	}
	//不再次赋值的话，返回值会报错，不知道狗日的啥意思
	realResult:=make([]byte,len(result))
	copy(realResult,result)
	return realResult
}

func (bd *BlockchainDB) Delete(k []byte, bt BucketType) bool {
	var DBFileName = "blockchain_" + bd.nodeID + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
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
	var DBFileName = "blockchain_" + bd.nodeID + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
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
