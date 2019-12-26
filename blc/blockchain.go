package block

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/corgi-kx/blockchain_golang/database"
	log "github.com/corgi-kx/logcustom"
	"math/big"
	"os"
	"time"
)

type blockchain struct {
	BD *database.BlockchainDB //封装的blot结构体
}

//创建区块链实例
func NewBlockchain() *blockchain {
	blockchain := blockchain{}
	bd := database.New()
	blockchain.BD = bd
	return &blockchain
}

//创建创世区块交易信息
func (bc *blockchain) CreataGenesisTransaction(address string, value int, send Sender) {
	//判断地址格式是否正确
	if !IsVaildBitcoinAddress(address) {
		log.Errorf("地址格式不正确:%s\n", address)
		return
	}
	//创世区块数据
	txi := TXInput{[]byte{}, -1, nil, nil}
	//本地一定要存创世区块地址的公私钥信息
	wallets := NewWallets(bc.BD)
	genesisKeys, ok := wallets.Wallets[address]
	if !ok {
		log.Fatal("没有找到地址对应的公私钥信息")
	}
	//通过地址获得rip160(sha256(publickey))
	publicKeyHash := generatePublicKeyHash(genesisKeys.PublicKey)
	txo := TXOutput{value, publicKeyHash}
	ts := Transaction{nil, []TXInput{txi}, []TXOutput{txo}}
	ts.hash()
	tss := []Transaction{ts}
	//开始生成区块链的第一个区块
	bc.newGenesisBlockchain(tss)
	//创世区块后,更新本地最新区块为1并,向全网节点发送当前区块链高度1
	NewestBlockHeight = 1
	send.SendVersionToPeers(1)
	fmt.Println("已成生成创世区块")
	//重置utxo数据库，将创世数据存入
	utxos := UTXOHandle{bc}
	utxos.ResetUTXODataBase()
}

//创建区块链
func (bc *blockchain) newGenesisBlockchain(transaction []Transaction) {
	//判断一下是否已生成创世区块
	if len(bc.BD.View([]byte(LastBlockHashMapping), database.BlockBucket)) != 0 {
		log.Fatal("不可重复生成创世区块")
	}
	//生成创世区块
	genesisBlock := newGenesisBlock(transaction)
	//添加到数据库中
	bc.AddBlock(genesisBlock)
}

//创建挖矿奖励地址交易
func (bc *blockchain) CreataRewardTransaction(address string) Transaction {
	if address == "" {
		log.Warn("没有设置挖矿奖励地址，如果出块则不会给予奖励代币")
		return Transaction{}
	}
	if !IsVaildBitcoinAddress(address) {
		log.Warnf("奖励地址格式不正确:%s\n", address)
		return Transaction{}
	}

	publicKeyHash := getPublicKeyHashFromAddress(address)
	txo := TXOutput{TokenRewardNum, publicKeyHash}
	ts := Transaction{nil, nil, []TXOutput{txo}}
	ts.hash()
	return ts
}

