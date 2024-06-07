
## setup
手っ取り早く確認する用
1. GoogleAIのAPIキーを取得
2. docker-composeのGOOGLE_AI_API_KEYにセット
3. docker-composeを起動\
    ```bash
    docker compose -f "dockercompose.yaml" up -d --build 
    ```

## dev

### 環境
- go 1.22.4
- Bun 1.1.9 or Node.js 21.6.2


### ElasticSearch起動
```bash
docker compose  -f "dockercompose.yaml" up -d --build elasticsearch 
```

### サーバー起動
```bash
go run backend/cmd/server/main.go
```
### フロントエンド起動
```bash
cd web

npm install
npm run dev
# or
bun install
bun dev
```