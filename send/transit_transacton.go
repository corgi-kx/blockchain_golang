package send

import (
	"bytes"
	"encoding/gob"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
)

type Transactions struct {
	Ts []Transaction
}

type Transaction struct {
	TxHash []byte
	//UTXO输入
	Vint []txInput
	//UTXO输出
	Vout []txOutput

	AddrFrom string
}
type txInput struct {
	TxHash    []byte
	Index     int
	Signature []byte
	PublicKey []byte
}

type txOutput struct {
	Value         int
	PublicKeyHash []byte
}



// 将transaction序列化成[]byte
func (t *Transactions) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(t)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func  (v *Transactions) Deserialize(d []byte){
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(v)
	if err != nil {
		log.Panic(err)
	}
}