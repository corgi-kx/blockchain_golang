package block

//用于network包向对等节点发送信息
type Sender interface {
	SendVersionToPeers(height int)
	SendTransToPeers(tss []Transaction)
}
