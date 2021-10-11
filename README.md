gRPCで動画のアップロードを試してみる
root直下にsample.mp4を作成して実行
### 生成コマンド
```bash
protoc --go_out=./proto --proto_path=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=require_unimplemented_servers=false --go-grpc_opt=paths=source_relative ./proto/upload.proto
```
protoc-gen-go-grpcではrequire_unimplemented_servers=falseに設定する([参考](https://github.com/grpc/grpc-go/tree/master/cmd/protoc-gen-go-grpc))


### memo
go_package="{import_path};{package_name}"でgoのpackage名を指定する


## TODO
- multpartで受け取ったファイルをアップロードして保存
- cloudrunで受け取りcloud storageに保存