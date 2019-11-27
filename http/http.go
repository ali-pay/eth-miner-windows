package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"ieth/conf"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
)

type rpc struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
	Result  string   `json:"result"`
}

//获取账户金额
func GetMoney(account string) (money string, err error) {
	req := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBalance","params":["%s", "latest"],"id":666}`, account)
	//util.Info.Println("req: ", req)

	httpClient := new(http.Client)
	resp, err := httpClient.Post("http://127.0.0.1:8545", "application/json", strings.NewReader(req))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	//util.Info.Println("resp: ", string(body))

	//解析应答数据
	res := new(rpc)
	if err = json.Unmarshal(body, res); err != nil {
		return
	}
	//util.Info.Printf("JsonRPC: %+v\n", res)

	//16进制转换10进制
	n := new(big.Int)
	n2, b := n.SetString(res.Result[2:], 16)
	if !b {
		err = errors.New("余额转换失败")
		return
	}

	money = n2.String()
	return
}

//获取区块数
func GetBlock() (block string, err error) {
	req := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":666}`
	//util.Info.Println("req: ", req)

	httpClient := new(http.Client)
	resp, err := httpClient.Post("http://127.0.0.1:8545", "application/json", strings.NewReader(req))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	//util.Info.Println("resp: ", string(body))

	//解析应答数据
	res := new(rpc)
	if err = json.Unmarshal(body, res); err != nil {
		return
	}
	//util.Info.Printf("JsonRPC: %+v\n", res)

	//16进制转换10进制
	n := new(big.Int)
	n2, b := n.SetString(res.Result[2:], 16)
	if !b {
		err = errors.New("区块数转换失败")
		return
	}

	block = n2.String()
	return
}

//挖矿
func Miner(state bool) (str string, err error) {
	//停止挖矿
	req := `{"jsonrpc":"2.0","method":"miner_stop","params":[],"id":666}`

	//开始挖矿
	if state {
		req = fmt.Sprintf(`{"jsonrpc":"2.0","method":"miner_start","params":[%s],"id":666}`,conf.Minerthreads)
	}

	httpClient := new(http.Client)
	resp, err := httpClient.Post("http://127.0.0.1:8545", "application/json", strings.NewReader(req))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	str = string(body)
	if len(str) != 41 {
		err = errors.New(str)
		return
	}
	return
}

//退出时，本地保存区块数据 需要在geth中加入退出的接口
func Exit() (str string, err error) {
	req := `{"jsonrpc":"2.0","method":"eth_exit","params":[],"id":666}`

	httpClient := new(http.Client)
	resp, err := httpClient.Post("http://127.0.0.1:8545", "application/json", strings.NewReader(req))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	str = string(body)
	if len(str) != 41 {
		err = errors.New(str)
		return
	}
	return
}

func GetPeer() (peer string, err error) {
	req := `{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":666}`
	//util.Info.Println("req: ", req)

	httpClient := new(http.Client)
	resp, err := httpClient.Post("http://127.0.0.1:8545", "application/json", strings.NewReader(req))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	//util.Info.Println("resp: ", string(body))

	//解析应答数据
	res := new(rpc)
	if err = json.Unmarshal(body, res); err != nil {
		return
	}
	//util.Info.Printf("JsonRPC: %+v\n", res)

	//16进制转换10进制
	n := new(big.Int)
	n2, b := n.SetString(res.Result[2:], 16)
	if !b {
		err = errors.New("节点数转换失败")
		return
	}

	peer = n2.String()
	return
}
