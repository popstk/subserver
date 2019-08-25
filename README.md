# subserver [![Go Report Card](https://goreportcard.com/badge/github.com/popstk/subserver)](https://goreportcard.com/report/github.com/popstk/subserver)
订阅服务器，自动按照配置，生成客户端用的url
- 支持v2ray配置文件

## 配置说明
```
{
    "addr": ":10086",
    "valid": {
        "SFZmIyaJDSAKgRVGf8YuiZst0": [
            {
                "type": "v2ray",
                "file": "config.json",
                "host": "www.example.com",
                "vmess-fmt": "{protocol}-{network}",
                "ss-fmt": "{protocol}"
            },
            {
                "type": "sub",
                "addr": "dnH4PH3I6ufsxFr1Bd3Ghryi"
            }
        ],
        "dnH4PH3I6ufsxFr1Bd3Ghryi": [
            {
                "type": "raw",
                "file": "account.json"
            },
            {
                "type": "url",
                "addr": "http://path/to/other/subserver"
            }
        ]
    }
}
```
- addr - 订阅服务器监听地址
- valid - `json object`，key是uuid，value是一组source配置

###  source
type 代表数据来源，以下是可选的类型和相应的字段:

- raw
  - `file`指定以换行符分割的url地址文本文件

- v2ray 
  - ` file` 指定 v2ray服务器端的配置文件
  - `host`指定生成url指定的`remote`地址
  - `vmess-fmt`生成vmess配置的备注格式
  - `ss-fmt`生成ss配置的备注格式

- sub
    - `addr`合并其他`uuid`的url

- url
    - `addr`合并其他`subserver`的url
 
## 特性
- 支持从v2ray配置文件，根据suberver配置，生成节点
- 自动过滤localhost节点
- 自动更新配置文件

