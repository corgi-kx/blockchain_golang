package network

import (
	"bytes"
	"fmt"
	blc "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/database"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
	"github.com/corgi-kx/blockchain_golang/send"
	"github.com/corgi-kx/blockchain_golang/util"
	"time"
	"unsafe"
)

func dataHandle(data []byte) {
	cmd,content:=splitMessage(data)
	log.Tracef("本节点已接收到命令：%s",cmd)
	switch command(cmd){
	case cVersion:
		handleVersion(content)
	case cGetHash:
		handleGetHash(content)
	case cHashMap:
		handleHashMap(content)
	case cGetBlock:
		handleGetBlock(content)
	case cBlock:
		handleBlock(content)
	case cTransaction :
		handleTransaction(content)
	case cMining :
		go handleMining(content)
	case cMyError :
		go handleMyError(content)
	}
}

func handleMyError(content []byte) {
	log.Error(string(content))
}

//挖矿节点进行挖矿
func handleMining(content []byte) {
	t:=send.Transactions{}
	t.Deserialize(content)
	bc := blc.NewBlockchain()
	defer bc.BD.Close()
	//将send下的transaction转换为blc下的transaction
	nTs:=[]blc.Transaction{}
	for _,v:=range t.Ts {
		nv:=(*blc.Transaction)(unsafe.Pointer(&v))
		nTs = append(nTs,*nv)
	}
	bc.Transfer(nTs)
}

//接收交易信息，满足条件后发送到挖矿节点进行挖矿
func handleTransaction(content []byte) {
	t:=send.Transactions{}
	t.Deserialize(content)
	//先检查交易池里是否存在相同交易地址，有的话则剔除掉（如果发送交易时，在将同一地址进行多笔转账时 写入同一个json则，则允许进行同一地址多笔转账）
	circle:
	for i,_ := range t.Ts {
		for _,v:=range tradePool.Ts {
			if bytes.Equal(t.Ts[i].Vint[0].PublicKey,v.Vint[0].PublicKey) {
				s:=fmt.Sprintf("当前交易池里，已存在此笔地址转账信息(%s)，顾暂不能进行转账，请等待上笔交易出块后在进行此地址转账操作",blc.GetAddressFromPublicKey(t.Ts[i].Vint[0].PublicKey))
				log.Error(s)
				data:=jointMessage(cMyError,[]byte(s))
				send.SendMessage(data,t.Ts[i].AddrFrom)
				t.Ts=append(t.Ts[:i],t.Ts[i+1:]...)
				goto circle
			}
		}
	}
	if len(t.Ts) == 0 {
		log.Error("没有满足条件的转账信息，顾不存入交易池")
		return
	}
	//类型转换为blc下的Transactions
	tradePool.Ts = append(tradePool.Ts,t.Ts...)
	if len(tradePool.Ts) >= tradePoolLength {
		//满足交易池规定的大小后发送交易列表到旷工节点进行挖矿
		log.Debugf("交易池已满足挖矿交易数量大小限制:%d,即将发往挖矿节点",tradePoolLength)
		mineTrans:=send.Transactions{make([]send.Transaction,tradePoolLength)}
		copy(mineTrans.Ts,tradePool.Ts[:tradePoolLength])
		data:=jointMessage(cMining,mineTrans.Serialize())
		for addr,_ :=range knowNodesMap {
			send.SendMessage(data,addr)
			log.Debugf("已将交易列表发送至%s：",addr)
		}
		//将已发送的交易删掉
		newTradePool := send.Transactions{}
		copy(newTradePool.Ts,tradePool.Ts[tradePoolLength:])
		tradePool = newTradePool
	}else {
		log.Debugf("已收到交易信息，当前交易池数量:%d，交易池未满%d，暂不进行挖矿操作",len(tradePool.Ts),tradePoolLength)
	}
}

