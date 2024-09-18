LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	 make generate-auth-api
	 make generate-user-api

generate-auth-api:
	mkdir -p backend/auth-service/api/auth
	protoc --proto_path shared/proto/auth \
            --go_out=backend/auth-service/api/auth --go_opt=paths=source_relative \
            --plugin=protoc-gen-go=bin/protoc-gen-go \
            --go-grpc_out=backend/auth-service/api/auth --go-grpc_opt=paths=source_relative \
            --plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
            shared/proto/auth/auth.proto``

generate-user-api:
	# mkdir -p backend/user-service/api/user
	protoc --proto_path shared/proto/user \
            --go_out=backend/user-service/api/user --go_opt=paths=source_relative \
            --plugin=protoc-gen-go=bin/protoc-gen-go \
            --go-grpc_out=backend/user-service/api/user --go-grpc_opt=paths=source_relative \
            --plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
            shared/proto/user/user.proto``

run:
	make run-auth-service

run-auth-service:
	cd backend/auth-service/ && go run cmd/main.go