FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  echo http://dl-cdn.alpinelinux.org/alpine/edge/testing >> /etc/apk/repositories && \
  apk --no-cache add ca-certificates && \
  apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
  echo "Asia/Shanghai" > /etc/timezone
  

# 设置时区为上海

WORKDIR /app

COPY ./bin/app /app/app
# RUN ls /app && cat /app/config/config.toml
# Command to run the executable
ENTRYPOINT ["/app/app"]
# CMD []
