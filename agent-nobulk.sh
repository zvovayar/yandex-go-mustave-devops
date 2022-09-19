#export ADDRESS=192.168.0.23:44751
#export ADDRESS=192.168.0.23:44751
#export ADDRESS=localhost:44751
export REPORT_INTERVAL=1s
export POLL_INTERVAL=1s
export RESTORE=false
export SHUTDOWN_TIMEOUT=5s
export KEY=/tmp/zK3deX0
export CRYPTO_KEY=./publickey.pem
export GRPC_SRV_ADDR=192.168.0.23:3222


go run ./cmd/agent/main.go -B=false

