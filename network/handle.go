package network

import (
	"bytes"
	"fmt"
	blc "github.com/corgi-kx/blockchain_golang/blc"
	log "github.com/corgi-kx/logcustom"
	"github.com/libp2p/go-libp2p-core/network"
	"io/ioutil"
	"sync"
	"time"
)

//对接收到的数据解析出命令,然后对不同的命令分别进行处理
func handleStream(stream network.Stream) {
	data, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Panic(err)
	}
	//取信息的前十二位得到命令
	cmd, content := splitMessage(data)
	log.Tracef("本节点已接收到命令：%s", cmd)
	switch command(cmd) {
	case cVersion:
		go handleVersion(content)
	case cGetHash:
		go handleGetHash(content)
	case cHashMap:
		go handleHashMap(content)
	case cGetBlock:
		go handleGetBlock(content)
	case cBlock:
		go handleBlock(content)
	case cTransaction:
		go handleTransaction(content)
	case cMyError:
		go handleMyError(content)
	}
}

//打印接收到的错误信息
func handleMyError(content []byte) {
	e := myerror{}
	e.deserialize(content)
	log.Warn(e.Error)
	peer := buildPeerInfoByAddr(e.Addrfrom)
	delete(peerPool, fmt.Sprint(peer.ID))
}

//接收交易信息，满足条件后进行挖矿
func handleTransaction(content []byte) {
	t := Transactions{}
	t.Deserialize(content)
	if len(t.Ts) == 0 {
		log.Error("没有满足条件的转账信息，顾不存入交易池")
		return
	}
	//交易池中只能存在每个地址的一笔交易信息
	//判断当前交易池中是否已有该地址发起的交易
	if len(tradePool.Ts) != 0 {
	circle:
		for i := range t.Ts {
			for _, v := range tradePool.Ts {
				if bytes.Equal(t.Ts[i].Vint[0].PublicKey, v.Vint[0].PublicKey) {
					s := fmt.Sprintf("当前交易池里，已存在此笔地址转账信息(%s)，顾暂不能进行转账，请等待上笔交易出块后在进行此地址转账操作", blc.GetAddressFromPublicKey(t.Ts[i].Vint[0].PublicKey))
					log.Error(s)
					t.Ts = append(t.Ts[:i], t.Ts[i+1:]...)
					goto circle
				}
			}
		}
	}

	if len(t.Ts) == 0 {
		return
	}
	mineBlock(t)
}

//调用区块模块进行挖矿操作
var lock = sync.Mutex{}

func mineBlock(t Transactions) {
	//锁上,等待上一个挖矿结束后才进行挖矿!
	lock.Lock()
	defer lock.Unlock()
	//将临时交易池的交易添加进交易池
	tradePool.Ts = append(tradePool.Ts, t.Ts...)

	for {
		//满足交易池规定的大小后进行挖矿
		if len(tradePool.Ts) >= TradePoolLength {
			log.Debugf("交易池已满足挖矿交易数量大小限制:%d,即将进行挖矿", TradePoolLength)
			mineTrans := Transactions{make([]Transaction, TradePoolLength)}
			copy(mineTrans.Ts, tradePool.Ts[:TradePoolLength])

			bc := blc.NewBlockchain()
			//如果当前节点区块高度小于网络最新高度，则等待节点更新区块后在进行挖矿
			for {
				currentHeight := bc.GetLastBlockHeight()
				if currentHeight >= blc.NewestBlockHeight {
					break
				}
				time.Sleep(time.Second * 1)
			}
			//将network下的transaction转换为blc下的transaction
			nTs := make([]blc.Transaction, len(mineTrans.Ts))
			for i := range mineTrans.Ts {
				nTs[i].TxHash = mineTrans.Ts[i].TxHash
				nTs[i].Vint = mineTrans.Ts[i].Vint
				nTs[i].Vout = mineTrans.Ts[i].Vout
			}
			//进行转帐挖矿
			bc.Transfer(nTs, send)
			//剔除已打包进区块的交易
			newTrans := []Transaction{}
			newTrans = append(newTrans, tradePool.Ts[TradePoolLength:]...)
			tradePool.Ts = newTrans
		} else {
			log.Infof("当前交易池数量:%d，交易池未满%d，暂不进行挖矿操作", len(tradePool.Ts), TradePoolLength)
			break
		}
	}
}

