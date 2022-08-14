簡易下單應用程式。

通過 Google [Cloud Functions](https://cloud.google.com/functions?hl=zh-tw) 在雲端環境中執行程式碼，無需管理伺服器

函式即服務 (FaaS)，等待訊號 webhook 來送出交易

收到「多方訊號！」就開倉，「平倉訊號！」就平倉

多、減、平、停

## 準備工作

註冊一個 Google Cloud 項目

綁定信用卡，啟用結算功能

### Cloud Functions 費用

128MB, asia-east1 區域, 一天交易50次

一個月大約 30 美金

## 發佈程式

整個流程設定耗時 ~7分鐘

### 建立函式

在建立好的 Google Cloud 項目中，點擊 [建立函式](https://console.cloud.google.com/functions/add?hl=zh-tw) 按鈕

### 設定

「函式名稱」輸入 `simple-trade-machine`

「地區」建議別選「美洲」，因為有些 token 是不能在美國的機器上交易的。

- asia-east1 台灣
- asia-east2 香港
- asia-southeast1 新加坡

「觸發條件」選 `HTTP`，並勾選 `允許未經驗證的叫用`

「執行階段」分配的記憶體，選 128 MB

「自動調度資源」執行個體數量下限， 填 2

「執行階段環境變數」新增變數，參考 [環境變數](#環境變數)

### 程式碼

下一步到 `2.程式碼` 部分

「執行階段」選 `Go1.16`

「進入點」鍵盤輸入 `AlertHandler` **必填**

「原始碼」選 `上傳 ZIP 檔案` ，把目錄帶有 `go.mod`的整個目錄壓縮成 zip
上傳，或是下載 [原始碼zip](https://github.com/snakehopper/simple-trade-machine/archive/refs/heads/master.zip)

「暫存值區」點擊「瀏覽」選取其中一個

點擊「部署」按鈕，等待雲端部署

### 部署完成

點擊函式名稱查看 `網址` (需保密)

格式大概為 https:// simple-trade-machine-blablabla-de.a.run.app

或 https: //asia-east1-my-blabla-project.cloudfunctions.net/simple-trade-machine

這網址就是 webhook url ，對其做POST 會觸發開倉

`https://<cloud-function-host>/<交易所>/<交易對>`

例如 FTX 的 BTC-PERP 收到多方訊號，則 POST url 如下：

https:// simple-trade-machine-blablabla.a.run.app/**ftx/BTC-PERP**

## 風險

webhook url `網址` 洩漏的話，別人就能控制機器下單！

Google Cloud functions 服務中斷的話，將會錯過開倉、平倉訊號！

## 環境變數

| 名稱             | 值                         |
|----------------|---------------------------|
| OPEN_PERCENT   | 每次開倉比例：10 代表每次使用十份之一的資金下單 |
| REDUCE_PERCENT | 每次減倉比例：50 代表每次減倉一半的倉位     |
| FTX_APIKEY     | FTX 網頁申請的 API Key         |   
| FTX_SECRET     | FTX 網頁申請的 API Secret      |   
| FTX_SUBACCOUNT | (選填) main 帳號填空 `""`       |   
| BINANCE_APIKEY | 幣安網頁申請的 API Key           |   
| BINANCE_SECRET | 幣安網頁申請的 API Secret        |   

## 限制

### 幣安

不支援月合約、季合約
如果要交易 BUSD 交易對，需要開啟「聯合保證金模式」Multi-Assets Mode