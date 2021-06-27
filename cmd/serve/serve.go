package serve

import (
	"context"
	"errors"
	"fmt"
	"github.com/dexidp/dex/api/v2"
	"github.com/norganna/dex-account-api/storage"
	"log"
	"net"
	"net/http"

	"github.com/spf13/viper"
)

type serveOptions struct {
	WebHTTPAddr string `mapstructure:"web-http-addr"`

	WebHTTPSAddr string `mapstructure:"web-https-addr"`
	WebTLSCert   string `mapstructure:"web-tls-cert"`
	WebTLSKey    string `mapstructure:"web-tls-key"`

	GrpcAddr          string `mapstructure:"grpc-addr"`
	GrpcTLSClientCert string `json:"grpc-tls-client-cert"`

	Store struct {
		Class string `mapstructure:"class"`
	} `mapstructure:"store"`

	configFile string
	store      storage.Storage
	dex        api.DexClient
	ctx        context.Context
	done       <-chan bool
}

func runServe(o *serveOptions) error {
	o.ctx = context.Background()

	if configFile := o.configFile; configFile != "" {
		viper.SetConfigFile(configFile)
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file: %s: %v", viper.ConfigFileUsed(), err)
		}
	}

	err := viper.Unmarshal(o)
	if err != nil {
		return fmt.Errorf("error parsing config: %v", err)
	}

	if o.Store.Class == "" {
		o.Store.Class = "memstore"
	}

	o.store = storage.NewStore(o.Store.Class)
	if o.store == nil {
		return fmt.Errorf("error creating %s store: not found", o.Store.Class)
	}

	_ = viper.Unmarshal(o.store.Config())

	if o.WebHTTPSAddr == "" && o.WebHTTPAddr == "" {
		return errors.New("neither http nor https address specified")
	}

	mux := http.NewServeMux()

	o.dex, err = newDexClient(o.GrpcAddr, o.GrpcTLSClientCert)
	if err != nil {
		return fmt.Errorf("failed creating dex client: %v ", err)
	}

	mux.HandleFunc("/challenge/", challengeHandler(o))
	mux.HandleFunc("/create/", createHandler(o))
	mux.HandleFunc("/update/", updateHandler(o))

	var tls bool
	proto := "http"
	addr := o.WebHTTPAddr
	if o.WebHTTPSAddr != "" {
		proto = "https"
		addr = o.WebHTTPSAddr
		tls = true
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := &http.Server{
		Handler: mux,
	}

	done := make(chan bool)
	o.done = done

	go func() {
		defer close(done)
		if tls {
			_ = server.ServeTLS(listener, o.WebTLSCert, o.WebTLSKey)
		} else {
			_ = server.Serve(listener)
		}
	}()

	log.Printf("Web listening on %s://%s", proto, listener.Addr().String())
	return nil
}

func (o *serveOptions) Wait() {
	select {
	case <-o.ctx.Done():
	case <-o.done:
	}
}

// Wait waits for the running server (if any) to finish.
func Wait() {
	if runningServer != nil {
		runningServer.Wait()
	}
}
