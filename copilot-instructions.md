# Copilot Instructions

## Lint 檢查規範

- 本專案請務必在每次新增或修改程式碼後，執行 lint 檢查，確保程式品質。
- 請使用 [golangci-lint](https://github.com/golangci/golangci-lint) 工具。
- 執行指令如下：

```bash
golangci-lint run ./...
```

- 若有 lint 問題，請先修正後再進行 commit。

---