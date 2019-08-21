package BLK
type TXInput struct {
	//承接的交易
	Hash []byte
	//承接交易的哪个输出
	vout int64
	//解锁脚本
	ScriptSig string
}