package network

import (
	"bytes"
	"encoding/gob"
	log "myCode/public_blockchain/part7-network/logcustom"
)

type hashMap  map[int][]byte

type hash struct {
	HashMap  hashMap
	AddrFrom string
}


func (v hash) serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(v)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func  (v *hash) deserialize(d []byte){
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(v)
	if err != nil {
		log.Panic(err)
	}
}