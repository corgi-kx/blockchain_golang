package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"myCode/public_blockchain/part7-network/util"
)

type proofOfWork struct {
	*Block
	Target *big.Int
}

func NewProofOfWork(block *Block) *proofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256-targetBits)
	pow := &proofOfWork{block, target}
	return pow
}

func (p *proofOfWork) run() (int, []byte) {
	nonce := 0
	var hashByte [32]byte
	var hashInt big.Int
	fmt.Println("Mining the block  ....")
	for nonce < maxInt {
		data := p.jointData(nonce)
		hashByte = sha256.Sum256(data)
		fmt.Printf("\r current hash : %x", hashByte)
		//将hash值转换为大数字
		hashInt.SetBytes(hashByte[:])
		if hashInt.Cmp(p.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println("")
	return nonce, hashByte[:]
}

func (p *proofOfWork) Verify() bool {
	target := big.NewInt(1)
	target.Lsh(target, 256-targetBits)
	data := p.jointData(p.Block.Nonce)
	hash := sha256.Sum256(data)
	var hashInt big.Int
	hashInt.SetBytes(hash[:])
	if hashInt.Cmp(target) == -1 {
		return true
	}
	return false
}

//将上一区块hash、数据、时间戳、难度位数、随机数 拼接成字节数组
func (p *proofOfWork) jointData(nonce int) (data []byte) {
	preHash := p.Block.PreHash
	timeStampByte := util.Int64ToBytes(p.Block.TimeStamp)
	heightByte := util.Int64ToBytes(int64(p.Block.Height))
	nonceByte := util.Int64ToBytes(int64(nonce))
	targetBitsByte := util.Int64ToBytes(int64(targetBits))
	//拼接成交易数组
	transData := [][]byte{}
	for _,v := range p.Block.Transactions {
		tBytes := v.getTransBytes()   //这里为什么要用到自己写的方法，而不是gob序列化，是因为gob同样的数据序列化后的字节数组有可能不一致，无法用于hash验证
		transData = append(transData, tBytes)
	}
	//获取交易数据的根默克尔节点
	mt := util.NewMerkelTree(transData)

	data = bytes.Join([][]byte{
		preHash,
		timeStampByte,
		heightByte,
		mt.MerkelRootNode.Data,
		nonceByte,
		targetBitsByte},
		[]byte(""))
	return
}
