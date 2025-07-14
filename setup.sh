#!/bin/bash

echo "🏢 勤怠管理システム セットアップスクリプト"
echo "================================================="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker が見つかりません。Dockerをインストールしてください。"
    exit 1
fi

# Check if Docker Compose is available
if ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose が見つかりません。Docker Composeをインストールしてください。"
    exit 1
fi

echo "✅ Docker環境を確認しました"

# Pull and start PostgreSQL
echo "📦 PostgreSQLコンテナを起動中..."
docker compose up postgres -d

# Wait for PostgreSQL to be ready
echo "⏳ PostgreSQLの起動を待機中..."
sleep 10

# Build and start backend
echo "🚀 バックエンドサーバーを起動中..."
docker compose up backend -d

# Wait for backend to be ready
echo "⏳ バックエンドサーバーの起動を待機中..."
sleep 15

# Build and start frontend
echo "🌐 フロントエンドサーバーを起動中..."
docker compose up frontend -d

echo ""
echo "🎉 セットアップ完了！"
echo ""
echo "📱 アプリケーションにアクセス:"
echo "   フロントエンド: http://localhost:3000"
echo "   バックエンドAPI: http://localhost:8080"
echo ""
echo "🔐 デフォルトアカウント:"
echo "   Email: admin@example.com"
echo "   Password: admin123"
echo ""
echo "📋 コマンド:"
echo "   ログ確認: docker compose logs -f"
echo "   停止: docker compose down"
echo "   再起動: docker compose restart"
echo ""