/*
	utxo数据库创建的意义在于,不会每次进行转帐时遍历整个区块链,
	而是去utxo数据库查找未消费的交易输出,这样会大大降低性能问题
*/
package block

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/corgi-kx/blockchain_golang/database"
	log "github.com/corgi-kx/logcustom"
)

type UTXOHandle struct {
	BC *blockchain
}

//重置UTXO数据库
func (u *UTXOHandle) ResetUTXODataBase() {
	//先查找全部未花费UTXO
	utxosMap := u.BC.findAllUTXOs()
	if utxosMap == nil {
		log.Debug("找不到区块,暂不重置UTXO数据库")
		return
	}
	//删除旧的UTXO数据库
	if database.IsBucketExist(u.BC.BD, database.UTXOBucket) {
		u.BC.BD.DeleteBucket(database.UTXOBucket)
	}
	//创建并将未花费UTXO循环添加
	for k, v := range utxosMap {
		u.BC.BD.Put([]byte(k), u.serialize(v), database.UTXOBucket)
	}
}

//根据地址未消费的utxo
func (u *UTXOHandle) findUTXOFromAddress(address string) []*UTXO {
	publicKeyHash := getPublicKeyHashFromAddress(address)
	utxosSlice := []*UTXO{}
	//获取bolt迭代器，遍历整个UTXO数据库
	//打开数据库
	var DBFileName = "blockchain_" + ListenPort + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(database.UTXOBucket))
		if b == nil {
			return errors.New("datebase view err: not find bucket ")
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			utxos := u.dserialize(v)
			for _, utxo := range utxos {
				if bytes.Equal(utxo.Vout.PublicKeyHash, publicKeyHash) {
					utxosSlice = append(utxosSlice, utxo)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//关闭数据库
	err = db.Close()
	if err != nil {
		log.Panic("db close err :", err)
	}
	return utxosSlice
}

//传入交易信息,将交易里的输出添加进utxo数据库,并剔除输入信息
func (u *UTXOHandle) Synchrodata(tss []Transaction) {
	//先将全部输入插入数据库
	for _, ts := range tss {
		utxos := []*UTXO{}
		for index, vOut := range ts.Vout {
			utxos = append(utxos, &UTXO{ts.TxHash, index, vOut})
		}
		u.BC.BD.Put(ts.TxHash, u.serialize(utxos), database.UTXOBucket)
	}

	//在用输出进行剔除
	for _, ts := range tss {
		for _, vIn := range ts.Vint {
			publicKeyHash := generatePublicKeyHash(vIn.PublicKey)
			//获取bolt迭代器，遍历整个UTXO数据库
			utxoByte := u.BC.BD.View(vIn.TxHash, database.UTXOBucket)
			if len(utxoByte) == 0 {
				log.Panic("Synchrodata err : do not find utxo")
			}
			utxos := u.dserialize(utxoByte)
			newUTXO := []*UTXO{}
			for _, utxo := range utxos {
				if utxo.Index == vIn.Index && bytes.Equal(utxo.Vout.PublicKeyHash, publicKeyHash) {
					continue
				}
				newUTXO = append(newUTXO, utxo)
			}
			u.BC.BD.Delete(vIn.TxHash, database.UTXOBucket)
			u.BC.BD.Put(vIn.TxHash, u.serialize(newUTXO), database.UTXOBucket)
		}
	}
}

func (u *UTXOHandle) serialize(utxos []*UTXO) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(utxos)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func (u *UTXOHandle) dserialize(d []byte) []*UTXO {
	var model []*UTXO
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&model)
	if err != nil {
		log.Panic(err)
	}
	return model
}
