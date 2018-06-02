package grpclib

import (
	"fmt"
	"log"
	"reflect"
)

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* ================================================================================
 * Client
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Rpc Client
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type RpcClient struct {
	Host        string
	Port        int
	Token       string
	credentials credentials.TransportCredentials
	conn        *grpc.ClientConn
	clients     map[string]interface{}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化Rpc Client
 * host       : 主机地址
 * port       : 主机端口
 * credentials: SSL证书
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRpcClient(host string, port int, token string, args ...credentials.TransportCredentials) (*RpcClient, error) {
	c := &RpcClient{
		Host:    host,
		Port:    port,
		Token:   token,
		clients: make(map[string]interface{}, 0),
	}

	if len(args) > 0 {
		c.credentials = args[0]
	}

	err := c.connection()

	return c, err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 链接Rpc Server
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RpcClient) connection() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	var opts []grpc.DialOption
	if s.credentials != nil {
		opts = append(opts, grpc.WithTransportCredentials(s.credentials))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	customAuthorize := &CustomAuthorize{}
	customAuthorize.Client = s
	opts = append(opts, grpc.WithPerRPCCredentials(customAuthorize))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	} else {
		s.conn = conn
	}

	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 注册RPC客户端服务
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RpcClient) RegisterClient(key string, value interface{}) {
	if _, isExists := s.clients[key]; !isExists {
		targetValueOf := reflect.ValueOf(value)
		typeOf := targetValueOf.Type()

		log.Printf("%s", typeOf.String())

		if targetValueOf.Kind() == reflect.Func {
			in := []reflect.Value{
				reflect.ValueOf(s.conn),
			}

			out := targetValueOf.Call(in)

			s.clients[key] = out[0].Interface()
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取RPC客户端
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RpcClient) GetClient(key string) interface{} {
	client, isExists := s.clients[key]
	if isExists {
		return client
	}
	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 关闭Rpc客户端连接
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RpcClient) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}
