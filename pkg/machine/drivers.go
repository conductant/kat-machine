package machine

import (
	"errors"
	"github.com/conductant/gohm/pkg/server"
	"github.com/docker/machine/libmachine/drivers"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"path"
)

var (
	ErrDriverNotFound = errors.New("err-driver-not-found")
)

type jsonFlags map[string]interface{}

func (f jsonFlags) String(key string) string {
	if v, has := f[key]; has {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (f jsonFlags) StringSlice(key string) []string {
	if v, has := f[key]; has {
		if s, ok := v.([]string); ok {
			return s
		}
	}
	return []string{}
}

func (f jsonFlags) Int(key string) int {
	if v, has := f[key]; has {
		switch v := v.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case uint64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return -1
}

func (f jsonFlags) Bool(key string) bool {
	if v, has := f[key]; has {
		if s, ok := v.(bool); ok {
			return s
		}
	}
	return false
}

func getStorePath(ctx context.Context, hostName string) string {
	// TOOD - Allow user to set this and storePath be loaded based on user
	wd, _ := os.Getwd()
	storePath := path.Join(wd, ".machine")

	machinePath := path.Join(storePath, "machines", hostName)
	err := os.MkdirAll(machinePath, 0755)
	if err != nil {
		panic(err)
	}

	return storePath
}

func getDriver(ctx context.Context, name, hostName string) (drivers.Driver, error) {
	factory, ok := driverFactories[name]
	if !ok {
		return nil, ErrDriverNotFound
	} else {
		_, driver := factory(hostName, getStorePath(ctx, hostName))
		return driver, nil
	}
}

func ListDrivers(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	list := []string{}
	for k, _ := range driverFactories {
		list = append(list, k)
	}
	server.Marshal(resp, req, list)
}

func DriverOptions(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	driverName := server.GetUrlParameter(req, "driver")
	driver, err := getDriver(ctx, driverName, "")
	if err != nil {
		server.HandleError(ctx, http.StatusNotFound, "not-found:"+driverName)
		return
	} else {
		server.Marshal(resp, req, driver.GetCreateFlags())
	}

}
