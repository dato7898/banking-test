go mod init payment
go mod tidy

protoc --go_out=. --go-grpc_out=. --go_opt=Mproto/payment.proto=payment-worker/proto --go-grpc_opt=Mproto/payment.proto=payment-worker/proto proto/payment.proto

protoc --go_out=. --go-grpc_out=. --go_opt=Mproto/payment.proto=integration/proto --go-grpc_opt=Mproto/payment.proto=integration/proto proto/payment.proto
