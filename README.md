
# Virtflow Scheduler (Go Version)

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

## 🛠 Backlog（開發任務總覽）

| 週期 | 任務項目                                       | 狀態     | 備註 / 說明                                  |
|------|------------------------------------------------|----------|-----------------------------------------------|
| Week 1 | 實作 Producer CLI（publish_task）             | ✅ 已完成 | 使用 RabbitMQ 發送 SchedulingRequest 任務    |
| Week 2 | 實作 Consumer CLI（worker 基礎版）             | ✅ 已完成 | Consume 任務，並印出 SchedulingRequest        |
| Week 3 | 設計 Strategy Pattern 架構                    | ✅ 已完成 | 定義 interface 與實作 `CPUStrategy`          |
| Week 3 | 建立空殼 `MemoryStrategy`, `HybridStrategy`   | ✅ 已完成 | 作為後續可插入演算法預留點                   |
| Week 3 | consumeLoop 中套用 strategy pattern          | ✅ 已完成 | 用 interface 呼叫可替換策略                   |
| Week 3 | 撰寫系統架構與專案 README                    | ✅ 已完成 | 包含功能總覽、檔案架構、JSON 與使用方式     |
| Week 4 | 撰寫 `MemoryStrategy`, `HybridStrategy`       | ⏳ 進行中 | CPU 以外邏輯選擇節點（記憶體、混合評分）     |
| Week 4 | 實作 SQLite 任務狀態儲存模組（成功 / 失敗）  | ⏳ 進行中 | 建立 `task_status` table，更新 task 狀態     |
| Week 5 | 封裝 Leader Election 模組                     | ⏳ 規劃中 | 使用 `k8s.io/client-go/tools/leaderelection` |
| Week 5 | Worker 啟動前需確認為 Leader 才進行任務處理   | ⏳ 規劃中 | 非 leader 則 idle（不可 consume）            |
| Week 6 | 透過 ConfigMap 或 ENV 選擇 strategy           | ⏳ 規劃中 | 支援 `SCHEDULING_STRATEGY=cpu` 類型控制       |
| Week 6 | 撰寫 Strategy Factory (`GetStrategy`)         | ⏳ 規劃中 | 將 string 對應到實體策略 struct              |
| Week 7 | 發送 Webhook（任務 success / failed）         | ⏳ 規劃中 | 用 HTTP POST 將結果通知外部系統              |
| Week 8 | 多 worker 支援 + Graceful shutdown           | ⏳ 規劃中 | 多 goroutine + context cancel 控制結束        |
| Week 9 | 撰寫單元測試（演算法 + Consumer + DB）       | ⏳ 規劃中 | 使用 table-driven test，提升穩定性            |


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