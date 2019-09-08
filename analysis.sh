
#!/bin/sh

BASEDIR=$(cd $(dirname $0); pwd)

sudo less /var/log/nginx/access.log | kataribe -f ${BASEDIR}/kataribe.toml
sudo /usr/bin/pt-query-digest /var/log/mysql/mysql-slow.log

