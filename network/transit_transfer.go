package network

import (
	"bytes"
	"encoding/gob"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
)

type transfer struct {
	from []string
	to []string
	amount []string
}

func (t *transfer) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(t)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func  (v *transfer) Deserialize(d []byte){
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(v)
	if err != nil {
		log.Panic(err)
	}
}