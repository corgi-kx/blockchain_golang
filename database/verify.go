package database

import (
	"github.com/boltdb/bolt"
	log "myCode/public_blockchain/part7-network/logcustom"
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
