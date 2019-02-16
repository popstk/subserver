# subserver
订阅服务器，自动按照配置，生成客户端用的url
- 支持v2ray配置文件

## 配置说明
```
{
    "addr": ":10086",
    "valid": {
        "SFZmIyaJDSAKgRVGf8YuiZst0": [
            {
                "type": "raw",
                "file": "config.json",
                "host": "www.example.com"
            },
            {
                "type": "v2ray",
                "file": "account.txt",
            }
        ]
    }
}
```
- addr - 订阅服务器监听地址
- valid - `json object`，key是uuid，value是一组source配置

###  source
- type: 来源类型raw | v2ray， raw代表普通文件，v2ray代表v2ray服务器端的配置文件
- file: 文件路径
- host: 仅type为v2ray有效，用于生成url的host


## TODO
- v2ray支持ws、http、quic解析

