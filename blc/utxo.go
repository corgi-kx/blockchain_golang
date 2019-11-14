package block

//UTXO输出的详细信息,便于直接在utxo数据库查找输出
type UTXO struct {
	Hash  []byte
	Index int
	Vout  TXOutput
}
