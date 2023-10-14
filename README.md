# CharFluent

## 実装メモ

### Userの作成

データベースの準備

```bash
$ psql postgres
> CREATE DATABASE chatfluent;
$ psql -f data/setup.sql -d chatfluent
```

データベースとの接続  

DbConnection.Exex, DbConnection.Prepareの使い分けがわからない。  

Userの作成

uuidの生成
`github.com/nu7hatch/gouuid`を参考に実装する。  
uuidはセッションの管理に使用する。  





