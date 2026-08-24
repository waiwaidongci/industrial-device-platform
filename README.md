# Industrial Device Platform

纯Go工业物联网设备管理与远程控制服务，提供设备目录、连接心跳、远程命令、固件发布、告警规则和租户领域的可扩展骨架。默认使用内存适配器便于本地启动；PostgreSQL、Redis和消息队列通过接口接入，SQL迁移位于`migrations/`。

## 启动

```sh
go run ./cmd/server
```

服务默认监听`http://localhost:8080`，配置文件为`configs/config.yaml`，环境变量`HTTP_ADDR`和`HEARTBEAT_TIMEOUT`可覆盖配置。

## API示例

```sh
curl localhost:8080/healthz
curl -X POST localhost:8080/v1/devices -H 'Content-Type: application/json' -d '{"tenant_id":"factory-a","name":"PLC-01","kind":"plc","group":"line-1","tags":{"area":"press"}}'
curl 'localhost:8080/v1/devices?tenant_id=factory-a'
curl -X POST localhost:8080/v1/commands -H 'Content-Type: application/json' -H 'Idempotency-Key: demo-1' -d '{"device_id":"dev-id","type":"set_parameter","payload":{"key":"speed","value":10}}'
curl localhost:8080/metrics
```

项目支持请求ID、JSON日志、超时、CORS、统一错误响应、健康/就绪/Prometheus指标和SIGTERM优雅关闭。`cmd/worker`提供独立心跳监控进程。
