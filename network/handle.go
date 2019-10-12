package network

import (
	blc "github.com/corgi-kx/blockchain_golang/blc"
	"github.com/corgi-kx/blockchain_golang/database"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
	"github.com/corgi-kx/blockchain_golang/send"
	"github.com/corgi-kx/blockchain_golang/util"
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
		utxos.ResetUTXODataBase()
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
	log.Tracef("本节点已将区块数据发送到%s，该块hash为%x",g.AddrFrom,g.BlockHash)
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
	log.Trace("已发送获取hash列表命令")
}

//发送版本信息
func handleVersion(content []byte) {
	log.Trace("进入到发送版本信息界面")
	v:=version{}
	v.deserialize(content)
	bc:=blc.NewBlockchain()
	defer bc.BD.Close()
	lastHeight:=bc.GetLastBlockHeight()
	if lastHeight > v.Height {
		log.Trace("目标高度比本链小，向目标发送版本信息")
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
