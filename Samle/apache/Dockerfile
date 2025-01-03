# ベースイメージ
FROM php:8.2-apache

# 必要なモジュールのインストール
RUN docker-php-ext-install mysqli

# 作業ディレクトリの設定
WORKDIR /var/www/html

# ファイルをコンテナ内にコピー
COPY html/ /var/www/html/

# Apacheの設定をリロード可能にする
RUN a2enmod rewrite
