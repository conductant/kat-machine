package machine

import (
	"github.com/conductant/gohm/pkg/server"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"path"
)

func ListAllHosts(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	p := getStoreRoot(ctx)
	drivers, err := ioutil.ReadDir(p)
	if err != nil {
		server.HandleError(ctx, http.StatusNotFound, "not-found:"+p)
		return
	}
	result := map[string][]string{}
	for _, driver := range drivers {
		list := []string{}
		visitDir(path.Join(getStorePath(ctx, driver.Name()), "machines"), func(e string) { list = append(list, e) })
		result[driver.Name()] = list
	}
	server.Marshal(resp, req, result)
}

func ListAllHostsByDriver(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	driver := server.GetUrlParameter(req, "driver")
	p := path.Join(getStorePath(ctx, driver), "machines")
	hosts := []string{}
	visitDir(p, func(e string) {
		hosts = append(hosts, e)
	})
	server.Marshal(resp, req, hosts)
}

func visitDir(path string, visit func(string)) {
	list, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, entry := range list {
		visit(entry.Name())
	}
}
