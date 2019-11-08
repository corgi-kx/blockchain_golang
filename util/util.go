package util

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
)

//int64转换成字节数组
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

//字节数组转换为int
func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

//生成随机数
func GenerateRealRandom() int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000000000000000))
	if err != nil {
		fmt.Println(err)
	}
	return n.Int64()
}
