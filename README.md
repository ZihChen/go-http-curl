# go-http-curl

#### 校驗 Http ConnectPool
**啟動3個Web Server:**
```bash
cd benchmark
go run main.go -port 8087
go run main.go -port 8088
go run main.go -port 8089
```
**Run Test:**

驗證ConnectPool在多線程請求結果是否正確
1. 各個Server端取得的query params及headers需與線程請求帶的參數一致
2. 各個Server端顯示的remote address的port號是否都來自同一個
3. 每個線程只會對Server建立一次連線，Server關閉後顯示的reconnect count需等於0
```bash
cd benchmark/service
go test -v
```