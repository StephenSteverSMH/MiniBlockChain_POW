package BLK

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"unsafe"
)

type TX struct {
	Hash []byte
	//输出
	VOuts []*TXOutput
	//输入
	VIns []*TXInput
}

//简易交易池
type TXPOOL []TX
var GLOBAL_TXPOOL *TXPOOL = new(TXPOOL)
//JSON序列化交易池
func (txpool *TXPOOL)Serialize() string{
	data, err := json.Marshal(txpool)
	if err!=nil{
		fmt.Printf("json.Marshal,err:",err);
	}
	return string(data)
}

//JSON反序列化交易池
func (txpool *TXPOOL)Deserialize(str string){
	err := json.Unmarshal([]byte(str), txpool)
	if err!=nil{
		fmt.Printf("json.Marshal,err:",err);
	}
}

//根据交易HASH寻找对应交易
//待写
func GetTxByHash(hash []byte) *Tx{
	return nil
}

//生成交易HASH
func GenerateTXHash(tx TX) []byte{
	var hash [32]byte
	TXOutput_bytes := *(*[]byte)(unsafe.Pointer(&tx.VOuts))
	TXInput_bytes := *(*[]byte)(unsafe.Pointer(&tx.VIns))
	tx_bytes := bytes.Join([][]byte{tx.Hash, TXOutput_bytes, TXInput_bytes}, []byte{})
	hash = sha256.Sum256(tx_bytes)
	return hash[:]
}

//签名
func SignTX(tx TX){
	var temp_VIns []*TXInput = []*TXInput{}
	for i:=0;i<len(tx.VIns);i++{
		temp_in := *tx.VIns[i]
		//修改输入中的解锁脚本为上一个交易输出中的加密脚本，待写
		//。。。
		temp_VIns = append(temp_VIns, &temp_in)
	}
	tx.VIns = temp_VIns

}
//生成交易
func NewTX(VOuts []*TXOutput, VIns []*TXInput) TX{
	var hash [32]byte
	tx := TX{hash[:], VOuts, VIns}
	//对输入脚本进行签名

	tx.Hash = GenerateTXHash(tx)
	return tx
}
//生成coinbase交易
func NewCoinbaseTX(address string) TX{
	var hash [32]byte
	txInput := TXInput{hash[:], -1, "genesis Script"}
	txOutput := TXOutput{1000, address}
	tx := NewTX([]*TXOutput{&txOutput}, []*TXInput{&txInput})
	//将交易放入全局交易池
	*GLOBAL_TXPOOL = append(*GLOBAL_TXPOOL, tx)
	return tx;
}
//生成普通交易 一输入多输出
func NewSimpleTransaction(from string, to []string, amount []int64) TX{
	var sum int64 = 0
	for _, piece := range amount{
		sum += piece
	}
	//查看发起者余额
	balance := GLOBAL_UTXOS.GetBalance(from)
	//确认发起者的余额足够
	if(balance < sum){
		log.Panic("账户余额不足")
	}
	//推算应该消耗发起者的哪些utxo
	USER_UTXOS := GLOBAL_UTXOS.GetUTXOByAddrAndBalance(from , balance)
	//构建交易
	inputs := []*TXInput{}
	outputs := []*TXOutput{}
	for _, utxo := range *USER_UTXOS{
		inputs = append(inputs, &TXInput{utxo.Hash, utxo.Vout, from})
	}
	for i, to_address := range to{
		outputs = append(outputs, &TXOutput{amount[i], to_address})
	}
	//余额回转项
	dic := balance - sum
	if dic>0{
		outputs = append(outputs, &TXOutput{dic, from})
	}
	tx := NewTX(outputs, inputs)
	//将交易放入全局交易池
	*GLOBAL_TXPOOL = append(*GLOBAL_TXPOOL, tx)
	return tx;
}
func (tx TX) IsCoinBaseTX() bool{
	zero := big.NewInt(0)
	hashInt := big.NewInt(0)
	hashInt.SetBytes(tx.VIns[0].Hash[:])
	if zero.Cmp(hashInt) == 0{
		return true
	}
	return false
}