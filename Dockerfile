FROM golang:1.19.5-alpine

# アップデートとgitのインストール
RUN apk update && apk add git
RUN apk add bash

# appディレクトリの作成
RUN mkdir /go/src/app

# ワーキングディレクトリの設定
WORKDIR /go/src/app

# Goモジュールの初期化
RUN go mod init example.com/mtg-notification

# godotenvライブラリのインストール
RUN go get github.com/joho/godotenv

# .envファイルのコピー
COPY .env /go/src/app

# ホストのファイルをコンテナの作業ディレクトリに移行
COPY . /go/src/app