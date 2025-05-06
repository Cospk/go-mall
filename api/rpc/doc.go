// Package rpc
/*
  rpc 包用于定义proto文件生成的rpc服务接口
  该包下的所有proto文件通过“protoc --go_out=./api/rpc --go-grpc_out=./api/rpc --proto_path=./api/rpc  文件名”
  命令生成对应的grpc服务接口到gen目录下，同时会生成对应的pb.go文件
  单独定义一个gen目录是为了避免生成的pb.go文件被覆盖，同时也方便管理
  每次生成两个文件，一个是pb.go文件，一个是pb.gw.go文件，其中pb.go文件是用于定义proto文件生成的rpc服务接口，pb.gw.go文件是用于定义proto文件生成的http服务接口
  开发客户端时，只需要引入pb.go文件，不需要引入pb.gw.go文件，因为pb.gw.go文件是用于定义http服务接口的，客户端不需要使用http服务接口，只需要使用rpc服务接口即可
  开发服务端时，需要引入pb.go文件和pb.gw.go文件，因为pb.gw.go文件是用于定义http服务接口的，服务端需要使用http服务接口，同时也需要使用rpc服务接口
*/
package rpc
