package block

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"math/big"
	log "myCode/public_blockchain/part7-network/logcustom"
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
	Nonce int
	//本区块hash
	Hash []byte
}

func newBlock(transaction []Transaction, preHash []byte, height int) *Block {
	timeStamp := time.Now().Unix()
	//hash数据+时间戳+上一个区块hash
	block := Block{preHash, transaction, timeStamp, height, 0, nil}
	pow := NewProofOfWork(&block)
	nonce, hash := pow.run()
	block.Nonce = nonce
	block.Hash = hash[:]
	fmt.Println("pow verify : ", pow.Verify())
	fmt.Println("已生成新的区块...")
	return &block
}

func newGenesisBlock(transaction []Transaction) *Block {
	preHash := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	genesisBlock := newBlock(transaction, preHash, 1)
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

func isGenesisBlock(block *Block) bool {
	var hashInt big.Int
	hashInt.SetBytes(block.PreHash)
	if big.NewInt(0).Cmp(&hashInt) == 0 {
		return true
	}
	return false
}

func Deserialize(d []byte, i interface{}) interface{} {
	var model interface{}
	switch i.(type) {
	case *Block:
		model = i.(*Block)
	case *addressList:
		model = i.(*addressList)
	case *bitcoinKeys:
		gob.Register(elliptic.P256())
		model = i.(*bitcoinKeys)
	default:
		log.Fatal("Deserialize err :没有可反序列化的类型")
	}
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(model)
	if err != nil {
		log.Panic(err)
	}
	return model
}
