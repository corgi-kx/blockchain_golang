package block

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"myCode/public_blockchain/part7-network/database"
	log "myCode/public_blockchain/part7-network/logcustom"
	"os"
	"time"
)

type blockchain struct {
	BD *database.BlockchainDB
}

func NewBlockchain() *blockchain {
	if !database.IsBlotExist(nodeID) {
		log.Fatal(" 没有找到对应的数据库！")
	}
	blockchain := blockchain{}
	bd := database.New(nodeID)
	blockchain.BD = bd
	return &blockchain
}

func (bc *blockchain) CreataGenesisTransaction(address string, value int) {
	if !isVaildBitcoinAddress(address) {
		log.Errorf("地址格式不正确:%s\n", address)
		return
	}
	//创世区块数据
	txi := txInput{[]byte{}, -1, nil, nil}
	wallets := NewWallets(bc.BD)
	genesisKeys := wallets.Wallets[address]
	if genesisKeys == nil {
		log.Fatal("没有找到地址对应的公钥信息")
	}
	//通过地址获取rip160(sha256(publickey))
	publicKeyHash := generatePublicKeyHash(genesisKeys.PublicKey)
	txo := txOutput{value, publicKeyHash}
	ts := transaction{nil, []*txInput{&txi}, []*txOutput{&txo}}
	ts.hash()
	tss := []*transaction{&ts}
	bc.NewGenesisBlockchain(tss)
	//重置未花费数据库，将创世数据存入
	utxos := UTXOHandle{bc}
	utxos.ResetUTXODataBase()
}

func (bc *blockchain) CreataRewardTransaction(address string) *transaction {
	if !isVaildBitcoinAddress(address) {
		log.Errorf("奖励地址格式不正确:%s\n", address)
		return nil
	}

	publicKeyHash := getPublicKeyHashFromAddress(address)
	txo := txOutput{tokenRewardNum, publicKeyHash}
	ts := transaction{nil, nil, []*txOutput{&txo}}
	ts.hash()
	return &ts
}
func (bc *blockchain) SetRewardAddress(address string) {
	bc.BD.Put([]byte(rewardAddrMapping), []byte(address), database.AddrBucket)
}
func (bc *blockchain) Transfer(from, to string, amount string) {
	//判断一下是否已生成创世区块
	if bc.BD.View([]byte(lastBlockHashMapping), database.BlockBucket) == nil {
		log.Fatal("还没有生成创世区块，不可进行转账操作 !")
	}
	if bc.BD.View([]byte(rewardAddrMapping), database.AddrBucket) == nil {
		log.Fatal("没有设置挖矿奖励地址，请前往设置!")
	}
	fromSlice := []string{}
	toSlice := []string{}
	amountSlice := []int{}
	err := json.Unmarshal([]byte(from), &fromSlice)
	if err != nil {
		log.Fatal("json err:", err)
	}
	err = json.Unmarshal([]byte(to), &toSlice)
	if err != nil {
		log.Fatal("json err:", err)
	}
	err = json.Unmarshal([]byte(amount), &amountSlice)
	if err != nil {
		log.Fatal("json err:", err)
	}
	if len(fromSlice) != len(toSlice) || len(fromSlice) != len(amountSlice) {
		log.Fatal("json err:", err)
	}

	for i, v := range fromSlice {
		if !isVaildBitcoinAddress(v) {
			log.Errorf(" %s,地址格式不正确！已将此笔交易剔除\n", v)
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
			}
		}
	}

	for i, v := range toSlice {
		if !isVaildBitcoinAddress(v) {
			log.Errorf(" %s,地址格式不正确！已将此笔交易剔除\n", v)
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
			}
		}
	}

	tss := bc.createTransaction(fromSlice, toSlice, amountSlice)
	if tss == nil {
		log.Error("没有符合规则的交易，不进行出块操作")
		return
	}
	rewardTs := bc.CreataRewardTransaction(string(bc.BD.View([]byte(rewardAddrMapping), database.AddrBucket)))
	tss = append(tss, rewardTs)
	bc.addBlockchain(tss)
}

