# 以太坊单机挖矿客户端

## 编译说明

- 生成syso文件，执行一次就行

```bash
rsrc -arch="amd64" -ico ICON.ico -manifest main.manifest -o main.syso
```

- 开启黑窗输出调试信息

```go
go build
```

- 关闭黑窗输出

```go
go build -ldflags="-H windowsgui -w -s"
```

- 错误处理: 注释掉go-ethereum@v1.9.6\ethdb\leveldb\leveldb.go:103:3: DisableSeeksCompaction

- 编译参考: https://blog.csdn.net/kgjn__/article/details/89280235



## 单机挖矿使用说明

### 生成账号（钱包）

输入密码，点击生成账号即可，牢记密码以及私钥，密钥文件保存着data文件夹下的keystore文件夹中。

### CPU挖矿

- 有了账号后，在钱包地址栏中填入账号，点击CPU挖矿即可

### 显卡挖矿

- 同CPU挖矿操作相同

### 打开控制台

- 打开官方的GETH控制台，这是一个js环境，在此处可以进行web3的操作

### 停止挖矿

- 停掉挖矿中的进程

### 配置文件说明

```js
//此处是软件启动时加载的钱包地址，当点击挖矿的时候，会自动将钱包地址栏中的地址替换到此处
钱包地址=0xeb22459524804361ab700f2552b066a1392b80ab

//网络的ID，必须跟其他节点一样才能进行通信
网络ID=100

//本机暴露的端口，跟其他节点通信
端口=61910

//挖矿的线程，默认是根据CPU核心数的，此处设置为1即可
挖矿线程=1

//区块的同步方式，full是全同步，fast是快速同步，如果fast同步失败，则在此改成使用full进行同步
区块同步=full

//节点启动时去连接的另一个节点，当然也可以打开控制台，使用admin.addPeer()函数进行添加，添加后下次会自动尝试连接
启动节点=enode://5e00f29f43637107067e5c0cf94004fd8f9475d4d4f86be7501acb20f54a3ad3d2d70ec5c660dbb86656560a5e19892bf1cf38c65ae8202e54a2eb7fe144bbb1@192.168.105.215:61910
```

### 数据保存说明

- 默认数据保存在data文件夹中

- geth：以太坊的链数据，保存区块的，如果遇到区块同步出错的情况，建议删除这个文件夹，重新同步区块

- keystore：保存密钥文件，请牢记密码以及备份好这个文件

