package grpclib

import (
	"fmt"
	"log"
	"net"
	"time"
)

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/* ================================================================================
 * Oauth Qq
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * RPC Server
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type RpcServer struct {
	Host               string
	Port               int
	rpc                *grpc.Server
	registerService    func(*grpc.Server)
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	authorizeHandler   func(context.Context) (context.Context, error)
	credentials        credentials.TransportCredentials
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化Rpc Server
 * host: 主机地址
 * port: 主机端口
 * auth: 授权
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRpcServer(host string, port int, args ...credentials.TransportCredentials) *RpcServer {
	s := &RpcServer{
		Host:               host,
		Port:               port,
		unaryInterceptors:  make([]grpc.UnaryServerInterceptor, 0),
		streamInterceptors: make([]grpc.StreamServerInterceptor, 0),
	}

	if len(args) > 0 {
		s.credentials = args[0]
	}

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 启动服务器
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RpcServer) Serve() {
	serverOpts := make([]grpc.ServerOption, 0)

	//SSL
	if s.credentials != nil {
		serverOpts = append(serverOpts, grpc.Creds(s.credentials))
	}

	//拦截器链
	if s.authorizeHandler != nil {
		s.unaryInterceptors = append(s.unaryInterceptors, grpc_auth.UnaryServerInterceptor(s.authorizeHandler))
		s.streamInterceptors = append(s.streamInterceptors, grpc_auth.StreamServerInterceptor(s.authorizeHandler))
	}

	s.unaryInterceptors = append(s.unaryInterceptors, grpc_recovery.UnaryServerInterceptor())
	s.streamInterceptors = append(s.streamInterceptors, grpc_recovery.StreamServerInterceptor())

	//配置
	serverOpts = append(serverOpts, grpc_middleware.WithUnaryServerChain(
		s.unaryInterceptors...,
	))
	serverOpts = append(serverOpts, grpc_middleware.WithStreamServerChain(
		s.streamInterceptors...,
	))

	//实例化RPC Server
	s.rpc = grpc.NewServer(serverOpts...)

	//注册Rpc服务
	s.registerService(s.rpc)

	//监听
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.rpc.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 注册服务
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RpcServer) RegisterService(registerService func(*grpc.Server)) {
	s.registerService = func(rpc *grpc.Server) {
		registerService(rpc)

		reflection.Register(rpc)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 是否启用日志
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RpcServer) Authorize(auth func(context.Context) (context.Context, error)) {
	s.authorizeHandler = auth
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 是否启用日志
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RpcServer) Logger() {
	zapLogger := zap.NewExample()
	zapLoggerOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}
	grpc_zap.ReplaceGrpcLogger(zapLogger)

	s.unaryInterceptors = append(s.unaryInterceptors, grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)))
	s.unaryInterceptors = append(s.unaryInterceptors, grpc_zap.UnaryServerInterceptor(zapLogger, zapLoggerOpts...))

	s.streamInterceptors = append(s.streamInterceptors, grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)))
	s.streamInterceptors = append(s.streamInterceptors, grpc_zap.StreamServerInterceptor(zapLogger, zapLoggerOpts...))
}