func (bc *blockchain) createTransaction(fromSlice, toSlice []string, amountSlice []int) []*transaction {
	var tss []*transaction
	wallets := NewWallets(bc.BD)
	for index, fromAddress := range fromSlice {
		fromKeys := wallets.Wallets[fromAddress]
		if fromKeys == nil {
			log.Errorf("没有找到地址%s所对应的公钥,跳过此笔交易\n", fromAddress)
			continue
		}
		toKeysPublicKeyHash := getPublicKeyHashFromAddress(toSlice[index])
		if fromAddress == toSlice[index] {
			log.Errorf("相同地址不能转账！！！:%s\n", fromAddress)
			return nil
		}
		u := UTXOHandle{bc}
		//获取数据库中的utxo
		utxos := u.findUTXOFromAddress(fromAddress)
		if utxos == nil {
			log.Errorf("%s 余额为0", fromAddress)
			return nil
		}
		//将utxos添加上未打包进区块的交易切片
		if tss != nil {
			for _, ts := range tss {
				//先添加未花费utxo 如果有就不加了
			tagVout:
				for index, vOut := range ts.Vout {
					if bytes.Compare(vOut.PublicKeyHash, generatePublicKeyHash(fromKeys.PublicKey)) != 0 {
						continue
					}
					for _, utxo := range utxos {
						if bytes.Equal(ts.TxHash, utxo.Hash) && index == utxo.Index {
							continue tagVout
						}
					}
					utxos = append(utxos, &UTXO{ts.TxHash, index, vOut})
				}
				//剔除已花费的utxo
				for _, vInt := range ts.Vint {
					for index, utxo := range utxos {
						if bytes.Equal(vInt.TxHash, utxo.Hash) && vInt.Index == utxo.Index {
							utxos = append(utxos[:index], utxos[index+1:]...)
						}
					}
				}

			}
		}

		newTXInput := []*txInput{}
		newTXOutput := []*txOutput{}
		var amount int
		for _, utxo := range utxos {
			amount += utxo.Vout.Value
			newTXInput = append(newTXInput, &txInput{utxo.Hash, utxo.Index, nil, fromKeys.PublicKey})
			if amount > amountSlice[index] {
				tfrom := &txOutput{}
				tfrom.Value = amount - amountSlice[index]
				tfrom.PublicKeyHash = generatePublicKeyHash(fromKeys.PublicKey)
				tTo := &txOutput{}
				tTo.Value = amountSlice[index]
				tTo.PublicKeyHash = toKeysPublicKeyHash
				newTXOutput = append(newTXOutput, tfrom)
				newTXOutput = append(newTXOutput, tTo)
				break
			} else if amount == amountSlice[index] {
				tTo := &txOutput{}
				tTo.Value = amountSlice[index]
				tTo.PublicKeyHash = toKeysPublicKeyHash
				newTXOutput = append(newTXOutput, tTo)
				break
			}
		}
		if amount < amountSlice[index] {
			log.Errorf(" 第%d笔交易%s余额不足\n", index+1, fromAddress)
			continue
		}
		ts := &transaction{nil, newTXInput, newTXOutput[:]}
		ts.hash()
		tss = append(tss, ts)
	}
	if tss == nil {
		return nil
	}
	bc.signatureTransactions(tss, wallets)
	return tss
}

func (bc *blockchain) GetBalance(address string) int {
	if !isVaildBitcoinAddress(address) {
		log.Errorf("地址格式不正确：%s\n", address)
		os.Exit(0)
	}
	var balance int
	uHandle := UTXOHandle{bc}
	utxos := uHandle.findUTXOFromAddress(address)
	for _, v := range utxos {
		balance += v.Vout.Value
	}
	return balance
}

func (bc *blockchain) findAllUTXOs() map[string][]*UTXO {
	utxosMap := make(map[string][]*UTXO)
	txInputmap := make(map[string][]*txInput)
	bcIterator := bc.newBlockchainIterator()
	for {
		currentBlock := bcIterator.next()
		//必须倒序 否则有的已花费不会被扣掉
		for i := len(currentBlock.Transactions) - 1; i >= 0; i-- {
			var utxos = []*UTXO{}
			ts := currentBlock.Transactions[i]
			for _, vInt := range ts.Vint {
				txInputmap[string(vInt.TxHash)] = append(txInputmap[string(vInt.TxHash)], vInt)
			}

		VoutTag:
			for index, vOut := range ts.Vout {
				if txInputmap[string(ts.TxHash)] == nil {
					utxos = append(utxos, &UTXO{ts.TxHash, index, vOut})
				} else {
					for _, vIn := range txInputmap[string(ts.TxHash)] {
						if vIn.Index == index {
							continue VoutTag
						}
					}
					utxos = append(utxos, &UTXO{ts.TxHash, index, vOut})
				}
				utxosMap[string(ts.TxHash)] = utxos
			}
		}

		if isGenesisBlock(currentBlock) {
			break
		}
	}
	return utxosMap
}

func (bc *blockchain) NewGenesisBlockchain(transaction []*transaction) {
	//判断一下是否已生成创世区块
	if bc.BD.View([]byte(lastBlockHashMapping), database.BlockBucket) != nil {
		log.Fatal("不可重复生成创世区块")
	}
	genesisBlock := newGenesisBlock(transaction)
	bc.BD.Put(genesisBlock.Hash, genesisBlock.serialize(), database.BlockBucket)
	bc.BD.Put([]byte(lastBlockHashMapping), genesisBlock.Hash, database.BlockBucket)
}

