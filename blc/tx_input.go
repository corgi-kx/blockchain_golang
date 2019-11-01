package block

type TXInput struct {
	TxHash    []byte
	Index     int
	Signature []byte
	PublicKey []byte
}