//创建UTXO交易实例
func (bc *blockchain) CreateTransaction(from, to string, amount string, send Sender) {
	//判断一下是否已生成创世区块
	if len(bc.BD.View([]byte(LastBlockHashMapping), database.BlockBucket)) == 0 {
		log.Error("还没有生成创世区块，不可进行转账操作 !")
		return
	}
	//检测是否设置了挖矿地址,没设置的话会给出提示
	if len(bc.BD.View([]byte(RewardAddrMapping), database.AddrBucket)) == 0 {
		log.Warn("没有设置挖矿地址，如果挖出区块将不会给予奖励代币!")
	}

	fromSlice := []string{}
	toSlice := []string{}
	amountSlice := []int{}

	//对传入的信息进行校验检测
	err := json.Unmarshal([]byte(from), &fromSlice)
	if err != nil {
		log.Error("json err:", err)
		return
	}
	err = json.Unmarshal([]byte(to), &toSlice)
	if err != nil {
		log.Error("json err:", err)
		return
	}
	err = json.Unmarshal([]byte(amount), &amountSlice)
	if err != nil {
		log.Error("json err:", err)
		return
	}
	if len(fromSlice) != len(toSlice) || len(fromSlice) != len(amountSlice) {
		log.Error("转账数组长度不一致")
		return
	}

	for i, v := range fromSlice {
		if !IsVaildBitcoinAddress(v) {
			log.Errorf(" %s,地址格式不正确！已将此笔交易剔除\n", v)
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
				amountSlice = append(amountSlice[:i], amountSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
				amountSlice = append(amountSlice[:i])
			}
		}
	}

	for i, v := range toSlice {
		if !IsVaildBitcoinAddress(v) {
			log.Errorf(" %s,地址格式不正确！已将此笔交易剔除\n", v)
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
				amountSlice = append(amountSlice[:i], amountSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
				amountSlice = append(amountSlice[:i])
			}
		}
	}
	for i, v := range amountSlice {
		if v < 0 {
			log.Error("转账金额不可小于0，已将此笔交易剔除")
			if i < len(fromSlice)-1 {
				fromSlice = append(fromSlice[:i], fromSlice[i+1:]...)
				toSlice = append(toSlice[:i], toSlice[i+1:]...)
				amountSlice = append(amountSlice[:i], amountSlice[i+1:]...)
			} else {
				fromSlice = append(fromSlice[:i])
				toSlice = append(toSlice[:i])
				amountSlice = append(amountSlice[:i])
			}
		}
	}

	var tss []Transaction
	wallets := NewWallets(bc.BD)
	for index, fromAddress := range fromSlice {
		fromKeys, ok := wallets.Wallets[fromAddress]
		if !ok {
			log.Errorf("没有找到地址%s所对应的公钥,跳过此笔交易", fromAddress)
			continue
		}
		toKeysPublicKeyHash := getPublicKeyHashFromAddress(toSlice[index])
		if fromAddress == toSlice[index] {
			log.Errorf("相同地址不能转账！！！:%s\n", fromAddress)
			return
		}
		u := UTXOHandle{bc}
		//获取数据库中的未消费的utxo
		utxos := u.findUTXOFromAddress(fromAddress)
		if len(utxos) == 0 {
			log.Errorf("%s 余额为0,不能进行转帐操作", fromAddress)
			return
		}
		//将utxos添加上未打包进区块的交易信息
		if tss != nil {
			for _, ts := range tss {
				//先添加未花费utxo 如果有的话就不添加
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

		//打包交易的核心操作
		newTXInput := []TXInput{}
		newTXOutput := []TXOutput{}
		var amount int
		for _, utxo := range utxos {
			amount += utxo.Vout.Value
			newTXInput = append(newTXInput, TXInput{utxo.Hash, utxo.Index, nil, fromKeys.PublicKey})
			if amount > amountSlice[index] {
				tfrom := TXOutput{}
				tfrom.Value = amount - amountSlice[index]
				tfrom.PublicKeyHash = generatePublicKeyHash(fromKeys.PublicKey)
				tTo := TXOutput{}
				tTo.Value = amountSlice[index]
				tTo.PublicKeyHash = toKeysPublicKeyHash
				newTXOutput = append(newTXOutput, tfrom)
				newTXOutput = append(newTXOutput, tTo)
				break
			} else if amount == amountSlice[index] {
				tTo := TXOutput{}
				tTo.Value = amountSlice[index]
				tTo.PublicKeyHash = toKeysPublicKeyHash
				newTXOutput = append(newTXOutput, tTo)
				break
			}
		}
		//如果余额不足则跳过不会打包进入交易
		if amount < amountSlice[index] {
			log.Errorf(" 第%d笔交易%s余额不足", index+1, fromAddress)
			continue
		}
		ts := Transaction{nil, newTXInput, newTXOutput[:]}
		ts.hash()
		tss = append(tss, ts)
	}
	if tss == nil {
		return
	}
	bc.signatureTransactions(tss, wallets)
	//向P2P节点发送交易数据
	send.SendTransToPeers(tss)
}

//交易转账
func (bc *blockchain) Transfer(tss []Transaction, send Sender) {
	//如果是创世区块的交易则无需进行数字签名验证
	if !isGenesisTransaction(tss) {
		//交易的数字签名验证
		bc.verifyTransactionsSign(&tss)
		if len(tss) == 0 {
			log.Error("没有通过的数字签名验证，不予挖矿出块！")
			return
		}
	}
	//进行余额验证
	bc.VerifyTransBalance(&tss)
	if len(tss) == 0 {
		log.Error("没有通过余额验证的交易，不予挖矿出块！")
		return
	}
	//如果设置了奖励地址，则挖矿成功后给予奖励代币
	rewardTs := bc.CreataRewardTransaction(string(bc.BD.View([]byte(RewardAddrMapping), database.AddrBucket)))
	if rewardTs.TxHash != nil {
		tss = append(tss, rewardTs)
	}
	bc.addBlockchain(tss, send)
}

//校验交易余额是否足够,如果不够则剔除
func (bc *blockchain) VerifyTransBalance(tss *[]Transaction) {
	//获取每个地址的UTXO余额，并存入字典中
	var balance = map[string]int{}
	for i := range *tss {
		fromAddress := GetAddressFromPublicKey((*tss)[i].Vint[0].PublicKey)
		//获取数据库中的utxo
		u := UTXOHandle{bc}
		utxos := u.findUTXOFromAddress(fromAddress)
		if len(utxos) == 0 {
			log.Warnf("%s 余额为0！", fromAddress)
			continue
		}
		aomunt := 0
		for _, v := range utxos {
			aomunt += v.Vout.Value
		}
		balance[fromAddress] = aomunt
	}

circle:
	for i := range *tss {
		fromAddress := GetAddressFromPublicKey((*tss)[i].Vint[0].PublicKey)
		u := UTXOHandle{bc}
		utxos := u.findUTXOFromAddress(fromAddress)
		var utxoAmount int //vint将要花费的总utxo
		var voutAmount int //vout剩余的utxo
		var costAmount int //vint将要花费的总utxo减去vout剩余的utxo等于花费的钱数
		//获取每笔vin的值
		for _, vIn := range (*tss)[i].Vint {
			for _, vUTXO := range utxos {
				if bytes.Equal(vIn.TxHash, vUTXO.Hash) && vIn.Index == vUTXO.Index {
					utxoAmount += vUTXO.Vout.Value
				}
			}
		}
		for _, vOut := range (*tss)[i].Vout {
			if bytes.Equal(getPublicKeyHashFromAddress(fromAddress), vOut.PublicKeyHash) {
				voutAmount += vOut.Value
			}
		}
		costAmount = utxoAmount - voutAmount
		if _, ok := balance[fromAddress]; ok {
			balance[fromAddress] -= costAmount
			if balance[fromAddress] < 0 {
				log.Errorf("%s 余额不够，已将此笔交易剔除！", fromAddress)
				*tss = append((*tss)[:i], (*tss)[i+1:]...)
				balance[fromAddress] += costAmount
				goto circle
			}
		} else {
			log.Errorf("%s 余额不够，已将此笔交易剔除！", fromAddress)
			*tss = append((*tss)[:i], (*tss)[i+1:]...)
			goto circle
		}
	}
	log.Debug("已完成UTXO交易余额验证")
}

//设置挖矿奖励地址
func (bc *blockchain) SetRewardAddress(address string) {
	bc.BD.Put([]byte(RewardAddrMapping), []byte(address), database.AddrBucket)
}

//将交易添加进区块链中(内含挖矿操作)
func (bc *blockchain) addBlockchain(transaction []Transaction, send Sender) {
	preBlockbyte := bc.BD.View(bc.BD.View([]byte(LastBlockHashMapping), database.BlockBucket), database.BlockBucket)
	preBlock := Block{}
	preBlock.Deserialize(preBlockbyte)
	height := preBlock.Height + 1
	//进行挖矿
	nb, err := mineBlock(transaction, bc.BD.View([]byte(LastBlockHashMapping), database.BlockBucket), height)
	if err != nil {
		log.Warn(err)
		return
	}
	//将区块添加到本地库中
	bc.AddBlock(nb)
	//将数据同步到UTXO数据库中
	u := UTXOHandle{bc}
	u.Synchrodata(transaction)
	//挖矿出块后 发送高度信息到其他节点
	send.SendVersionToPeers(nb.Height)
}

//添加区块信息到数据库，并更新lastHash
func (bc *blockchain) AddBlock(block *Block) {
	bc.BD.Put(block.Hash, block.Serialize(), database.BlockBucket)
	bci := NewBlockchainIterator(bc)
	currentBlock := bci.Next()
	if currentBlock == nil || currentBlock.Height < block.Height {
		bc.BD.Put([]byte(LastBlockHashMapping), block.Hash, database.BlockBucket)
	}
}

//对交易信息进行数字签名
func (bc *blockchain) signatureTransactions(tss []Transaction, wallets *wallets) {
	for i := range tss {
		copyTs := tss[i].customCopy()
		for index := range tss[i].Vint {
			//获取地址
			bk := bitcoinKeys{nil, tss[i].Vint[index].PublicKey, nil}
			address := bk.getAddress()
			//从数据库或者为打包进数据库的交易数组中,找到vint所对应的交易信息
			trans, err := bc.findTransaction(tss, tss[i].Vint[index].TxHash)
			if err != nil {
				log.Fatal(err)
			}
			copyTs.Vint[index].Signature = nil
			//将拷贝后的交易里面的公钥替换为公钥hash
			copyTs.Vint[index].PublicKey = trans.Vout[tss[i].Vint[index].Index].PublicKeyHash
			//对拷贝后的交易进行整体hash
			copyTs.TxHash = copyTs.hashSign()
			copyTs.Vint[index].PublicKey = nil
			privKey := wallets.Wallets[string(address)].PrivateKey
			//进行签名操作
			tss[i].Vint[index].Signature = ellipticCurveSign(privKey, copyTs.TxHash)
		}
	}
}

//数字签名验证
func (bc *blockchain) verifyTransactionsSign(tss *[]Transaction) {
circle:
	for i := range *tss {
		copyTs := (*tss)[i].customCopy()
		for index, Vin := range (*tss)[i].Vint {
			findTs, err := bc.findTransaction(*tss, Vin.TxHash)
			if err != nil {
				log.Fatal(err)
			}
			//先验证输入地址的公钥hash与指定的utxo输出的公钥hash是否相同
			if !bytes.Equal(findTs.Vout[Vin.Index].PublicKeyHash, generatePublicKeyHash(Vin.PublicKey)) {
				log.Errorf("签名验证失败 %x笔交易的vin并非是本人", (*tss)[i].TxHash)
				*tss = append((*tss)[:i], (*tss)[i+1:]...)
				goto circle
			}
			copyTs.Vint[index].Signature = nil
			copyTs.Vint[index].PublicKey = findTs.Vout[Vin.Index].PublicKeyHash
			copyTs.TxHash = copyTs.hashSign()
			copyTs.Vint[index].PublicKey = nil
			//进行签名验证
			if !ellipticCurveVerify(Vin.PublicKey, Vin.Signature, copyTs.TxHash) {
				log.Errorf("此笔交易：%x没通过签名验证", (*tss)[i].TxHash)
				*tss = append((*tss)[:i], (*tss)[i+1:]...)
				goto circle
			}
		}
	}
	log.Debug("已完成数字签名验证")
}

//查找交易id对应的交易信息
func (bc *blockchain) findTransaction(tss []Transaction, ID []byte) (Transaction, error) {
	//先查找未插入数据库的交易
	if len(tss) != 0 {
		for _, tx := range tss {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return tx, nil
			}
		}
	}
	bci := NewBlockchainIterator(bc)
	//在查找数据库中存在的交易
	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return tx, nil
			}
		}
		//一直迭代到创世区块后退出
		var hashInt big.Int
		hashInt.SetBytes(block.PreHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("FindTransaction err : Transaction is not found")
}

