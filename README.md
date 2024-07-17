# gin-naiveui

## [api 接口](./api.md)

## 本地运行

#### 运行

docker 方式

```shell
docker compose -f docker-compose-env.yaml
```

或者自行修改 .env 文件 Mysql 连接参数，并导入 init.sql 数据库及表结构

##### 运行前端

```shell
cd vue-naive-front && npm install && npm run dev
```

##### 运行后端

```shell
go run main.go
```
