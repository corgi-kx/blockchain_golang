package network

import (
	"bytes"
	"encoding/gob"
	block "github.com/corgi-kx/blockchain_golang/blc"
	log "github.com/corgi-kx/logcustom"
)

type Transactions struct {
	Ts []Transaction
}

type Transaction struct {
	TxHash []byte
	//UTXO输入
	Vint []block.TXInput
	//UTXO输出
	Vout []block.TXOutput

	AddrFrom string
}

func (t *Transactions) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(t)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func (v *Transactions) Deserialize(d []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(v)
	if err != nil {
		log.Panic(err)
	}
}
