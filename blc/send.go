package block

type Sender interface {
	SendVersionToPeers(height int)
	SendTransToPeers(tss []Transaction)
}
