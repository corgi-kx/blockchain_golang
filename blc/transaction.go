package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/corgi-kx/blockchain_golang/util"
	"github.com/ethereum/go-ethereum/log"

)

type Transaction struct {
	TxHash []byte
	//UTXO输入
	Vint []txInput
	//UTXO输出
	Vout []txOutput
}

func (t *Transaction) hash() {
	tBytes := t.Serialize()
	//加入随机数byte
	randomNumber := util.GenerateRealRandom()
	randomByte := util.Int64ToBytes(randomNumber)
	sumByte := bytes.Join([][]byte{tBytes, randomByte}, []byte(""))
	hashByte := sha256.Sum256(sumByte)
	t.TxHash = hashByte[:]
}

//作为数字签名的hash方法，为什么不用gob序列化后hash，因为涉及到tcp传输gob直接序列化有问题，所以单独拼接成byte数组再hash
func (t *Transaction) hashSign() []byte {
	t.TxHash = nil
	nHash:=[]byte{}
	for _,v:=range t.Vint {
		nHash=append(nHash,v.TxHash...)
		nHash=append(nHash,v.PublicKey...)
		nHash=append(nHash,util.Int64ToBytes(int64(v.Index))...)
	}
	for _,v:=range t.Vout {
		nHash=append(nHash,v.PublicKeyHash...)
		nHash=append(nHash,util.Int64ToBytes(int64(v.Value))...)
	}
	hashByte := sha256.Sum256(nHash)
	return hashByte[:]
}

// 将transaction序列化成[]byte
func (t *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(t)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}


func (t *Transaction) getTransBytes() []byte {
	if t.TxHash == nil || t.Vint == nil || t.Vout == nil{
		log.Error("交易信息不完整，无法拼接成字节数组")
		return nil
	}
	transBytes:=[]byte{}
	transBytes = append(transBytes,t.TxHash...)
	for _,v := range t.Vint {
		transBytes = append(transBytes,v.TxHash...)
		transBytes = append(transBytes,util.Int64ToBytes(int64(v.Index))...)
		transBytes = append(transBytes,v.Signature...)
		transBytes = append(transBytes,v.PublicKey...)
	}
	for _,v := range t.Vout {
		transBytes = append(transBytes,util.Int64ToBytes(int64(v.Value))...)
		transBytes = append(transBytes,v.PublicKeyHash...)
	}
	return transBytes
}

func (t *Transaction) customCopy() Transaction {
	newVin := []txInput{}
	newVout := []txOutput{}
	for _, vin := range t.Vint {
		newVin = append(newVin, txInput{vin.TxHash, vin.Index, nil, nil})
	}
	for _, vout := range t.Vout {
		newVout = append(newVout, txOutput{vout.Value, vout.PublicKeyHash})
	}
	return Transaction{t.TxHash, newVin, newVout}
}

func isGenesisTransaction(tss []Transaction) bool {
	if tss != nil {
		if tss[0].Vint[0].Index == -1 {
			return true
		}
	}
	return false
}
