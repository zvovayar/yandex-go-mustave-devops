
#export ADDRESS=localhost:44751
#export ADDRESS=192.168.0.23:44751
export ADDRESS=0.0.0.0:44751
export REPORT_INTERVAL=10s
export POLL_INTERVAL=2s
export RESTORE=false
export SHUTDOWN_TIMEOUT=5s
export STORE_INTERVAL=10s
export STORE_FILE=./devops-metrics-db.json
export RESTORE=true
#export KEY=q
export KEY=/tmp/zK3deX0
export DATABASE_DSN=postgres://postgres:qweasd@localhost:5432/yandex?sslmode=disable
#export DATABASE_DSN="'***postgres:9999/praktikum?easter_egg_msg=you_must_prefer_this_incorrect_settings_to_those_obtained_through_arguments'"
export CRYPTO_KEY=./privatekey.pem

go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildDate=23.08.2022 -X main.buildCommit=CWCWCWC" ./cmd/server/main.go 
