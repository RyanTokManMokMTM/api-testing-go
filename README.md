# API 測試框架 🚀

一個強大且靈活的 API 測試框架，讓您能夠使用 YAML 配置來定義、執行和維護 API 測試工作流程。本框架旨在使 API 測試更高效、可維護且可擴展。

## ✨ 功能特點

- **基於 YAML 的配置**：使用簡單的 YAML 語法定義測試工作流程
- **動態變量支持**：使用預定義和從響應中提取的變量
- **全面的驗證**：多種斷言類型，實現全面的 API 測試
- **靈活的鉤子系統**：全局和場景特定的設置和清理鉤子
- **可擴展架構**：輕鬆添加新路由和測試場景
- **基於 Go 的實現**：高性能和可靠性

## 📋 目錄

- [開始使用](#開始使用)
  - [前置條件](#前置條件)
  - [安裝](#安裝)
  - [快速開始](#快速開始)
- [核心概念](#核心概念)
  - [API 路由配置](#api-路由配置)
  - [測試工作流程結構](#測試工作流程結構)
  - [變量和動態值](#變量和動態值)
  - [鉤子系統](#鉤子系統)
- [編寫測試](#編寫測試)
  - [基本測試結構](#基本測試結構)
  - [請求配置](#請求配置)
  - [響應驗證](#響應驗證)
  - [高級功能](#高級功能)
- [最佳實踐](#最佳實踐)
- [貢獻指南](#貢獻指南)
- [支持](#支持)

## 🚀 開始使用

### 前置條件

- Go 1.x 或更高版本
- REST API 基礎知識
- YAML 語法基礎
- Git 版本控制

### 安裝

1. 克隆倉庫：
```bash
git clone https://github.com/your-org/github.com/RyanTokManMokMTM/api-testing-go.git
cd github.com/RyanTokManMokMTM/api-testing-go
```

2. 安裝依賴：
```bash
go mod download
```

### 快速開始

1. 進入項目目錄
2. 運行測試套件：
```bash
go test ./test/...
```

3. 查看 `/config/etc/api-test/` 中的示例測試工作流程以了解結構

## 🎯 核心概念

### API 路由配置

#### 可用路由
```yaml
- 訂單路由: /api/merchants/:mid/orders
- 訂閱路由: /api/merchants/:mid/subscriptions
- 可訂閱實體路由: /api/subscribables_entities
```

#### 添加新路由

1. 在 `/utils/util/var.go` 中定義路由：
```go
const (
    // 現有路由
    SubscriptionPrefix = "/api"
    SubscribablePrefix = "/api/subscribables_entities"
    OrderPrefix = "/api/merchants/:mid/orders"

    // 添加新路由
    NewRoutePrefix = "/api/your/new/route"
)
```

2. 在 `/test/api_test_suite_test.go` 中註冊路由：
```go
func initRoute() {
    // ... 現有代碼 ...
    
    // 註冊新路由
    newService := new.NewService(uow, eventService)
    newHandler := newhandler.NewHandler(newService)
    newHandler.Router(app.Group(apiutil.NewRoutePrefix))
}
```

## 📝 編寫測試

### 基本測試結構

在 `/config/etc/api-test/` 中創建測試工作流程 YAML 文件：

```yaml
api_test:
  name: my_test_workflow
  description: "此測試工作流程的說明"
  scenarios:
    - name: test_scenario
      description: "此特定場景的說明"
      workflow:
        - step: health_check
          request:
            method: GET
            headers: '{"Content-Type": "application/json"}'
            uri: /health_check
          expect_response:
            code: SUCCESS
            status_code: 200
```

### 請求配置

```yaml
request:
  method: POST                    # HTTP 方法 (GET, POST, PUT, PATCH, DELETE)
  uri: /api/endpoint             # API 端點
  headers: '{"key": "value"}'    # 請求頭
  query: '{"param": "value"}'    # 查詢參數
  body: '{"data": "value"}'      # 請求體
  timeout: 30                    # 請求超時時間（秒）（可選）
  retry:                         # 重試配置（可選）
    attempts: 3
    delay: 1
```

### 響應驗證

```yaml
expect_response:
  code: SUCCESS                  # 預期響應代碼
  status_code: 200              # 預期 HTTP 狀態碼
  timeout: 5                    # 響應超時時間（秒）
  equals:                       # 精確值匹配
    - field: data.id
      value: "expected-id"
      message: "ID 應匹配預期值"
  matches:                      # 正則表達式匹配
    - field: data.id
      value: "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
      message: "ID 應匹配 UUID 格式"
  presents:                     # 字段存在性檢查
    - field: data.id
      message: "ID 字段應存在"
  not_presents:                 # 字段不存在性檢查
    - field: data.error
      message: "錯誤字段不應存在"
```

### 變量和動態值

#### 預定義變量
在 `/config/etc/api-test/config/config.yaml` 中定義變量：
```yaml
environment:
  development:
    mid: "6552f99b99821c568c0115cc"
    legacy_id: "SL101PRO6828965740321448017_SKU6828965740942204962"
  production:
    mid: "your-production-mid"
    legacy_id: "your-production-legacy-id"
```

#### 從響應中提取動態變量
```yaml
request:
  vars:
    - name: mid
  from_response:
    - step: create_subscription
      name: subscription_id
      from_field: data.id
      default: "fallback-value"    # 可選默認值
  uri: /api/merchants/{{.mid}}/subscriptions/{{.subscription_id}}
```

### 鉤子系統

#### 全局鉤子
```yaml
api_test:
  name: workflow_with_hooks
  description: "帶全局鉤子的測試工作流程"
  global_hook:
    before:
      workflows:
        - step: setup_environment
          description: "準備測試環境"
    after:
      workflows:
        - step: cleanup_environment
          description: "清理測試環境"
  scenarios:
    - name: test_scenario
      # ... 場景定義
```

#### 場景特定鉤子
```yaml
api_test:
  name: workflow_with_scenario_hooks
  scenarios:
    - name: test_scenario
      hook:
        before:
          workflows:
            - step: scenario_setup
              description: "準備場景特定數據"
        after:
          workflows:
            - step: scenario_cleanup
              description: "清理場景特定數據"
      workflows:
        - step: main_test_step
```

## 💡 最佳實踐

1. **組織結構**
   - 邏輯性地分組相關場景
   - 使用清晰、描述性的場景和步驟名稱
   - 為工作流程和場景添加說明
   - 保持一致的目錄結構

2. **變量管理**
   - 使用環境特定的變量
   - 記錄所有變量依賴關係
   - 在適當的地方提供默認值
   - 使用有意義的變量名稱

3. **測試設計**
   - 保持測試獨立和隔離
   - 使用鉤子進行設置和清理
   - 包含正向和負向測試用例
   - 驗證成功和錯誤場景

4. **驗證策略**
   - 驗證狀態碼和響應體
   - 使用適當的斷言類型
   - 包含有意義的錯誤消息
   - 測試邊界情況和錯誤條件

5. **維護**
   - 定期審查測試用例
   - 移除過時的測試
   - 在 API 變更時更新測試
   - 記錄測試依賴關係

## 🤝 貢獻指南

我們歡迎貢獻！請查看我們的[貢獻指南](CONTRIBUTING.md)了解詳情。

## 📞 支持

- **文檔**：[詳細文檔鏈接]
- **問題反饋**：[GitHub Issues](https://github.com/your-org/github.com/RyanTokManMokMTM/api-testing-go/issues)
- **聯繫方式**：development-team@your-org.com

---

由您的組織 ❤️ 製作 