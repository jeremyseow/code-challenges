
create_discriptor:
	protoc --proto_path=. --include_imports --descriptor_set_out=./descriptor pb/event.proto

run_app:
	go run main.go

create_gogo_code:
	protoc --proto_path=. --proto_path="C:/Users/Jeremy Seow/go/pkg/mod/github.com/gogo/protobuf@v1.3.2" --gogo_out=. ./pb/event.proto

create_code:
	protoc --proto_path=. --go_out=. ./pb/event.proto