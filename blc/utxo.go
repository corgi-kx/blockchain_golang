package block

type UTXO struct {
	Hash  []byte
	Index int
	Vout  *txOutput
}
