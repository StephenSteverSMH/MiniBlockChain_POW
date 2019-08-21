package BLK

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type Wallets map[string]*Wallet

func NewWallets() *Wallets{
	wallets := make(Wallets)
	return &wallets
}

func (wallets *Wallets)AddWallet(wallet *Wallet){
	(*wallets)[wallet.GetAddress()] = wallet
}

func (Wallet *Wallets)GetWallet(address string) (*Wallet, bool){
	val, ok := (*Wallet)[address]
	if ok == false{
		return nil, false
	}
	return val, true
}

func (wallets *Wallets)Serialize() string{
	b_wallets, err := json.Marshal(*wallets)
	if err != nil{
		log.Panic("json Marshal error")
	}
	return string(b_wallets)
}
func (wallets *Wallets)Deserialize(json_str string){
	json.Unmarshal([]byte(json_str), wallets)
}

//将JSON化的UTXO写入数据库
func (wallets *Wallets)WriteToDB(db *bolt.DB){
	err:=db.Update(func(tx *bolt.Tx) error{
		b:= tx.Bucket([]byte("BlockBucket"))
		if b!=nil{
			b.Put([]byte("wallets"), []byte(wallets.Serialize()))
		}
		//返回nil，以便数据库处理相应操作
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
}
//从数据库读入UTXO
func (wallets *Wallets)ReadFromDB(db *bolt.DB){
	err:=db.View(func(tx *bolt.Tx) error{
		b:= tx.Bucket([]byte("BlockBucket"))
		if b!=nil{
			//获取最新区块
			wallets.Deserialize(string(b.Get([]byte("wallets"))))
		}
		//返回nil，以便数据库处理相应操作
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
}
//输出钱包所有地址
func (wallets *Wallets)AddressList(){
	for key, _ := range *wallets{
		fmt.Println(key)
	}
}