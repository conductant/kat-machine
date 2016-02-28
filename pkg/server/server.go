package server

import (
	"github.com/conductant/gohm/pkg/resource"
	"github.com/conductant/gohm/pkg/server"
	"github.com/conductant/gohm/pkg/version"
	"github.com/conductant/kat-machine/pkg/machine"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"net/http"
)

type ServerOptions struct {
	Port         int    `json:"port" yaml:"port" flag:"port, The server listening port"`
	PublicKeyUrl string `json:"public_key_url,omitempty" yaml:"public_key_url" flag:"public_key_url,Url for fetching the public key for auth token"`
}

type Server struct {
	ServerOptions

	publicKey []byte
}

func (this *Server) Init() error {
	// TODO - this is just for dev
	if this.PublicKeyUrl == "" {
		return nil
	}

	buff, err := resource.Fetch(context.Background(), this.PublicKeyUrl)
	if err != nil {
		return err
	}
	this.publicKey = buff
	return nil
}

func (this *Server) Start() <-chan error {
	shutdown := make(chan struct{})
	stop, stopped := server.NewService().
		WithAuth(
			server.Auth{
				VerifyKeyFunc: func() []byte {
					return this.publicKey
				},
			}.Init()).
		ListenPort(this.Port).
		Route(
			server.Endpoint{
				UrlRoute:   "/info",
				HttpMethod: server.GET,
				AuthScope:  server.AuthScopeNone,
			}).
		To(
			func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
				glog.Infoln("Showing version info.")
				server.Marshal(resp, req, version.BuildInfo())
			}).
		Route(
			server.Endpoint{
				UrlRoute:   "/v1/driver/",
				HttpMethod: server.GET,
				AuthScope:  server.AuthScopeNone,
			}).
		To(machine.ListDrivers).
		Route(
			server.Endpoint{
				UrlRoute:   "/v1/driver/{driver}/options",
				HttpMethod: server.GET,
				AuthScope:  server.AuthScopeNone,
			}).
		To(machine.DriverOptions).
		Route(
			server.Endpoint{
				UrlRoute:   "/v1/machine/{driver}/host/{name}",
				HttpMethod: server.POST,
				AuthScope:  server.AuthScopeNone,
			}).
		To(machine.CreateInstance).
		Route(
			server.Endpoint{
				UrlRoute:   "/quitquitquit",
				HttpMethod: server.POST,
				AuthScope:  server.AuthScope("admin-kill"),
			}).
		To(
			func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
				glog.Infoln("Stopping the server....")
				close(shutdown)
			}).
		OnShutdown(
			func() error {
				glog.Infoln("Executing user custom shutdown...")
				return nil
			}).
		Start()

	// For stopping the server
	go func() {
		<-shutdown
		stop <- 1
	}()
	return stopped
}
