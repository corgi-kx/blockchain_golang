package block

import "github.com/corgi-kx/blockchain_golang/database"

type blockchainIterator struct {
	CurrentBlockHash []byte
	BD               *database.BlockchainDB
}
func  NewBlockchainIterator(bc *blockchain) *blockchainIterator {
	blockchainIterator := &blockchainIterator{bc.BD.View([]byte(lastBlockHashMapping), database.BlockBucket), bc.BD}
	return blockchainIterator
}

func (bi *blockchainIterator) Next() *Block {
	currentByte := bi.BD.View(bi.CurrentBlockHash, database.BlockBucket)
	if currentByte == nil {
		return nil
	}
	b := Deserialize(currentByte, &Block{})
	block := b.(*Block)
	bi.CurrentBlockHash = block.PreHash
	return block
}
