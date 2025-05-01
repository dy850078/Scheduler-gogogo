
# 📅 明日進公司任務規劃

目標：完成 PostgreSQL 任務儲存整合 + 測試排程 API Server

---

## ✅ 上午任務：啟動與驗證環境

| 項目 | 預估時間 | 工具 | 備註 |
|------|----------|------|------|
| 啟動 PostgreSQL（本地或 VM）| 20 min | docker / systemd | 確保可連線 |
| 匯入 schema (`task_status`) | 10 min | psql / GUI 工具 | `InitSchema()` 會自動執行 |
| 設定 `POSTGRES_CONN` 環境變數 | 5 min | `.env` / bashrc | ✅ 必備 |
| 執行 `worker.go` 檢查能寫入任務結果 | 15 min | `go run cmd/worker.go` | 先從手動送出任務開始（Producer or curl） |

---

## ✅ 下午任務：API Server 實作與串接測試

| 項目 | 預估時間 | 工具 | 備註 |
|------|----------|------|------|
| 解壓 & 整合 API 套件 | 10 min | VS Code | 可與主專案合併 |
| 設定 `RABBITMQ_URL` 並啟動 API | 10 min | `go run cmd/api_server.go` | 埠為 8080 |
| 測試 `/schedule` API 是否能送出任務 | 10 min | curl / Postman | `Accepted` 回應即成功 |
| 在 worker 查看是否能接收到任務 | 10 min | log 應看到 TaskID 與排程結果 | ✅ 串接成功 |

---

## ⏱️ 額外 Bonus 任務（可選）

| 項目 | 預估時間 | 建議等級 | 備註 |
|------|----------|----------|------|
| 自行撰寫 `MemoryStrategy` | 30–60 min | 🔥實戰練習 | 可先手寫邏輯、回家後再讓我 review |
| 用 DB GUI 工具檢查寫入資料 | 10 min | 建議使用 pgAdmin / DBeaver | 觀察 task 狀態更新結果 |
| 設定 log 格式（加上時間、TaskID）| 15 min | 好維運習慣 | 日後好用於 tracing |

---

## 📦 預期完成產出：

- [ ] PostgreSQL 初始化與連線驗證
- [ ] API Server 可正常 POST `/schedule`
- [ ] 任務送出 → Worker consume → 正確寫入 task_status
- [ ] 若有餘裕，完成 `MemoryStrategy` 初版實作




# 🛠 virtflow-scheduler-go - 空機器環境建置與開發啟動指南

此文件用於指導在尚未安裝 Golang、PostgreSQL、RabbitMQ 的 Linux 開發環境中，快速部署與啟動 virtflow-scheduler-go。

---

## ✅ 環境建置（建議作業系統：Ubuntu 20.04+）

### 0. 系統更新與工具安裝

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install curl wget unzip git -y
```

### 1. 安裝 Golang（使用 Snap 或手動）

```bash
sudo snap install go --classic
go version  # 確認版本
```

或手動安裝最新版（建議 Go 1.20 以上）

---

### 2. 安裝 PostgreSQL

```bash
sudo apt install postgresql postgresql-contrib -y
sudo systemctl enable postgresql
sudo systemctl start postgresql
```

建立資料庫與使用者（預設 user/password 可根據程式碼設定）

```bash
sudo -u postgres psql
CREATE USER virtflow WITH PASSWORD 'password';
CREATE DATABASE virtflow OWNER virtflow;
\q
```

---

### 3. 安裝 RabbitMQ

```bash
sudo apt install rabbitmq-server -y
sudo systemctl enable rabbitmq-server
sudo systemctl start rabbitmq-server
```

檢查是否成功啟動：

```bash
sudo rabbitmqctl status
```

---

### 4. 環境變數設定

請在 shell 或 `.bashrc` 中加入：

```bash
export POSTGRES_CONN="postgres://virtflow:password@localhost:5432/virtflow?sslmode=disable"
export RABBITMQ_URL="amqp://guest:guest@localhost:5672/"
```

---

## 🧪 驗證工作任務清單（明日開發項目）

### 1️⃣ 啟動 worker + 檢查 DB 連線與建表

```bash
go run cmd/worker.go
```

- ☑ 應顯示 `[INFO] Worker running...`
- ☑ 會自動建立 `task_status` 資料表（無錯誤即成功）

---

### 2️⃣ 啟動 API Server

```bash
go run cmd/api_server.go
```

- ☑ 埠號為 `:8080`
- ☑ 日誌應顯示 `[INFO] API Server running on :8080`

---

### 3️⃣ 測試 `/schedule` 任務流程

使用 curl：

```bash
curl -X POST http://localhost:8080/schedule   -H "Content-Type: application/json"   -d '{
    "requested_cpu": 4,
    "requested_memory": 8192,
    "requested_pool": "default",
    "dedicated": false,
    "task_id": "demo-001"
  }'
```

- ☑ 回應為 `202 Accepted`
- ☑ `worker.go` log 顯示接收與處理
- ☑ `psql` 查詢資料：

```bash
psql virtflow virtflow
SELECT * FROM task_status;
```

---

## ⛳ 建議後續任務

| 任務 | 說明 |
|------|------|
| 完成 `MemoryStrategy` 實作 | 練習你對策略模式與資源條件邏輯的理解 |
| 編寫 DB GUI 工具設定（pgAdmin） | 方便日後查表與管理 |
| 記錄 README + 今日操作筆記 | 留下設定與遇到的問題方便回報與日後查閱 |

---

📘 如需手動建立 queue，可執行：

```bash
rabbitmqctl list_queues
```

大部分 queue 會由程式執行時自動建立。
