package serve

import (
	"fmt"
	"github.com/dexidp/dex/api/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func newDexClient(hostAndPort, clientCert string) (api.DexClient, error) {
	if clientCert != "" {
		creds, err := credentials.NewClientTLSFromFile(clientCert, "")
		if err != nil {
			return nil, fmt.Errorf("load dex cert: %v", err)
		}

		conn, err := grpc.Dial(hostAndPort, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, fmt.Errorf("dial: %v", err)
		}

		return api.NewDexClient(conn), nil
	}

	conn, err := grpc.Dial(hostAndPort, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}

	return api.NewDexClient(conn), nil
}
