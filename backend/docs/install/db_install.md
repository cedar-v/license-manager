# MySQL安装

## Docker安装

```bash

docker run --name lm-mysql -d \
  --restart always \
  -p 13306:3306 \
  -e MYSQL_ROOT_PASSWORD=root@123 \
  -e MYSQL_DATABASE=license_manager \
  -v mysql-data:/var/lib/mysql \
  -v mysql-conf:/etc/mysql/conf.d \
  -v mysql-logs:/var/log/mysql \
  mysql:8.0 \
  --default-authentication-plugin=mysql_native_password
```

## 卸载清理

```bash
docker rm -f lm-mysql
docker volume rm mysql-data mysql-conf mysql-logs
```