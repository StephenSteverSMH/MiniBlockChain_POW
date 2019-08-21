package BLK

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type UTXO struct {
	Hash []byte
	Vout int64
	Output *TXOutput
}
type UTXOS []UTXO
var GLOBAL_UTXOS *UTXOS = new(UTXOS)

//JSON序列化UTXO
func (utxos *UTXOS)Serialize() string{
	data, err := json.Marshal(utxos)
	if err!=nil{
		fmt.Printf("json.Marshal,err:",err);
	}
	return string(data)
}

//JSON反序列化UTXO
func (utxos *UTXOS)Deserialize(str string){
	err := json.Unmarshal([]byte(str), utxos)
	if err!=nil{
		fmt.Printf("json.Marshal,err:",err);
	}
}

//获取某用户的余额
func (utxos *UTXOS)GetBalance(address string) int64{
	var balance int64 = 0
	for _, utxo := range *utxos{
		if utxo.Output.ScriptPubKey == address{
			balance += utxo.Output.Money
		}
	}
	return balance
}

//获取某用户的UTXO
func (utxos *UTXOS)GetUTXOByAddress(address string) *UTXOS{
	var USER_UTXOS *UTXOS = new(UTXOS)
	for _, utxo := range *utxos{
		if utxo.Output.ScriptPubKey == address{
			*USER_UTXOS = append(*USER_UTXOS, utxo)
		}
	}
	return USER_UTXOS
}

//获取某用户应该拿出的输出项
func (utxos *UTXOS)GetUTXOByAddrAndBalance(address string, balance int64) *UTXOS{
	var USER_UTXOS *UTXOS = new(UTXOS)
	var sum int64 = 0
	for _, utxo := range *utxos{
		if utxo.Output.ScriptPubKey == address{
			*USER_UTXOS = append(*USER_UTXOS, utxo)
			sum+=utxo.Output.Money
			if(sum>=balance){
				return USER_UTXOS
			}
		}
	}
	if(sum<balance){
		return nil
	}
	return USER_UTXOS
}
//挖到区块后，需更新当前UTXO
func (utxos *UTXOS)UpdateUTXO(tx_pool *TXPOOL) {
	//用于比较[]byte
	//b1 := big.NewInt(0)
	//b2 := big.NewInt(0)
	//删除被用掉的utxo
	for _, tx := range *tx_pool {
		for _, in := range tx.VIns {
			//删除切片中的一个元素
			for j:=0;j<len(*utxos);j++{
				if (*utxos)[j].Output.ScriptPubKey == in.ScriptSig{
					*utxos = append((*utxos)[:j],(*utxos)[j+1:]...)
					j--
				}
			}
		}
	}
	//更新新加入的utxo
	for _, tx := range *tx_pool {
		//需要更新utxo,将这个交易的所有输出更新入UTXO
		for i, out := range tx.VOuts {
			var utxo UTXO
			utxo.Hash = tx.Hash
			utxo.Vout = int64(i)
			utxo.Output = out
			*utxos = append(*utxos, utxo)
		}
	}
}
//从数据库读入UTXO
func (utxos *UTXOS)ReadFromDB_UTXOS(db *bolt.DB){
	err:=db.View(func(tx *bolt.Tx) error{
		b:= tx.Bucket([]byte("BlockBucket"))
		if b!=nil{
			//获取最新区块
			utxos.Deserialize(string(b.Get([]byte("utxo"))))
		}
		//返回nil，以便数据库处理相应操作
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
}
//将JSON化的UTXO写入数据库
func (utxos *UTXOS)WriteToDB_UTXOS(db *bolt.DB){
	err:=db.Update(func(tx *bolt.Tx) error{
		b:= tx.Bucket([]byte("BlockBucket"))
		if b!=nil{
			b.Put([]byte("utxo"), []byte(utxos.Serialize()))
		}
		//返回nil，以便数据库处理相应操作
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
}