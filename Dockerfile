# ---------- build stage ----------
FROM golang:1.23 AS builder

WORKDIR /app

# 先拷贝依赖文件，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 再拷贝全部源码
COPY . .

# 编译入口：你的入口是 cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go


# ---------- run stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/server /app/server

# Railway 会注入 PORT
ENV PORT=8080
EXPOSE 8080

CMD ["/app/server"]
