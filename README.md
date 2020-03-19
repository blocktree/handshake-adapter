# handshake-adapter

handshake-adapter适配了openwallet.AssetsAdapter接口，给应用提供了底层的区块链协议支持。

## 项目依赖库

- [go-owcrypt](https://github.com/blocktree/go-owcrypt.git)
- [go-owcdrivers](https://github.com/blocktree/.git)

## 如何测试

openwtester包下的测试用例已经集成了openwallet钱包体系，创建conf文件，新建HNS.ini文件，编辑如下内容：

```ini


# node api url, if RPC Server Type = 1, use bitbay insight-api
nodeAPI = "http://ip:port"
# RPC Authentication Username
rpcUser = ""
# RPC Authentication Password
rpcPassword = ""
# minimum transaction fees
minFeeRate = "0.0001"
# Cache data file directory, default = "", current directory: ./data
dataDir = ""

```
