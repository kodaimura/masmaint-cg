# masmaint
## 準備 (ローカルで起動する場合)
* docker 起動

#### sqlite3の場合
config/env/local.env
```
DB_NAME="sqlite3のファイルのパス"
```
## 起動
```
cd path/to/masmaint
make dev
```
http://localhost:3000  

ユーザ名　：user  
パスワード：pass  
（簡易ログインのためカスタム推奨 internal/server/router.go 参照）

## その他
* Makefile 参照
