package block

import (
	"bytes"
	"encoding/gob"
	log "github.com/corgi-kx/logcustom"
	"math/big"
	"time"
)

type Block struct {
	//上一个区块的hash
	PreHash []byte
	//数据data
	Transactions []Transaction
	//时间戳
	TimeStamp int64
	//区块高度
	Height int
	//随机数
	Nonce int64
	//本区块hash
	Hash []byte
}

//进行挖矿来生成区块
func mineBlock(transaction []Transaction, preHash []byte, height int) (*Block, error) {
	timeStamp := time.Now().Unix()
	//hash数据+时间戳+上一个区块hash
	block := Block{preHash, transaction, timeStamp, height, 0, nil}
	pow := NewProofOfWork(&block)
	nonce, hash, err := pow.run()
	if err != nil {
		return nil, err
	}
	block.Nonce = nonce
	block.Hash = hash[:]
	log.Info("pow verify : ", pow.Verify())
	log.Infof("已生成新的区块,区块高度为%d", block.Height)
	return &block, nil
}

//生成创世区块
func newGenesisBlock(transaction []Transaction) *Block {
	//创世区块的上一个块hash默认设置成下面的样子
	preHash := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	//生成创世区块
	genesisBlock, err := mineBlock(transaction, preHash, 1)
	if err != nil {
		log.Error(err)
	}
	return genesisBlock
}

// 将Block对象序列化成[]byte
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func (v *Block) Deserialize(d []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(v)
	if err != nil {
		log.Panic(err)
	}
}

func isGenesisBlock(block *Block) bool {
	var hashInt big.Int
	hashInt.SetBytes(block.PreHash)
	if big.NewInt(0).Cmp(&hashInt) == 0 {
		return true
	}
	return false
}
