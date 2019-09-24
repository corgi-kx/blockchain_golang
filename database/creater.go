package database

import (
	"github.com/boltdb/bolt"
	"log"
)

// 仓库类型
type BucketType string

const (
	BlockBucket BucketType = "blocks"
	AddrBucket  BucketType = "address"
	UTXOBucket  BucketType = "utxo"
)

type BlockchainDB struct {
	DB     *bolt.DB
	nodeID string
}

func New(nodeID string) *BlockchainDB {
	var DBFileName = "blockchain_" + nodeID + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	bd := &BlockchainDB{db, nodeID}
	return bd
}

func (bd *BlockchainDB) Close() {
	err := bd.DB.Close()
	if err != nil {
		log.Panic("db close err :", err)
	}
}
