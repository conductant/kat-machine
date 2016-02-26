package main

import (
	"fmt"
	"github.com/conductant/gohm/pkg/auth"
	"github.com/conductant/gohm/pkg/command"
	"github.com/conductant/gohm/pkg/resource"
	"github.com/conductant/gohm/pkg/runtime"
	"github.com/conductant/kat-machine/pkg/server"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"io"
	"time"
)

type run struct {
	server.Server
}

func (r *run) Help(w io.Writer) {
	fmt.Fprintln(w, "Run the kat-machine server")
}

func (r *run) Run(args []string, w io.Writer) error {
	glog.Infoln("Starting server:", r.Server)
	err := r.Server.Init()
	if err != nil {
		panic(err)
	}

	stopped := r.Server.Start()
	<-stopped
	glog.Infoln("Bye")
	return nil
}

func (t *run) Close() error {
	return nil
}

type token struct {
	PrivateKeyUrl string        `flag:"private_key_url,The url to private key"`
	Scopes        []string      `flag:"scope, The auth scope"`
	Ttl           time.Duration `flag:"ttl,The token ttl to expiration"`
}

func (t *token) Help(w io.Writer) {
	fmt.Fprintln(w, "Generates auth token for accessing the server.")
}

func (t *token) Run(args []string, w io.Writer) error {
	buff, err := resource.Fetch(context.Background(), t.PrivateKeyUrl)
	if err != nil {
		return err
	}
	token := auth.NewToken(t.Ttl)
	for _, scope := range t.Scopes {
		token.Add(scope, 1)
	}
	signed, err := token.SignedString(func() []byte { return buff })
	if err != nil {
		return err
	}
	fmt.Print(signed)
	return nil
}

func (t *token) Close() error {
	return nil
}

func main() {

	command.Register("run", func() (command.Module, command.ErrorHandling) {
		return new(run), command.PanicOnError
	})
	command.Register("token", func() (command.Module, command.ErrorHandling) {
		return new(token), command.PanicOnError
	})

	runtime.Main()
}
