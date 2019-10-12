package database

import (
	"github.com/boltdb/bolt"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
	"os"
)

func IsBlotExist(nodeID string) bool {
	var DBFileName = "blockchain_" + nodeID + ".db"
	_, err := os.Stat(DBFileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsBucketExist(bd *BlockchainDB, bt BucketType) bool {
	var isBucketExist bool
	err := bd.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			isBucketExist = false
		} else {
			isBucketExist = true
		}
		return nil
	})
	if err != nil {
		log.Panic("datebase IsBucketExist err:", err)
	}
	return isBucketExist
}
