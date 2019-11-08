package block

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"github.com/corgi-kx/blockchain_golang/util"
	log "github.com/corgi-kx/logcustom"
	"math/big"
)

//工作量证明(pow)结构体
type proofOfWork struct {
	*Block
	Target *big.Int
}

//获取POW实例
func NewProofOfWork(block *Block) *proofOfWork {
	target := big.NewInt(1)
	//返回1 << 256-TargetBits 的一个大数
	target.Lsh(target, 256-TargetBits)
	pow := &proofOfWork{block, target}
	return pow
}

//进行hash运算,获取到当前区块的hash值
func (p *proofOfWork) run() (int,[]byte,error) {
	nonce := 0
	var hashByte [32]byte
	var hashInt big.Int
	log.Info("正在挖矿...")
	for nonce < maxInt {
		//检测网络上其他节点是否已经挖出了区块
		if p.Height <= NewestBlockHeight {
			return 0,nil,errors.New("检测到当前节点已接收到最新区块，所以终止此块的挖矿操作")
		}
		data := p.jointData(nonce)
		hashByte = sha256.Sum256(data)
		//fmt.Printf("\r current hash : %x", hashByte)
		//将hash值转换为大数字
		hashInt.SetBytes(hashByte[:])
		//如果hash后的data值小于设置的挖矿难度大数字,则代表挖矿成功!
		if hashInt.Cmp(p.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	log.Infof("挖到区块了,区块hash为: %x", hashByte)
	return nonce, hashByte[:],nil
}

//检验区块是否有效
func (p *proofOfWork) Verify() bool {
	target := big.NewInt(1)
	target.Lsh(target, 256-TargetBits)
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
	targetBitsByte := util.Int64ToBytes(int64(TargetBits))
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
