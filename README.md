Go微服务开发框架使用指南

### 一、基础环境配置

##### 1、安装protobuf v3.17.3

地址：https://github.com/protocolbuffers/protobuf/releases 

下载系统对应版本 [protoc-3.17.3-osx-x86_64.zip]，并解压

拷贝bin目录下protoc到/usr/local/bin

拷贝include目录下的google到/usr/local/include

##### 2、安装golang protobuf

```shell
go install github.com/golang/protobuf/protoc-gen-go@v1.5.2
```

##### 3、安装grpc

```shell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
```

##### 4、安装grpc-gateway

```shell
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.4.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.4.0
```

##### 5、安装protoc-go-inject-tag

```shell
go install github.com/favadi/protoc-go-inject-tag@latest
```

### 二、microctl工具安装与使用

##### 1、安装microctl工具

```shell
go get -u https://github.com/imind-lab/micro/microctl@latest
# 把microctl加入系统PATH
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zprofile
source ~/.zprofile

# 查看microctl版本
microctl version
```

##### 2、运行工具，生成实例代码

```shell
cd $GOPATH/src

# init 子命令 初始化微服务
# -p 项目名 默认值imind-lab
# -s 微服务名 默认值greet
# -a 是否生成api-gateway 默认值是true
microctl init -p imind-lab -s greet -a=true

cd github.com/imind-lab/greet/build

# 部署greet服务到kubernetes
make deploy
```

### 
