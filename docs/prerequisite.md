## Environment configuration

#### 1. Install protobuf v3.20.1

- Intel chips macbook

  ```shell
  # Download protoc-3.20.1-osx-x86_64.zip
  wget https://github.com/protocolbuffers/protobuf/releases/download/v3.20.1/protoc-3.20.1-osx-x86_64.zip
  
  # Unzip the protoc-3.20.1-osx-x86_64.zip
  unzip protoc-3.20.1-osx-x86_64.zip -d protoc-3.20.1-osx-x86_64
  
  cd protoc-3.20.1-osx-x86_64
  
  # Copy files to system directory
  sudo cp -rf bin/protoc /usr/local/bin
  sudo cp -rf include/google /usr/local/include/
  ```

- M1 chips macbook

  1. Install the Xcode command line tool

     ```shell
     sudo xcode-select --install
     ```

  2. Install port tool (https://www.macports.org), and then execute the following command

     ```shell
     sudo /opt/local/bin/port install autoconf automake libtool
     ```

  3. Build and install Protobuf

     ```shell
     git clone https://github.com/protocolbuffers/protobuf.git
     cd protobuf
     git submodule update --init --recursive
     ./autogen.sh
     ./configure
     make
     make check
     sudo make install
     ```

#### 2. Install golang protobuf

```shell
go install github.com/golang/protobuf/protoc-gen-go@v1.5.2
```

#### 3. Install gRPC

```shell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
```

#### 4. Install gRPC-Gateway

```shell
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.10.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10.0
```

#### 5. Add golang bin path to system PATH

```shell
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zprofile

source ~/.zprofile
```