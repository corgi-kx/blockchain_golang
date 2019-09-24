package block

type txInput struct {
	TxHash    []byte
	Index     int
	Signature []byte
	PublicKey []byte
}
