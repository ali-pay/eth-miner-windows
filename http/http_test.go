package http

import (
	"fmt"
	"testing"
)

//测试获取金额
func TestGetMoney(T *testing.T) {
	//req:  {"jsonrpc":"2.0","method":"eth_getBalance","params":["0xeb22459524804361ab700f2552b066a1392b80ab", "latest"],"id":666}
	//resp:  {"jsonrpc":"2.0","id":666,"result":"0x38e78d476008a66000"}
	money, err := GetMoney("0xeb22459524804361ab700f2552b066a1392b80ab")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mlen := len(money)
	if mlen > 18 {
		point := mlen - 18
		s := money[:point]
		e := money[point:]
		money = s + "." + e
	}
	fmt.Println("money: ", money) //money:  1049.702738800000000000
}

//测试获取区块
func TestGetBlock(T *testing.T) {
	//req:  {"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":666}
	//resp:  {"jsonrpc":"2.0","id":666,"result":"0x67fb"}
	block, err := GetBlock()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("block: ", block) //block:  26619
}

//测试http挖矿
func TestStartMiner(t *testing.T) {
	resp, err := Miner(false)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(len(resp))
	fmt.Println("resp: ", resp)

	//resp:  {"jsonrpc":"2.0","id":666,"result":null} true
	//resp:  {"jsonrpc":"2.0","id":666,"result":null} false
}

//测试退出
func TestExit(t *testing.T) {
	resp, err := Exit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(len(resp))
	fmt.Println("resp: ", resp) //resp:  {"jsonrpc":"2.0","id":666,"result":null}
}

func TestGetPeer(t *testing.T) {
	resp, err := GetPeer()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("peer: ", resp) //peer:  0
}
