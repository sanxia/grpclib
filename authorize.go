package grpclib

import (
	"errors"
	"fmt"
)

import (
	"google.golang.org/grpc/credentials"
)

import (
	"golang.org/x/net/context"
)

/* ================================================================================
 * Oauth Qq
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Rpc 身份认证
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type CustomAuthorize struct {
	Client *RpcClient
}

func (c CustomAuthorize) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("%s %v", "Bearer", c.Client.Token),
	}, nil
}

func (c CustomAuthorize) RequireTransportSecurity() bool {
	if c.Client.credentials != nil {
		return true
	}

	return false
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * SSL证书
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetTLSCredentials(pemPath, keyPath string) (credentials.TransportCredentials, error) {
	if pemPath != "" && keyPath != "" {
		return credentials.NewClientTLSFromFile(pemPath, keyPath)
	}

	return nil, errors.New("credentials error")
}
