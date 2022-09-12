export ADDRESS=192.168.0.23:44751
#export ADDRESS=localhost:44751
export REPORT_INTERVAL=1s
export POLL_INTERVAL=1s
export RESTORE=false
export SHUTDOWN_TIMEOUT=5s
export KEY=/tmp/zK3deX0
export CRYPTO_KEY=./publickey.pem

go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildDate=23.08.2022 -X main.buildCommit=CWCWCWC" ./cmd/agent/main.go -B=true

