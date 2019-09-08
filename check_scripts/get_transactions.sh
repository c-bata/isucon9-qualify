#!/bin/bash

HOST=${HOST:-https://isucon9.catatsuy.org}

curl -XGET ${HOST}/initialize # 初期化
id=$(curl -XPOST ${HOST}/register --data '{"account_name":"foobar","address":"my address","password":"foobar"}' -H "Content-Type: application/json" | jq '.id')
curl -XPOST -c cookie.txt ${HOST}/login --data '{"account_name": "foobar", "password": "foobar"}' -H "Content-Type: application/json" # ログイン
cat cookie.txt

curl -XGET -b cookie.txt ${HOST}/users/transactions.json | jq .

