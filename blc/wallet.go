package block

import (
	"bytes"
	"encoding/gob"
	"myCode/public_blockchain/part7-network/database"
)

type addressList [][]byte

func (a *addressList) serliazle() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(a)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

type wallets struct {
	Wallets map[string]*bitcoinKeys
}

func NewWallets(bd *database.BlockchainDB) *wallets {
	w := &wallets{make(map[string]*bitcoinKeys)}
	if database.IsBucketExist(bd, database.AddrBucket) {
		addressList := GetAllAddress(bd)
		for _, v := range *addressList {
			keys := deserialize(bd.View(v, database.AddrBucket), &bitcoinKeys{})
			w.Wallets[string(v)] = keys.(*bitcoinKeys)
		}
		return w
	}
	return w
}

func (w *wallets) GenerateWallet(bd *database.BlockchainDB) string {
	bitcoinKeys := newBitcoinKeys()
	addressByte := bitcoinKeys.getAddress()
	w.storage(addressByte, bitcoinKeys, bd)
	//将地址存入实例
	address := string(addressByte)
	w.Wallets[address] = bitcoinKeys
	return address
}

func (w *wallets) storage(address []byte, keys *bitcoinKeys, bd *database.BlockchainDB) {
	//将公私钥以地址为键 存入数据库
	bd.Put(address, keys.serliazle(), database.AddrBucket)
	//将地址存入地址导航
	listBytes := bd.View([]byte(addrListMapping), database.AddrBucket)
	if listBytes == nil {
		a := addressList{address}
		bd.Put([]byte(addrListMapping), a.serliazle(), database.AddrBucket)
	} else {
		a := deserialize(listBytes, &addressList{})
		addressList := a.(*addressList)
		*addressList = append(*addressList, address)
		bd.Put([]byte(addrListMapping), addressList.serliazle(), database.AddrBucket)
	}
}

//获取全部地址信息
func GetAllAddress(bd *database.BlockchainDB) *addressList {
	listBytes := bd.View([]byte(addrListMapping), database.AddrBucket)
	a := deserialize(listBytes, &addressList{})
	addressList := a.(*addressList)
	return addressList
}
