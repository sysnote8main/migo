# エンドポイントラッパー再生成手順

Misskey の OpenAPI spec (`api.json`) を元に、全エンドポイントの Go 型定義・Service メソッドを自動生成します。

## 前提

- プロジェクトルートに `api.json` が配置されていること
- `tools/gen/main.go` が存在すること（コード生成ツール）

## 再生成手順

### 1. OpenAPI spec を取得

```bash
curl -s 'https://crafters.aosankaku.net/api.json' -o api.json
```

### 2. 既存の生成ファイルを削除

```bash
# 全 gen_all.go を削除
rm -f */gen_all.go
```

### 3. コード生成を実行

```bash
go run tools/gen/main.go
```

### 4. 生成物を確認

```bash
# ビルドが通ることを確認
go build ./...

# テストが通ることを確認
go test ./...
```

### 5. 生成ファイルの diff を確認してコミット

```bash
git diff
git add -A
git commit -m "feat: regenerate endpoint wrappers from OpenAPI spec"
```

### 6. api.json はコミットしない

`api.json` は `.gitignore` に含まれているため、コミットされません。  
必要に応じて再取得してください。

## カスタマイズ

### covered の更新

`tools/gen/main.go` 内の `covered` マップに、手書き実装済みで生成から除外したいエンドポイントを追加します。

例：

```go
"/admin/emoji/add": true,  // 手書き実装済みのため生成対象外
```

### knownTypes の更新

`tools/gen/main.go` 内の `knownTypes` マップに、Go の型と OpenAPI schema 名の対応を追加します。

例：

```go
"NewType": "types.NewType",
```

### 新しいパッケージを services.go に追加

生成ツールは既存パッケージに対しては Service 型のボイラープレートを出力しませんが、
新規パッケージに対しては `Service` struct + `NewService` を生成します。

新規パッケージを `Services` バンドルに追加するには `services.go` を手動で編集してください。

## アーキテクチャ

```
api.json (OpenAPI spec)
    │
    ▼
tools/gen/main.go (Go コード生成ツール)
    │
    ▼
各パッケージ/gen_all.go (生成されたラッパー)
```

- 各パッケージの `gen_all.go` は **"DO NOT EDIT"** です。再生成すると上書きされます。
- 手書きの実装が必要なエンドポイントは各パッケージの `.go` ファイル（`admin.go`, `notes.go` など）に記述し、`covered` マップに追加することで生成対象から除外します。
- 型解決できないフィールドは `any` になります。必要に応じて手書きの型定義と差し替えてください。
