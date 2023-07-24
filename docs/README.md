# MASMAINT-CG 
### マスタメンテナンス自動生成
https://masmaint-cg.carkodr.com/


## CSV設定
```
T, テーブル名, テーブル和名  
C, カラム名, 型分類, PKフラグ, NOTNILLフラグ, 登録可フラグ, 更新可フラグ  
```

|型分類||
|:---:|----|
|S|文字列型|
|I|整数型|
|F|少数型|
|T|日付型|

|フラグ||
|:---:|----|
|1|YES|
|0|NO| 

|CSVサンプル||
|----|----|
|T,employee,従業員マスタ|
|C,id,I,1,0,0,0|PK 登録・更新不可(AUTOINCREMENT)|
|C,name,S,0,1,1,1|NOTNULL 登録・更新可能|
|C,hire_date,T,0,1,1,0|NOTNULL 登録可能 更新不可|
|C,saraly,F,0,0,1,1|登録・更新可能|
|C,created_at,T,0,1,0,0|登録・更新不可(トリガー更新)|
|C,updated_at,T,0,1,0,0|登録・更新不可(トリガー更新)|

`*続けて複数テーブル指定可能*`



## プログラム起動設定
### Go 1.18~ (Gin)
#### 依存モジュールインストール
```bash
cd masmaint
go mod init masmaint
go mod tidy
```
#### 環境変数ファイル設定
```css
masmaint
└── config
     └── env
         └── local.env
```
```
LOG_LEVEL=DEBUG
APP_HOST=localhost
APP_PORT=3000        //環境に合わせて
DB_NAME=masmaint     //環境に合わせて （SQLite3の場合は masmaintフォルダからの相対パス or 絶対パス）
DB_HOST=localhost    //環境に合わせて
DB_PORT=5432         //環境に合わせて
DB_USER=postgres     //環境に合わせて
DB_PASSWORD=postgres //環境に合わせて
```

#### 起動
```bash
ENV=local go run cmd/masmaint/main.go
```
ブラウザでアクセス  
http://localhost:3000/mastertables

### PHP 7 ~ 8
バージョンは細かく検証したわけではないため、動かない可能性もある。
#### 依存モジュールインストール
```bash
composer install
```
#### 環境変数ファイル設定
```css
masmaint
└── env
     └── .env
```
```
DB_DRIVER=mysql      //環境に合わせて
DB_HOST=localhost    //環境に合わせて
DB_PORT=3307         //環境に合わせて
DB_NAME=masmaint     //環境に合わせて （SQLite3の場合は masmaintフォルダからの相対パス or 絶対パス）
DB_USER=root         //環境に合わせて
DB_PASS=root         //環境に合わせて
```

#### 起動
```bash
composer start
```
ブラウザでアクセス  
http://localhost:8080/mastertables

## License
Copyright © 2023 Murakami Koudai