//获取最新区块高度
func (bc *blockchain) GetLastBlockHeight() int {
	bcl := NewBlockchainIterator(bc)
	lastblock := bcl.Next()
	if lastblock == nil {
		return 0
	}
	return lastblock.Height
}

//通过高度获取区块hash
func (bc *blockchain) GetBlockHashByHeight(height int) []byte {
	bcl := NewBlockchainIterator(bc)
	for {
		currentBlock := bcl.Next()
		if currentBlock == nil {
			return nil
		} else if currentBlock.Height == height {
			return currentBlock.Hash
		} else if isGenesisBlock(currentBlock) {
			return nil
		}
	}
}

//通过区块hash获取区块信息
func (bc *blockchain) GetBlockByHash(hash []byte) []byte {
	return bc.BD.View(hash, database.BlockBucket)
}

//传入地址 返回地址余额信息
func (bc *blockchain) GetBalance(address string) int {
	if !IsVaildBitcoinAddress(address) {
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

//查找数据库中全部未花费的UTXO
func (bc *blockchain) findAllUTXOs() map[string][]*UTXO {
	utxosMap := make(map[string][]*UTXO)
	txInputmap := make(map[string][]TXInput)
	bcIterator := NewBlockchainIterator(bc)
	for {
		currentBlock := bcIterator.Next()
		if currentBlock == nil {
			return nil
		}
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

//打印区块链详细信息
func (bc *blockchain) PrintAllBlockInfo() {
	blcIterator := NewBlockchainIterator(bc)
	for {
		block := blcIterator.Next()
		if block == nil {
			log.Error("还未生成创世区块!")
			return
		}
		fmt.Println("========================================================================================================")
		fmt.Printf("本块hash         %x\n", block.Hash)
		fmt.Println("  	------------------------------交易数据------------------------------")
		for _, v := range block.Transactions {
			fmt.Printf("   	 本次交易id:  %x\n", v.TxHash)
			fmt.Println("   	  tx_input：")
			for _, vIn := range v.Vint {
				fmt.Printf("			交易id:  %x\n", vIn.TxHash)
				fmt.Printf("			索引:    %d\n", vIn.Index)
				fmt.Printf("			签名信息:    %x\n", vIn.Signature)
				fmt.Printf("			公钥:    %x\n", vIn.PublicKey)
				fmt.Printf("			地址:    %s\n", GetAddressFromPublicKey(vIn.PublicKey))
			}
			fmt.Println("  	  tx_output：")
			for index, vOut := range v.Vout {
				fmt.Printf("			金额:    %d    \n", vOut.Value)
				fmt.Printf("			公钥Hash:    %x    \n", vOut.PublicKeyHash)
				fmt.Printf("			地址:    %s\n", GetAddressFromPublicKeyHash(vOut.PublicKeyHash))
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
