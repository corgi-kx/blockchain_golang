package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"myCode/public_blockchain/part7-network/util"
)

type transaction struct {
	TxHash []byte
	//UTXO输入
	Vint []*txInput
	//UTXO输出
	Vout []*txOutput
}

func (t *transaction) hash() {
	tBytes := t.serialize()
	//加入随机数byte
	randomNumber := util.GenerateRealRandom()
	randomByte := util.Int64ToBytes(randomNumber)
	sumByte := bytes.Join([][]byte{tBytes, randomByte}, []byte(""))
	hashByte := sha256.Sum256(sumByte)
	t.TxHash = hashByte[:]
}

func (t *transaction) hashSign() []byte {
	t.TxHash = nil
	tBytes := t.serialize()
	//加入随机数byte
	hashByte := sha256.Sum256(tBytes)
	return hashByte[:]
}

// 将Block对象序列化成[]byte
func (t *transaction) serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(t)
	if err != nil {
		panic(err)
	}

	return result.Bytes()
}

func (t *transaction) customCopy() *transaction {
	newVin := []*txInput{}
	newVout := []*txOutput{}
	for _, vin := range t.Vint {
		newVin = append(newVin, &txInput{vin.TxHash, vin.Index, nil, nil})
	}
	for _, vout := range t.Vout {
		newVout = append(newVout, &txOutput{vout.Value, vout.PublicKeyHash})
	}
	return &transaction{t.TxHash, newVin, newVout}
}

func isGenesisTransaction(tss []*transaction) bool {
	if tss != nil {
		if tss[0].Vint[0].Index == -1 {
			return true
		}
	}
	return false
}
