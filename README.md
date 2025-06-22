
# Scheduler-gogogo (Go Version)

Virtflow Scheduler 是一個可擴展、具備策略模式與支援 Kubernetes Leader Election 的任務排程系統，  
目標是將 VM 節點依據可用資源與策略規則，排程至最合適的 Bare Metal Server。

---

## ✨ 核心功能

- ✅ RabbitMQ 任務佇列（Producer / Consumer）
- ✅ Strategy Pattern 演算法模組
- ✅ Golang 並發處理（goroutine + context）
- ✅ Mock BareMetalNode 節點與 SchedulingRequest 任務
- ☸ 預留支援 Kubernetes Leader Election（Leader 才執行任務）
- 📦 模組化架構，方便測試與團隊協作
- 🧩 預留 ConfigMap 設定動態策略支援（backlog 中）

---

## 📁 專案架構

```
virtflow-scheduler-go/
├── cmd/
│   ├── producer.go          # 發送任務 CLI
│   └── worker.go            # 消費任務主程式
├── internal/
│   ├── algorithm/           # 排程策略模組
│   │   ├── strategy.go
│   │   ├── cpu_strategy.go
│   │   ├── memory_strategy.go   # TODO
│   │   └── hybrid_strategy.go   # TODO
│   ├── model/               # 資料模型
│   │   └── request.go
│   ├── db/                  # TODO: 任務狀態記錄
│   ├── service/             # TODO: 組合策略與處理流程
│   ├── queue/               # TODO: RabbitMQ 封裝
│   ├── elector/             # TODO: Leader Election 模組
│   └── util/                # TODO: 工具庫
├── configs/
│   └── config.yaml          # 設定檔（未來支援）
└── test/                    # 單元測試
```

---

## 🧪 SchedulingRequest 範例

```json
{
  "requested_cpu": 4,
  "requested_memory": 8192,
  "requested_pool": "default",
  "dedicated": false,
  "task_id": "abcd-1234"
}
```

---

## 🎯 Strategy Pattern 使用方式

```go
strategy := &algorithm.CPUStrategy{}
selected := strategy.SelectBestNode(request, nodes)
```

---

# virtflow-scheduler-go Backlog Overview

## 🔧 Development Backlog

| 週期     | 任務項目                                 | 狀態    |  備註                                 |
| ------ | ------------------------------------ | ----- | ---------------------------------------- |
| Week 1 | 實作 Producer CLI (發送 publish\_task)   | ✅ 已完成 | 用 RabbitMQ 發送 SchedulingRequest 任務       |
| Week 2 | 實作 Consumer CLI (執行 worker)          | ✅ 已完成 | 讀取任務並輸出 log                              |
| Week 3 | 設計 Strategy Pattern 架構               | ✅ 已完成 | 定義 SchedulingStrategy interface + CPU 策略 |
| Week 3 | 建立 MemoryStrategy, HybridStrategy 空殼 | ✅ 已完成 | TODO 傳續實作                                |
| Week 3 | consumeLoop 套用 Strategy Pattern      | ✅ 已完成 | 用 interface 呼叫可切換策略                      |
| Week 3 | 撰寫系統 README + 架構圖                    | ✅ 已完成 | markdown 包含條狠 clear                      |
| Week 4 | 實作 MemoryStrategy + HybridStrategy   | ⏳ 進行中 | 據 memory/混合比量 算分邏輯                       |
| Week 4 | 換成 PostgreSQL (替代 SQLite)            | ✅ 已完成 | internal/db + task\_status 表             |
| Week 5 | Leader Election 機制                 | ⏳ 規劃中 | 會用 k8s.io/client-go                      |
| Week 5 | Worker 必須為 leader 才啟用                | ⏳ 規劃中 | idle follower, active leader             |
| Week 6 | 支援 ENV / ConfigMap 切換 Strategy       | ⏳ 規劃中 | SCHEDULING\_STRATEGY=xxx                 |
| Week 6 | 撰寫 Strategy Factory (GetStrategy)    | ⏳ 規劃中 | string -> struct 映射器                     |
| Week 7 | 發送 Webhook (成功/失敗)                   | ⏳ 規劃中 | task callback HTTP POST                  |
| Week 8 | 多 Worker + Graceful Shutdown         | ⏳ 規劃中 | goroutine + context 關閉                   |
| Week 9 | 撰寫單元測試                               | ⏳ 規劃中 | algorithm / DB / worker                  |

---

## 💡 Concurrency Safety / Performance Enhancement Backlog

| ID | 項目                                | 狀態    | 備註                      |
| -- | --------------------------------- | ----- | ----------------------- |
| C1 | 加入處理耗時經過記錄                        | ✅ 已完成 | time.Since() + log      |
| C2 | goroutine 內 recover 保護機制          | ⏳ 待開發 | 避免 panic 擊敗全體 worker    |
| C3 | RabbitMQ channel 加入 prefetchCount | ⏳ 待開發 | 限制協約同時處理數               |
| C4 | 分粘 publishTask 錯誤類型               | ✅ 已完成 | timeout 與網路錯誤分開         |
| C5 | publishTask 加入 retry (3 次)        | ⏳ 待開發 | 用 for loop + backoff 重試 |

---


## 🧭 Golang Developer 能力升級表（Weekly Path）

| 週數 | 能力主題                   | 學習重點                                       | 任務對應 & 建議練習                         | 解鎖狀態 |
|------|----------------------------|------------------------------------------------|----------------------------------------------|----------|
| Week 1 | Golang 基本語法 + struct/json | struct 定義、json tag、marshal/unmarshal     | 完成 `SchedulingRequest`, 發送任務 payload   | ✅       |
| Week 2 | goroutine / channel       | `go func()`、非同步處理、`select{}` 控制 loop | `worker.go` 使用 goroutine + select 消費任務 | ✅       |
| Week 3 | interface + Strategy Pattern | interface 定義、結構抽象、演算法封裝         | 實作 `SchedulingStrategy` + `CPUStrategy`    | ✅       |
| Week 4 | 排程邏輯設計與排序技巧     | slice 過濾、`sort.Slice()`、條件判斷邏輯      | 自主實作 `MemoryStrategy`、排序反轉邏輯      | ⏳       |
| Week 5 | SQLite 資料庫操作         | `database/sql`、prepare/exec、error check     | 更新 `task_status`、記錄 success/failed     | ⏳       |
| Week 6 | Golang config & factory   | `os.Getenv()`、封裝 Factory 模式              | 根據 env 決定策略 → `GetStrategy(name)`     | ⏳       |
| Week 7 | Leader Election 與 context cancel | K8s Leader 模組、context 控制中止           | 整合 `elector.go` → Leader 才能 run worker   | ⏳       |
| Week 8 | 錯誤處理與 retry 模式     | 分類 error、重試設計、如何用 log 分層         | Retry 任務策略 + 可調 backoff（optional）   | ⏳       |
| Week 9 | 單元測試與 table-driven test | `testing` package、範例表格測試設計         | 為策略模組寫測試 + edge case 模擬           | ⏳       |
| Week 10 | 重構與模組設計總整理     | package 拆分、依賴倒轉、邊界設計與文件結構     | 將 worker 拆成 service 層，撰寫總結 README   | ⏳       |
