package block

import (
	"bytes"
	"encoding/gob"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
)

type Transactions struct {
	Ts []Transaction
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