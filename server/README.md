下载phpstudy 打开Mysql 
go run ./cmd/main.go
go mod tidy
go clean -modcache


$env:GOOS="linux"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0"; go build -o app ./cmd/main.go   //打包 放上宝塔 
运行 ./1.sh

## 支付接入

易支付环境变量、模拟回调、真实支付测试和上线检查见 [PAYMENT_TEST.md](./PAYMENT_TEST.md)。
