package network

import (
	"bytes"
	"encoding/gob"
	log "github.com/corgi-kx/logcustom"
)

type getBlock struct {
	BlockHash []byte
	AddrFrom  string
}

func (v getBlock) serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(v)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func (v *getBlock) deserialize(d []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(v)
	if err != nil {
		log.Panic(err)
	}
}
