# simple short url

这是一个面向开发者个人的简单实现的短链服务，预期短链数量规模较少，不会产生大的压力，使用 hertz http 框架实现。

短链示意：[https://s.mcoder.cc/s/home](https://s.mcoder.cc/s/home)

## 原理

1. 在内存中有个全局 map，存储全部短链；
2. 使用 json 文件持久化数据

## 能力

1. 支持添加短链、删除、查询短链
2. 支持通过短链 307 临时跳转到短链地址
3. 支持设置一个 token，仅有开发者自己知道的 token，可以避免恶意请求

## 预期部署方式

### 编译

1. 需要安装 golang 环境，并确保 golang 版本大于 1.21
2. 运行 `sh ./build.sh` 编译服务，编译完成的内容位于 `./output` 目录下

### 运行时的一些说明

1. 本服务实现时仅绑定了 HTTP 端口，建议开发者使用 nginx 反向代理将 HTTP 包装成 HTTPS。
关于反向代理的配置可以自己查阅资料，这里不展开。
2. 请注意看 `./script/bootstrap.sh` 中的一些环境变量，根据自己的实际情况进行修改。
3. 运行时配置请注意看 `./conf/config_demo.yaml`，需要将文件名称修改为 `confg.yaml` 方可使用
4. 短链的存储文件位于 `./conf/short.json`，需要在程序启动前将文件创建好，最简单的方法就是修改
`./conf/short_demo.json`名称为 `short.json`

### 接口

接口参考 `idl/service.thrift` 文件。

```cgo
    RedirectShortURLResponse RedirectShortURL(1: RedirectShortURLRequest request)(api.get="/s/:url");
    AddShortURLResponse AddShortURL(1: AddShortURLRequest request)(api.post="/s/short_url");
    DeleteShortURLResponse DeleteShortURL(1: DeleteShortURLRequest request)(api.delete="/s/:url");
    RefreshResponse Refresh(1: RefreshRequest request)(api.post="/s/refresh");
    ListShortURLResponse ListShortURL(1: ListShortURLRequest request)(api.get="/s/list");
```
