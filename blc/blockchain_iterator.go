package block

import "github.com/corgi-kx/blockchain_golang/database"

type blockchainIterator struct {
	CurrentBlockHash []byte
	BD               *database.BlockchainDB
}
func  NewBlockchainIterator(bc *blockchain) *blockchainIterator {
	blockchainIterator := &blockchainIterator{bc.BD.View([]byte(LastBlockHashMapping), database.BlockBucket), bc.BD}
	return blockchainIterator
}

func (bi *blockchainIterator) Next() *Block {
	currentByte := bi.BD.View(bi.CurrentBlockHash, database.BlockBucket)
	if len(currentByte) == 0 {
		return nil
	}
	block:=Block{}
	block.Deserialize(currentByte)
	bi.CurrentBlockHash = block.PreHash
	return &block
}
