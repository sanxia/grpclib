# grpclib
gRPC and Portobuf for Client and Server Lib
--------------------------------------------

### Server
import (

    "github.com/sanxia/grpclib"

)

const (

    rpcHost = "127.0.0.1"

    rpcPort = 5111

)

func main() {

    RpcServer(rpcHost, rpcPort)

}

func RpcServer(host string, port int) {

    //ssl证书
    creds, _ := grpclib.GetTLSCredentials("", "")

    //rpc server new
    rpcServer := grpclib.NewRpcServer(host, port, creds)

    //register rpc service
    rpcServer.RegisterService(func(rpcSrv *grpc.Server) {

        protobuf.RegisterActionServer(rpcSrv, &serverImpl{})

    })

    //authorize enabled
    rpcServer.Authorize(authorize)

    //logger enabled
    rpcServer.Logger()

    rpcServer.Serve()
}

### Client
import (

    "github.com/sanxia/grpclib"

)

const (

    rpcHost = "127.0.0.1"

    rpcPort = 5111

    token   = "test-token-2018"

)

func main() {

    //ssl证书
    tlsCreds, _ := grpclib.GetTLSCredentials("", "")

    //实例化RpcClientImpl
    rpcClientImpl, err := NewRpcClientImpl(rpcHost, rpcPort, token, tlsCreds)

    if err != nil {

        log.Fatalf("Failed connection rpc server: %v", err)

        return

    }

    //注册Rpc客户端服务

    rpcClientImpl.RegisterClient("action", protobuf.NewActionClient)

    defer rpcClientImpl.Close()

    //调用Rpc方法
    rpcClientImpl.Song("圣诞歌")

    rpcClientImpl.MakeAWish("老刘", "亲们", "元节快乐!")

    rpcClientImpl.MakeAWish("老刘", "才子佳人", "圣诞节快乐!")

}

### Protobuf

syntax = "proto3";

package protobuf;

service Action {

  rpc Song (SongRequest) returns (SongReply) {}

  rpc MakeAWish (MakeAWishRequest) returns (stream MakeAWishReply){}

}

message SongRequest {

  string title = 1;

}

message SongReply {

  string lyric = 1;

  string singer = 2;

  uint32 year = 3;

}

message MakeAWishRequest {

    string from = 1;

    string to = 2;

    string content = 3;

}

message MakeAWishReply {

    string content = 1;

}
