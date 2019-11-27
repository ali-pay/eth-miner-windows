package acc

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/mobile"
	"io/ioutil"
)

//传入密码生成密钥文件，返回文件存放位置
func NewAccount(password string) (path string, err error) {
	//N:1<<12 P:6 易
	//N:1<<18 P:1 难
	ks := geth.NewKeyStore("./data/keystore", 1<<18, 1)
	acc, err := ks.NewAccount(password)
	if err != nil {
		return
	}
	path = acc.GetURL()
	path=path[11:]
	return
}

//传入密钥文件和密码，生成地址和私钥明文
func DecryptKeystore(file, password string) (address string, privateKey string, err error) {
	keyjson, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		return
	}
	address = key.Address.Hex()
	privateKey = hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))
	return
}
