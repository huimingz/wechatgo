# wechatgo
企业微信开发工具Go语言版本

## 安装

```
go get -v github.com/huimingz/wechatgo
```

使用`vendor`:
```
govendor fetch github.com/huimingz/wechatgo
```

## 文档
文档地址：[https://godoc.org/github.com/huimingz/wechatgo](https://godoc.org/github.com/huimingz/wechatgo)

## 特性

- [x] 通讯录管理（大部分）
- [x] 外部联系人管理
- [x] 身份验证 
- [x] 应用管理 
- [x] 消息推送（部分）
- [x] 素材管理
- [x] OA数据接口
- [x] 电子发票


## 关于测试
少部分功能由于缺乏测试环境，未进行测试。

```
ok  	github.com/huimingz/wechatgo/client	1.102s	coverage: 40.5% of statements
?   	github.com/huimingz/wechatgo/entapi	[no test files]
ok  	github.com/huimingz/wechatgo/entapi/app	1.457s	coverage: 90.3% of statements
ok  	github.com/huimingz/wechatgo/entapi/dept	2.034s	coverage: 100.0% of statements
ok  	github.com/huimingz/wechatgo/entapi/extcontact	0.857s	coverage: 11.9% of statements
?   	github.com/huimingz/wechatgo/entapi/invoice	[no test files]
ok  	github.com/huimingz/wechatgo/entapi/media	1.716s	coverage: 42.0% of statements
ok  	github.com/huimingz/wechatgo/entapi/msg	1.041s	coverage: 33.8% of statements
ok  	github.com/huimingz/wechatgo/entapi/oa	1.459s	coverage: 100.0% of statements
ok  	github.com/huimingz/wechatgo/entapi/oauth	0.688s	coverage: 93.9% of statements
ok  	github.com/huimingz/wechatgo/entapi/tag	1.642s	coverage: 100.0% of statements
ok  	github.com/huimingz/wechatgo/entapi/user	2.975s	coverage: 97.1% of statements
ok  	github.com/huimingz/wechatgo/log	(cached)	coverage: 100.0% of statements
ok  	github.com/huimingz/wechatgo/storage	0.008s	coverage: 95.7% of statements
```

## 可能存在的问题
企业微信文档中"获取公费电话拨打记录"的示例中存在错误，加之无测试环境，无法作出有效判断。

## 版权
使用MIT许可证授权，详细内容查看[LICENSE](LICENSE)文件。
