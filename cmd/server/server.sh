
export ADDRESS=localhost:44751
export REPORT_INTERVAL=10s
export POLL_INTERVAL=2s
export RESTORE=false
export SHUTDOWN_TIMEOUT=5s
export STORE_INTERVAL=10s
export STORE_FILE=.\devops-metrics-db.json
export RESTORE=true
export KEY=/tmp/zK3deX0
export DATABASE_DSN=postgres://postgres:qweasd@localhost:5432/yandex?sslmode=disable

go run ./cmd/server/main.go 