go mod init payment
go mod tidy

protoc --go_out=. --go-grpc_out=. --go_opt=Mproto/payment.proto=payment-worker/proto --go-grpc_opt=Mproto/payment.proto=payment-worker/proto proto/payment.proto

protoc --go_out=. --go-grpc_out=. --go_opt=Mproto/payment.proto=integration/proto --go-grpc_opt=Mproto/payment.proto=integration/proto proto/payment.proto

protoc --go_out=. --go-grpc_out=. --go_opt=Mproto/account.proto=account/proto --go-grpc_opt=Mproto/account.proto=account/proto proto/account.proto

protoc --go_out=. --go-grpc_out=. --go_opt=Mproto/account.proto=payment/proto --go-grpc_opt=Mproto/account.proto=payment/proto proto/account.proto