//接收到区块数据,进行验证后加入数据库
func handleBlock(content []byte) {
	block := &blc.Block{}
	block.Deserialize(content)
	log.Infof("本节点已接收到来自其他节点的区块数据，该块hash为：%x", block.Hash)
	bc := blc.NewBlockchain()
	pow := blc.NewProofOfWork(block)
	//重新计算本块hash,进行pow验证
	if pow.Verify() {
		log.Infof("POW验证通过,该区块高度为：%d", block.Height)
		//如果是创世区块则直接添加进本地库
		currentHash := bc.GetBlockHashByHeight(block.Height)
		if block.Height == 1 && currentHash == nil {
			bc.AddBlock(block)
			utxos := blc.UTXOHandle{bc}
			utxos.ResetUTXODataBase() //重置utxo数据库
			log.Info("创世区块验证通过,已存入本地数据库...")
		}
		//验证上一个区块的hash与本块中prehash是否一致
		lastBlockHash := bc.GetBlockHashByHeight(block.Height - 1)
		if lastBlockHash == nil {
			//如果找不到上一个区块,可能是还未同步,建立个循环等待同步
			for {
				time.Sleep(time.Second)
				lastBlockHash = bc.GetBlockHashByHeight(block.Height - 1)
				if lastBlockHash != nil {
					log.Debugf("区块高度%d尚未同步,等待同步...", block.Height-1)
					break
				}
			}
		}
		//如果上一块的hash等于本块prehash则通过存入本地库
		if bytes.Equal(lastBlockHash, block.PreHash) {
			bc.AddBlock(block)
			utxos := blc.UTXOHandle{bc}
			//重置utxo数据库
			utxos.ResetUTXODataBase()
			log.Infof("prehash验证通过,该区块高度为:%d,", block.Height)
			log.Infof("总验证通过已存入本地库,区块高度%d,哈希%x", block.Height, block.Hash)
		} else {
			log.Infof("上一个块高度为%d的hash值为:%x,与本块中的prehash值:%x不一致,固不存入区块链中", block.Height-1, lastBlockHash, block.Hash)
		}
	} else {
		log.Errorf("POW验证不通过，无法将此块：%x加入数据库", block.Hash)
	}
}

//接收到获取区块命令,通过hash值 找到该区块 然后把该区块发送过去
func handleGetBlock(content []byte) {
	g := getBlock{}
	g.deserialize(content)
	bc := blc.NewBlockchain()
	blockBytes := bc.GetBlockByHash(g.BlockHash)
	data := jointMessage(cBlock, blockBytes)
	log.Debugf("本节点已将区块数据发送到%s，该块hash为%x", g.AddrFrom, g.BlockHash)
	send.SendMessage(buildPeerInfoByAddr(g.AddrFrom), data)
}

//从对面节点处获取到本地区块链所没有的区块hash列表,然后依次发送"获取区块命令"到该节点
func handleHashMap(content []byte) {
	h := hash{}
	h.deserialize(content)
	hm := h.HashMap
	bc := blc.NewBlockchain()
	lastHeight := bc.GetLastBlockHeight()
	targetHeight := lastHeight + 1
	for {
		hash := hm[targetHeight]
		if hash == nil {
			break
		}
		g := getBlock{hash, localAddr}
		data := jointMessage(cGetBlock, g.serialize())
		send.SendMessage(buildPeerInfoByAddr(h.AddrFrom), data)
		log.Debugf("已发送获取区块信息命令,目标高度为：%d", targetHeight)
		targetHeight++
	}
}

//接收到"获取hash列表"命令,返回对面节点所没有的区块的hash信息(两条链的高度差)
func handleGetHash(content []byte) {
	g := getHash{}
	g.deserialize(content)
	bc := blc.NewBlockchain()
	lastHeight := bc.GetLastBlockHeight()
	hm := hashMap{}
	for i := g.Height + 1; i <= lastHeight; i++ {
		hm[i] = bc.GetBlockHashByHeight(i)
	}
	h := hash{hm, localAddr}
	data := jointMessage(cHashMap, h.serialize())
	send.SendMessage(buildPeerInfoByAddr(g.AddrFrom), data)
	log.Debug("已发送获取hash列表命令")
}

//接收到其他节点的区块高度信息,与本地区块高度进行对比
func handleVersion(content []byte) {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()
	v := version{}
	v.deserialize(content)
	bc := blc.NewBlockchain()
	if blc.NewestBlockHeight > v.Height {
		log.Info("目标高度比本链小，准备向目标发送版本信息")
		for {
			currentHeight := bc.GetLastBlockHeight()
			if currentHeight < blc.NewestBlockHeight {
				log.Info("当前正在更新区块信息,稍后将发送版本信息...")
				time.Sleep(time.Second)
			} else {
				newV := version{versionInfo, currentHeight, localAddr}
				data := jointMessage(cVersion, newV.serialize())
				send.SendMessage(buildPeerInfoByAddr(v.AddrFrom), data)
				break
			}
		}
	} else if blc.NewestBlockHeight < v.Height {
		log.Debugf("对方版本比咱们大%v,发送获取区块的hash信息！", v)
		gh := getHash{blc.NewestBlockHeight, localAddr}
		blc.NewestBlockHeight = v.Height
		data := jointMessage(cGetHash, gh.serialize())
		send.SendMessage(buildPeerInfoByAddr(v.AddrFrom), data)
	} else {
		log.Debug("接收到版本信息，双方高度一致，无需处理！")
	}
}
