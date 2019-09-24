package block

import (
	"io/ioutil"
	log "myCode/public_blockchain/part7-network/logcustom"
	"net"
)

var knowNode = []string{"localhost:3000"}

func StartNode() {
	addr := "localhost:" + nodeID
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()
	if addr == knowNode[0] {
		log.Trace("已启动中心节点监听...")
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
				log.Debug("RemoteAddr:", conn.RemoteAddr())
				log.Debug("接收到的信息：", string(b))
			}
		}()
	} else {
		log.Trace("启动挖矿节点发送消息...")
		conn, err := net.Dial("tcp", knowNode[0])
		if err != nil {
			log.Panic(err)
		}

		_, err = conn.Write([]byte("我是挖矿节点，收到我的信息了吗"))
		if err != nil {
			log.Panic(err)
		}
		conn.Close()
		log.Trace("挖矿节点已发送完消息...")
	}
	select {}
}