func handleBlock(content []byte) {
	b := blc.Deserialize(content, &blc.Block{})
	block := b.(*blc.Block)
	log.Debugf("本节点已接收到来自中心节点的区块数据，该块hash为：%x",block.Hash)
	bc := blc.NewBlockchain()
	defer bc.BD.Close()
	pow := blc.NewProofOfWork(block)
	if pow.Verify() {
		bc.AddBlock(block)
		utxos := blc.UTXOHandle{bc}
		//如果是创世区块则重置utxo数据库，否则执行同步操作
		if block.Height == 1 {
			utxos.ResetUTXODataBase()
		}else {
			utxos.Synchrodata(block.Transactions)
		}

		log.Debugf("POW验证通过,已将区块%x加入数据库,该区块高度为：%d", block.Hash, block.Height)
		if localAddr == centerNode {
			log.Debugf("中心节点已更新区块数据,高度为%d！\n",block.Height)
			//如果当前区块高度达到了最新高度，才会像其他节点发送版本信息
			if block.Height == util.BytesToInt(bc.BD.View([]byte(blc.MineNodeLastHeightMapping),database.BlockBucket)) {
				log.Debugf("中心节点已更新至最新区块数据,高度为%d，现在向其他节点发送版本信息！\n",block.Height)
				send.SendVersionToKnows(block.Height,knowNodesMap)
			}
		}
	} else {
		log.Errorf("POW验证不通过，无法将此块：%x加入数据库",block.Hash)
	}
}
func handleGetBlock(content []byte) {
	g:=getBlock{}
	g.deserialize(content)
	bc:=blc.NewBlockchain()
	blockBytes:=bc.GetBlockByHash(g.BlockHash)
	defer bc.BD.Close()
	data:=jointMessage(cBlock,blockBytes)
	log.Debugf("本节点已将区块数据发送到%s，该块hash为%x",g.AddrFrom,g.BlockHash)
	send.SendMessage(data,g.AddrFrom)
}

func handleHashMap(content []byte) {
	h:=hash{}
	h.deserialize(content)
	hm:=h.HashMap
	bc:=blc.NewBlockchain()
	defer bc.BD.Close()
	lastHeight:=bc.GetLastBlockHeight()
	targetHeight:=lastHeight+1
	for {
		hash:=hm[targetHeight]
		if hash == nil {
			break
		}
		g:=getBlock{hash,localAddr}
		data:=jointMessage(cGetBlock,g.serialize())
		send.SendMessage(data,h.AddrFrom)
		log.Debugf("已发送获取区块信息命令,目标高度为：%d",targetHeight)
		targetHeight ++
	}
}

//发送hash字典
func handleGetHash(content []byte) {
	g:=getHash{}
	g.deserialize(content)
	bc:=blc.NewBlockchain()
	defer bc.BD.Close()
	lastHeight:=bc.GetLastBlockHeight()
	hm:=hashMap{}
	for i:=g.Height+1;i<=lastHeight;i++ {
		hm[i] = bc.GetBlockHashByHeight(i)
	}
	h:=hash{hm,localAddr}
	data:=jointMessage(cHashMap,h.serialize())
	send.SendMessage(data,g.AddrFrom)
	log.Debug("已发送获取hash列表命令")
}

//发送版本信息
func handleVersion(content []byte) {
	v:=version{}
	v.deserialize(content)
	//更新网络中节点最新高度值，以方便挖矿线程检测是否还需要进行挖矿操作
    if v.Height > blc.NewestBlockHeight {
		blc.NewestBlockHeight = v.Height
	}
	//睡一下，等待挖矿线程将bolt数据库关闭
	time.Sleep(5*time.Millisecond)
	log.Debugf("所有区块链网络节点中目前最新高度为:%d",blc.NewestBlockHeight)
	bc:=blc.NewBlockchain()
	defer bc.BD.Close()
	lastHeight:=bc.GetLastBlockHeight()
	if lastHeight > v.Height {
		log.Debug("目标高度比本链小，向目标发送版本信息")
		newV:=version{versionInfo,lastHeight,localAddr}
		data:=jointMessage(cVersion,newV.serialize())
		send.SendMessage(data,v.AddrFrom)
	}else if lastHeight < v.Height {
		log.Debugf("对方版本比咱们大%v,发送获取区块hash信息！",v)
		bc.BD.Put([]byte(blc.MineNodeLastHeightMapping),util.Int64ToBytes(int64(v.Height)),database.BlockBucket)   //将最新高度存入数据库，方便中心节点更新区块后发送版本信息
		gh:=getHash{lastHeight,localAddr}
		data:=jointMessage(cGetHash,gh.serialize())
		send.SendMessage(data,v.AddrFrom)
	}else {
		log.Debug("接收到版本信息，双方高度一致，无需处理！")
	}
	//将来信节点加入到已知节点字典里
	if v.AddrFrom != centerNode {
		knowNodesMap[v.AddrFrom] = ""
	}
}


//默认前十二位为命令名称
func jointMessage(cmd command, content []byte) []byte {
	b := make([]byte, prefixCMDLength)
	for i, v := range []byte(cmd) {
		b[i] = v
	}
	joint := make([]byte, 0)
	joint = append(b, content...)
	return joint
}

//默认前十二位为命令名称
func splitMessage(message []byte) (cmd string, content []byte) {
	cmdBytes := message[:prefixCMDLength]
	newCMDBytes := make([]byte, 0)
	for _, v := range cmdBytes {
		if v != byte(0) {
			newCMDBytes = append(newCMDBytes, v)
		}
	}
	cmd = string(newCMDBytes)
	content = message[prefixCMDLength:]
	return
}
