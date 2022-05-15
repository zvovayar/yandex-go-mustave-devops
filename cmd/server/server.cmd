rem set ADDRESS=localhost:8081
rem set POLL_INTERVAL=3
rem set REPORT_INTERVAL=11
set ADDRESS=localhost:44751
set REPORT_INTERVAL=10s
set POLL_INTERVAL=2s
set RESTORE=false
set SHUTDOWN_TIMEOUT=5s
set STORE_INTERVAL=10s
set STORE_FILE=.\devops-metrics-db.json
set RESTORE=true

go run .\cmd\server\main.go