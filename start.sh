sudo docker-compose up -d
go build -o out/consumer cmd/kafkaconsumer/consumer.go
./out/consumer