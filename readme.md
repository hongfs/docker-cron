Docker-Cron 用于在 Docker 容器中运行 Cron 任务，执行定期清除应用日志、定期下载最新配置文件等任务，需配合 Docker-Compose 使用。

```yaml
version: "3"
services:
  cron:
    image: ghcr.io/hongfs/docker-cron:main
    environment:
      CONFIG: |-
        list:
        - name: 清除 ThinkPHP 过期日志
          cron: "0 0 * * *"
          command: find /var/www/html/runtime/log -mtime +7 -name "*.log" -exec rm -rf {} \;
        - name: 清除 ThinkPHP 过期缓存
          cron: "0 0 * * *"
          command: find /var/www/html/runtime/cache -mtime +7 -name "*.php" -exec rm -rf {} \;
```
