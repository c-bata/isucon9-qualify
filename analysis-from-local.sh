#!/bin/sh

USER=${ISUCON_USER:-isucon}
PEM=${PEM:-~/.ssh/isucon9.pem}
IPADDR=${IPADDR:-47.91.19.40}

ssh -i $PEM ${USER}@${IPADDR} <<EOF
sudo less /var/log/nginx/access.log | kataribe -f ~/isucari/kataribe.toml
sudo /usr/bin/pt-query-digest /var/log/mysql/mysql-slow.log
EOF

