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

#### 基本
「環境」選 `第2代`

「函式名稱」輸入 `simple-trade-machine`

「地區」建議別選「美洲」，因為有些 token 是不能在美國的機器上交易的。

- asia-east1 台灣
- asia-east2 香港
- asia-southeast1 新加坡

#### 觸發條件
1. 「觸發條件類型」選 `HTTP`
   _（留意這裡的`網址`選項，後面會用到）_
2. 勾選 `允許未經驗證的叫用`
3. 儲存  

#### 執行階段、建構作業、連線和安全性設定
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

### webhook範例

| 範例    | 交易所             | 交易對      | webhook                                                                    |
|-------|-----------------|----------|----------------------------------------------------------------------------|
| 幣安現貨  | binance         | BTCUSDT  | `https://simple-trade-machine-blablabla.a.run.app/binance/BTCUSDT`         |
| 幣安合約  | binance-futures | BTCUSDT  | `https://simple-trade-machine-blablabla.a.run.app/binance-futures/BTCUSDT` |
| FTX現貨 | ftx             | ETH/USDT | `https://simple-trade-machine-blablabla.a.run.app/ftx/ETH/USDT`            |
| FTX合約 | ftx             | ETH-PERP | `https://simple-trade-machine-blablabla.a.run.app/ftx/ETH-PERP`            |

_(此處 blablabla 應該會是[建立函式](#建立函式)後 Google 給的網址)_


## 風險

webhook url `網址` 洩漏的話，別人就能控制機器下單！

Google Cloud functions 服務中斷的話，將會錯過開倉、平倉訊號！

## 環境變數

| 名稱             | 值                                    |
|----------------|--------------------------------------|
| OPEN_PERCENT   | ! 每次開倉比例：10 代表每次使用十份之一的資金下單          |
| REDUCE_PERCENT | ! 每次減倉比例：50 代表每次減倉一半的倉位              |
| SPOT_OPEN_X    | 現貨的開倉量＝可用資金＊OPEN_PERCENT＊SPOT_OPEN_X |
| ORDER_TYPE     | ! 使用限價單或市價開單，`limit` `market`        |
| FOLLOWUP_LIMIT_ORDER     | 限價單多久沒成交改市價單 例：50s, 1m, 2m30s        |
| FTX_APIKEY     | FTX 網頁申請的 API Key                    |   
| FTX_SECRET     | FTX 網頁申請的 API Secret                 |   
| FTX_SUBACCOUNT | (選填) main 帳號填空 `""`                  |   
| BINANCE_APIKEY | 幣安網頁申請的 API Key                      |   
| BINANCE_SECRET | 幣安網頁申請的 API Secret                   |   

! 表示可根據策略客製化，例如一般策略都是減倉10%，唯左側拐點多方減倉訊號每次30%" 寫作 `REDUCE_PERCENT: "10"`  `COUNTER_REDUCE_PERCENT: "30"`
 
## 限制

### 幣安

不支援月合約、季合約
如果要交易 BUSD 交易對，需要開啟「聯合保證金模式」Multi-Assets Mode

## FAQ

```
設定流程好多步驟，有沒有簡單點的方法？
```
有。參考 [命令行部署](GCLOUD.md)

```
現貨交易對，碰到做空訊號如何處理？
```
程式碰到「多轉空訊號」「空方訊號」只會平掉手上的現貨倉位，並順利結束；不會有實際做空的動作

```
如何設定槓桿？
```
槓桿無需在程式裡配置，直接到交易所的網頁裡設置即可。

系統收到指標訊號後，會根據 「帳號可用資金 x OPEN_PERCENT x 帳號槓桿」 去開倉下單

現貨則採「帳號可用資金 x OPEN_PERCENT x SPOT_OPEN_X」

```
可以同時多個交易所下單嗎？
```
支援。通過新增多個「快訊」配置不同交易所的 webhook 即可。 [參考webhook範例](#webhook)

