package machine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conductant/gohm/pkg/server"
	"github.com/docker/machine/libmachine/drivers"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
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

func getStorePath(ctx context.Context, provider string) string {
	// TOOD - Allow user to set this and storePath be loaded based on user
	wd, _ := os.Getwd()
	storePath := path.Join(wd, ".machine", provider)
	err := os.MkdirAll(storePath, 0755)
	if err != nil {
		panic(err)
	}
	return storePath
}

func getMachinePath(ctx context.Context, provider, hostName string) string {
	machinePath := path.Join(getStorePath(ctx, provider), "machines", hostName)
	err := os.MkdirAll(machinePath, 0755)
	if err != nil {
		panic(err)
	}
	return machinePath
}

func getMachineLogPath(ctx context.Context, provider, hostName string) string {
	logPath := path.Join(getMachinePath(ctx, provider, hostName), "log")
	err := os.MkdirAll(logPath, 0755)
	if err != nil {
		panic(err)
	}
	return logPath
}

func saveDriver(ctx context.Context, driver drivers.Driver, operation, hostName string) error {
	state, err := json.Marshal(driver)
	if err != nil {
		return err
	}

	p := path.Join(getMachineLogPath(ctx, driver.DriverName(), hostName),
		fmt.Sprintf("%d-%s.json", time.Now().Unix(), operation))
	err = ioutil.WriteFile(p, state, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getLastState(ctx context.Context, provider, hostName string) ([]byte, error) {
	logs := getMachineLogPath(ctx, provider, hostName)
	list, err := ioutil.ReadDir(logs)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		last := path.Join(logs, list[len(list)-1].Name())
		return ioutil.ReadFile(last)
	}
	return []byte{}, nil
}

func getDriver(ctx context.Context, provider, hostName string) (driver drivers.Driver, restored bool, err error) {
	factory, ok := driverFactories[provider]
	if !ok {
		return nil, false, ErrDriverNotFound
	} else {
		_, driver = factory(hostName, getStorePath(ctx, provider))

		lastState, err := getLastState(ctx, provider, hostName)
		if err != nil {
			return nil, false, err
		}
		if len(lastState) > 0 {
			if err := json.Unmarshal(lastState, driver); err == nil {
				restored = true
			}
		}
	}
	return
}

func driverToJSON(driver drivers.Driver) string {
	buff, _ := json.MarshalIndent(driver, " ", " ")
	return string(buff)
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
	driver, _, err := getDriver(ctx, driverName, "")
	if err != nil {
		server.HandleError(ctx, http.StatusNotFound, "not-found:"+driverName)
		return
	} else {
		server.Marshal(resp, req, driver.GetCreateFlags())
	}

}
