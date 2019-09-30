package network

import (
	blc "myCode/public_blockchain/part7-network/blc"
	log "myCode/public_blockchain/part7-network/logcustom"
	"net"
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
	log.Debugf("本节点已接收到来自中心节点的区块数据，该块hash为：%x\n",block.Hash)
	bc := blc.NewBlockchain()
	defer bc.BD.Close()
	pow := blc.NewProofOfWork(block)
	if pow.Verify() {
		bc.AddBlock(block)
		utxos := blc.UTXOHandle{bc}
		utxos.ResetUTXODataBase()
		log.Debugf("POW验证通过,已将区块%x加入数据库,该区块高度为：%d", block.Hash, block.Height)
	} else {
		log.Errorf("POW验证不通过，无法将此块：%x加入数据库",block.Hash)
	}
}
func handleGetBlock(content []byte) {
	g:=getBlock{}
	g.deserialize(content)
	bc:=blc.NewBlockchain()
	defer bc.BD.Close()
	blockBytes:=bc.GetBlockByHash(g.BlockHash)
	data:=jointMessage(cBlock,blockBytes)
	log.Trace("本节点已将区块数据发送到%s，该块hash为%x",g.AddrFrom,g.BlockHash)
	sendMessage(data,g.AddrFrom)
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
		sendMessage(data,h.AddrFrom)
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
	sendMessage(data,g.AddrFrom)
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
		sendMessage(data,v.AddrFrom)
	}else if lastHeight < v.Height {
		log.Debugf("对方版本比咱们大%v,发送获取区块hash信息！",v)
		gh:=getHash{lastHeight,localAddr}
		data:=jointMessage(cGetHash,gh.serialize())
		sendMessage(data,v.AddrFrom)
	}else {
		log.Debug("接收到版本信息，双方高度一致，无需处理！")
	}
}

func sendMessage(data []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	_, err = conn.Write(data)
	if err != nil {
		log.Panic(err)
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
