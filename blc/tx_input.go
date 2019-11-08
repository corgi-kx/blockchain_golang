package block

//UTXO输入
type TXInput struct {
	TxHash    []byte
	Index     int
	Signature []byte
	PublicKey []byte
}
