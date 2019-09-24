package block

import "myCode/public_blockchain/part7-network/database"

type blockchainIterator struct {
	CurrentBlockHash []byte
	BD               *database.BlockchainDB
}

func (bi *blockchainIterator) next() *block {
	currentByte := bi.BD.View(bi.CurrentBlockHash, database.BlockBucket)
	b := deserialize(currentByte, &block{})
	block := b.(*block)
	bi.CurrentBlockHash = block.PreHash
	return block
}
