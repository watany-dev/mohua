# SageMaker Monitor

## 概要

SageMaker Monitorは、AWS SageMakerのコンピュートリソースを監視し、リアルタイムでコスト分析を行うCLIツールです。このツールは、以下のSageMakerリソースの状態と関連コストを追跡します：

- エンドポイント (Endpoints)
- ノートブックインスタンス (Notebook Instances)
- Studioアプリケーション (Studio Applications)

## 主な機能

- 🔍 リアルタイムのリソース状態監視
- 💰 詳細なコスト分析
  - 現在までの累積コスト
  - 時間あたりのコスト
  - 月間予測コスト
- ⚠️ コスト超過警告
- 📊 柔軟な出力フォーマット（テーブル/JSON）

## 前提条件

- Go 1.16以上
- AWS CLI設定済み
- AWS IAMアクセス権限

## インストール

### 方法1: ソースからビルド

```bash
# リポジトリをクローン
git clone https://github.com/yourusername/sagemaker-monitor.git
cd sagemaker-monitor

# 依存関係のインストール
go mod tidy

# ビルド
go build -o sagemaker-monitor

# インストール（オプション）
go install
```

### 方法2: バイナリダウンロード

[Releases](https://github.com/yourusername/sagemaker-monitor/releases)ページから最新のバイナリをダウンロードしてください。

## 使用方法

### 基本的な使用方法

```bash
# テーブル形式で表示
./sagemaker-monitor --region us-east-1

# JSON形式で出力
./sagemaker-monitor --region us-east-1 --json
```

### コマンドラインオプション

- `--region, -r`: AWS リージョンを指定（必須）
- `--json, -j`: JSON形式で出力

## 環境設定

AWS認証情報は以下のいずれかの方法で設定できます：

1. AWS CLI設定
```bash
aws configure
```

2. 環境変数
```bash
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_DEFAULT_REGION=us-east-1
```

3. IAMロール（EC2またはECS）

## 出力例

### テーブル形式
```
Type            Name                          Status     Instance       Running Time  Hourly($)  Current($)  Projected($)
Endpoint        my-ml-endpoint                InService  ml.m5.xlarge   72h 15m       $1.24      $89.54      $912.80
Notebook        dev-notebook                  Running    ml.t3.medium   168h 30m      $0.11      $18.54      $81.40

Total Current Cost: $108.08    Projected Monthly Cost: $994.20

WARNING: Endpoint 'my-ml-endpoint' has a high projected monthly cost: $912.80
```

### JSON形式
```json
[
  {
    "resourceType": "Endpoint",
    "name": "my-ml-endpoint",
    "status": "InService",
    "instanceType": "ml.m5.xlarge",
    "runningTime": "72h 15m",
    "hourlyCost": 1.24,
    "currentCost": 89.54,
    "projectedMonthlyCost": 912.80
  },
  ...
]
```

## 注意事項

- コスト計算は概算であり、実際の請求額とは異なる場合があります
- 最新の料金情報に基づいて`configs/pricing.yaml`を定期的に更新してください

## トラブルシューティング

- AWS認証エラー: IAMポリシーと権限を確認
- リージョン指定エラー: 正確なリージョン名を使用
- 予期せぬ結果: AWS SDKのバージョンを確認

## 貢献

プルリクエストや機能提案を歓迎します。詳細は`CONTRIBUTING.md`を参照してください。

## ライセンス

このプロジェクトは[MITライセンス](LICENSE)の下で公開されています。

## 免責事項

このツールは情報提供のみを目的としており、正確な請求情報については常にAWSコンソールを確認してください。
