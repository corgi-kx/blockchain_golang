package network

import (
	"io/ioutil"
	block "myCode/public_blockchain/part7-network/blc"
	log "myCode/public_blockchain/part7-network/logcustom"
	"net"
)

var knowNode = []string{"localhost:3000"}
var localAddr string
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
			log.Debugf("接收到来自%s的信息:", conn.RemoteAddr())
			dataHandle(b)
		}
	}()
	//如果不是中心节点则，启动时向中心节点发送版本信息
	if localAddr != knowNode[0] {
		go func() {
			log.Tracef("启动挖矿节点%s发送消息...",localAddr)
			bc:=block.NewBlockchain()
			lastHeight := bc.GetLastBlockHeight()
			bc.BD.Close()
			v := version{versionInfo, lastHeight, localAddr}
			versionBytes := v.serialize()
			sendMessage(jointMessage(cVersion, versionBytes), knowNode[0])
			log.Trace("挖矿节点已发送完消息...")
		}()
	}
	select {}
}
