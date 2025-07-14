# 勤怠管理システム (Workforce Tracker)

従業員の勤怠管理を行うWebアプリケーションです。出退勤時刻の記録、勤怠データの管理、ユーザー管理機能を提供します。

## 🌟 主な機能

### 認証・認可
- メールアドレス・パスワードによるログイン
- JWT トークンベースの認証
- ロール型アクセス制御（管理者・一般ユーザー）

### 勤怠管理
- ワンクリックでの出勤・退勤打刻
- リアルタイム時刻表示
- 本日の勤怠状況表示
- 勤務時間の自動計算

### ユーザー管理（管理者のみ）
- ユーザーの登録・編集・削除
- ユーザー一覧・検索機能
- ロール管理

## 🛠 技術仕様

### バックエンド
- **言語**: Go 1.24.4
- **フレームワーク**: Echo v4
- **データベース**: PostgreSQL 17.5
- **認証**: JWT (JSON Web Token)
- **パスワードハッシュ化**: bcrypt

### フロントエンド
- **フレームワーク**: Next.js 15.3.3
- **UIライブラリ**: shadcn/ui
- **スタイリング**: Tailwind CSS
- **状態管理**: Zustand
- **フォーム管理**: React Hook Form + Zod

## 🚀 セットアップ・起動方法

### 必要な環境
- Docker & Docker Compose
- Node.js 18+ (開発環境)
- Go 1.24+ (開発環境)

### Docker Composeによる起動

1. リポジトリをクローン
```bash
git clone <repository-url>
cd workforce-tracker
```

2. Docker Composeで全サービスを起動
```bash
docker compose up --build
```

3. アプリケーションにアクセス
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080

### 開発環境での起動

#### データベース（PostgreSQL）の起動
```bash
docker compose up postgres -d
```

#### バックエンドの起動
```bash
cd backend
go mod download
go run ./cmd/server
```

#### フロントエンドの起動
```bash
cd frontend
npm install
npm run dev
```

## 🔐 デフォルトアカウント

初回起動時に管理者アカウントが自動作成されます：

- **Email**: admin@example.com
- **Password**: admin123

⚠️ **重要**: 本番環境では必ずパスワードを変更してください。

## 📱 アプリケーション画面

### ログイン画面
![Login Page](https://github.com/user-attachments/assets/6de6b4dc-2333-46f0-aabc-319fe903ed94)

日本語対応のログインフォームで、デモアカウント情報も表示されます。

### ダッシュボード
![Dashboard](https://github.com/user-attachments/assets/e3de8754-038b-42e3-8ea1-ca0aa4fcc315)

リアルタイム時刻表示と出退勤ボタン、本日の勤怠状況を確認できます。

### 出退勤完了後
![After Clock Out](https://github.com/user-attachments/assets/1c6c380b-7445-448e-8bd1-635881b2e90d)

打刻完了後は成功メッセージが表示され、勤務時間が自動計算されます。

## 🗄 データベース設計

### 主要テーブル

#### users（ユーザー）
- 基本情報（ID、メール、名前、ロール）
- パスワードハッシュ
- 作成・更新・削除日時

#### attendances（勤怠記録）
- ユーザーID、日付
- 出勤・退勤時刻
- 勤務時間（自動計算）

#### modification_requests（修正申請）
- 勤怠データの修正申請
- 承認・却下ワークフロー

#### audit_logs（監査ログ）
- システム操作の記録
- データ変更履歴

## 🔧 API エンドポイント

### 認証API
- `POST /api/auth/login` - ログイン
- `POST /api/auth/logout` - ログアウト
- `GET /api/auth/me` - ユーザー情報取得

### 勤怠API
- `POST /api/attendance/clock-in` - 出勤打刻
- `POST /api/attendance/clock-out` - 退勤打刻
- `GET /api/attendance/today` - 本日の勤怠情報
- `GET /api/attendance/history` - 勤怠履歴

### ユーザー管理API（管理者のみ）
- `GET /api/users` - ユーザー一覧
- `POST /api/users` - ユーザー作成
- `PUT /api/users/:id` - ユーザー更新
- `DELETE /api/users/:id` - ユーザー削除

## 🔒 セキュリティ機能

- JWT による認証
- bcrypt によるパスワードハッシュ化
- ロール型アクセス制御
- CORS 対応
- SQLインジェクション対策
- 入力データ検証

## 📝 開発・運用

### ビルド
```bash
# バックエンド
cd backend && go build -o main ./cmd/server

# フロントエンド
cd frontend && npm run build
```

### テスト
```bash
# バックエンド
cd backend && go test ./...

# フロントエンド
cd frontend && npm run test
```

### 環境変数

#### バックエンド
- `DATABASE_URL`: PostgreSQL接続URL
- `JWT_SECRET`: JWT署名用秘密鍵
- `PORT`: サーバーポート番号

#### フロントエンド
- `NEXT_PUBLIC_API_URL`: バックエンドAPIのURL

## 🤝 今後の機能拡張

- [ ] 勤怠修正申請機能
- [ ] 全ユーザー勤怠管理画面
- [ ] 勤怠履歴閲覧機能
- [ ] CSVエクスポート機能
- [ ] メール通知機能
- [ ] 残業時間管理
- [ ] 有給管理
- [ ] レポート機能

## 📄 ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。