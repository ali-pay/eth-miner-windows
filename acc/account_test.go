package acc

import (
	"fmt"
	"testing"
)

var ps = "123456"
var ks = `D:\gopath\src\ieth\acc\data\keystore\UTC--2019-11-12T08-09-12.455021200Z--194d0a27980b969f43dc3e73ab0851a6ffabb767`

//创建账户
func TestNewAccount(t *testing.T) {
	path, err := NewAccount(ps)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Path:",path)
	//Path: D:\gopath\src\ieth\acc\data\keystore\UTC--2019-11-12T08-09-12.455021200Z--194d0a27980b969f43dc3e73ab0851a6ffabb767
}

//解析地址与私钥
func TestDecryptKeystore(t *testing.T) {
	address, privateKey, err := DecryptKeystore(ks, ps)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Address: %s\r\nPrivateKey: %s\r\n",address, privateKey)
	//Address: 0x194d0A27980b969f43dc3E73Ab0851A6fFabb767
	//PrivateKey: e34512b7c4f533f1dc65e8e369d4aeb79fc77cccbc29ce60e75b9bfbec9cecd4
}
