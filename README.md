# virtflow-scheduler-go

```bash
virtflow-go/
├── go.mod
├── go.sum
├── cmd/                          # 🟢 CLI 或 main 執行點
│   ├── producer.go               # 任務發送入口（Week 1）
│   └── worker.go                 # 任務消費主程式（Week 2+）
│
├── internal/                     # 🧩 專案內部邏輯
│   ├── model/                    # 任務定義、Node 定義
│   │   └── request.go            # SchedulingRequest, BareMetalNode struct
│   ├── queue/                    # RabbitMQ 封裝
│   │   ├── publisher.go
│   │   └── consumer.go
│   ├── algorithm/                # 排程演算法（Week 3）
│   │   └── selector.go
│   ├── db/                       # SQLite Task Status 管理
│   │   └── task_store.go
│   ├── service/                  # Worker 業務邏輯（Week 5 整合點）
│   │   └── scheduler.go
│   └── util/                     # 公用工具：log, uuid, etc.
│       └── idgen.go
│
├── test/                         # 單元測試與範例（Week 9）
│   └── selector_test.go
│
├── configs/
│   └── config.yaml               # RabbitMQ URL / DB path
│
└── README.md
```