
#!/bin/sh

BASEDIR=$(cd $(dirname $0); pwd)

sudo less /var/log/nginx/access.log | kataribe -f ${BASEDIR}/kataribe.toml
sudo /bin/pt-query-digest /var/log/mariadb/mariadb-slow.log
