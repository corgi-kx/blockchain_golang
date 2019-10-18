package network

import (
	blc "github.com/corgi-kx/blockchain_golang/blc"
	log "github.com/corgi-kx/blockchain_golang/logcustom"
	"github.com/corgi-kx/blockchain_golang/send"
	"io/ioutil"
	"net"
)

var centerNode = "localhost:3000"
var knowNodesMap = map[string]string{}
var localAddr string
var tradePool  =send.Transactions{}
func StartNode(nodeID string) {
	localAddr = "localhost:" + nodeID
	listen, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()
	log.Tracef("已启动%s节点监听...", localAddr)
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Panic(err)
			}
			b, err := ioutil.ReadAll(conn)
			if err != nil {
				log.Panic(err)
			}
			go dataHandle(b)
		}
	}()
	//如果不是中心节点则，启动时向中心节点发送版本信息
	if localAddr != centerNode {
		go func() {
			log.Tracef("启动挖矿节点%s发送消息...",localAddr)
			bc:=blc.NewBlockchain()
			lastHeight := bc.GetLastBlockHeight()
			blc.NewestBlockHeight = lastHeight
			v := version{versionInfo, lastHeight, localAddr}
			versionBytes := v.serialize()
			send.SendMessage(jointMessage(cVersion, versionBytes), centerNode)
			log.Trace("挖矿节点已发送完消息...")
		}()
	}
	select {}
}