func (blockchain *blockchain) addBlockchain(transaction []*transaction) {
	//进行数字签名验证
	if !isGenesisTransaction(transaction) {
		if !blockchain.verifyTransactions(transaction) {
			os.Exit(0)
		}
	}
	preBlockbyte := blockchain.BD.View(blockchain.BD.View([]byte(lastBlockHashMapping), database.BlockBucket), database.BlockBucket)
	p := deserialize(preBlockbyte, &block{})
	preBlock := p.(*block)
	height := preBlock.Height + 1
	newblock := newBlock(transaction, blockchain.BD.View([]byte(lastBlockHashMapping), database.BlockBucket), height)
	blockchain.BD.Put(newblock.Hash, newblock.serialize(), database.BlockBucket)
	blockchain.BD.Put([]byte(lastBlockHashMapping), newblock.Hash, database.BlockBucket)
	//将数据同步到UTXO数据库中
	u := UTXOHandle{blockchain}
	u.synchrodata(transaction)
}

func (bc *blockchain) signatureTransactions(tss []*transaction, wallets *wallets) {
	for _, ts := range tss {
		copyTs := ts.customCopy()
		for index, vIn := range ts.Vint {
			//获取地址
			bk := bitcoinKeys{nil, vIn.PublicKey}
			address := bk.getAddress()
			trans, err := bc.findTransaction(tss, vIn.TxHash)
			if err != nil {
				log.Fatal(err)
			}
			copyTs.Vint[index].Signature = nil
			copyTs.Vint[index].PublicKey = trans.Vout[vIn.Index].PublicKeyHash
			copyTs.TxHash = copyTs.hashSign()
			copyTs.Vint[index].PublicKey = nil
			privKey := wallets.Wallets[string(address)].PrivateKey
			vIn.Signature = ellipticCurveSign(privKey, copyTs.TxHash)
		}
	}
}

func (bc *blockchain) verifyTransactions(tss []*transaction) bool {
	for _, ts := range tss {
		copyTs := ts.customCopy()
		for index, Vin := range ts.Vint {
			findTs, err := bc.findTransaction(tss, Vin.TxHash)
			if err != nil {
				log.Fatal(err)
			}
			if !bytes.Equal(findTs.Vout[Vin.Index].PublicKeyHash, generatePublicKeyHash(Vin.PublicKey)) {
				log.Errorf("签名验证失败 %x笔交易的vin并非是本人", ts.TxHash)
				return false
			}
			copyTs.Vint[index].Signature = nil
			copyTs.Vint[index].PublicKey = findTs.Vout[Vin.Index].PublicKeyHash
			copyTs.TxHash = copyTs.hashSign()
			copyTs.Vint[index].PublicKey = nil
			if !ellipticCurveVerify(Vin.PublicKey, Vin.Signature, copyTs.TxHash) {
				log.Errorf("此笔交易：%x没通过签名验证", ts.TxHash)
				return false
			}
		}
	}
	log.Debug("数字签名已验证通过")
	return true
}

/**
  查找交易id对应的交易信息
  先查找未插入数据库的交易切片
  在查找已插入数据库的交易
*/
func (bc *blockchain) findTransaction(tss []*transaction, ID []byte) (transaction, error) {
	if tss != nil {
		for _, tx := range tss {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return *tx, nil
			}
		}
	}
	bci := bc.newBlockchainIterator()
	for {
		block := bci.next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PreHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

	return transaction{}, errors.New("FindTransaction err : Transaction is not found")
}

func (bc *blockchain) newBlockchainIterator() *blockchainIterator {
	blockchainIterator := &blockchainIterator{bc.BD.View([]byte(lastBlockHashMapping), database.BlockBucket), bc.BD}
	return blockchainIterator
}

func (bc *blockchain) PrintAllBlockInfo() {
	blcIterator := bc.newBlockchainIterator()
	for {
		block := blcIterator.next()
		fmt.Println("========================================================================================================")
		fmt.Printf("本块hash         %x\n", block.Hash)
		fmt.Println("  	------------------------------交易数据------------------------------")
		for _, v := range block.Transactions {
			fmt.Printf("   	 本次交易id:  %x\n", v.TxHash)
			fmt.Println("   	  已花费UTXO：")
			for _, vIn := range v.Vint {
				fmt.Printf("			交易id:  %x\n", vIn.TxHash)
				fmt.Printf("			索引:    %d\n", vIn.Index)
				fmt.Printf("			签名信息:    %x\n", vIn.Signature)
				fmt.Printf("			公钥:    %x\n", vIn.PublicKey)
			}
			fmt.Println("  	  未花费UTXO：")
			for index, vOut := range v.Vout {
				fmt.Printf("			金额:    %d    \n", vOut.Value)
				fmt.Printf("			公钥Hash:    %x    \n", vOut.PublicKeyHash)
				if len(v.Vout) != 1 && index != len(v.Vout)-1 {
					fmt.Println("			---------------")
				}
			}
		}
		fmt.Println("  	--------------------------------------------------------------------")
		fmt.Printf("时间戳           %s\n", time.Unix(block.TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("区块高度         %d\n", block.Height)
		fmt.Printf("随机数           %d\n", block.Nonce)
		fmt.Printf("上一个块hash     %x\n", block.PreHash)
		var hashInt big.Int
		hashInt.SetBytes(block.PreHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	fmt.Println("========================================================================================================")
}
