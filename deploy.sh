#!/bin/sh

set -ex

USER=${ISUCON_USER:-isucon}
PEM=${PEM:-~/.ssh/isucon9.pem}
IPADDR=${IPADDR:-47.91.19.40}

BRANCH=master
if [ $# -eq 1 ]; then
  BRANCH=$1
fi

ssh -i $PEM ${USER}@${IPADDR} <<EOF

sudo su - isucon
cd ~/isucari  # HOME環境変数はlocalが使われるので注意

# git branch 切り替え
git fetch origin -p
git checkout origin/$BRANCH
git branch
git log --graph --all --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit --date=relative | head -n 10

# go build
cd ./webapp/go/
make all
cd ~/isucari

# log rotate
if [ -f /var/log/nginx/access.log ]; then
     sudo mv /var/log/nginx/access.log /var/log/nginx/access.log.$(date +"%Y%m%d_%H%M%S")
fi
if [ -f /var/log/mysql/mysql-slow.log ]; then
    sudo mv /var/log/mysql/mysql-slow.log /var/log/mysql/mysql-slow.log.$(date +"%Y%m%d_%H%M%S")
fi

# service restart
sudo systemctl restart isucari.golang.service
sudo systemctl restart nginx.service
sudo systemctl restart mysql.service
EOF
