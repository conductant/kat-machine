package machine

import (
	"github.com/conductant/gohm/pkg/server"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"net/http"
)

func loadDriver(ctx context.Context, resp http.ResponseWriter, req *http.Request) (string, drivers.Driver, error) {
	driverName := server.GetUrlParameter(req, "driver")
	hostName := server.GetUrlParameter(req, "name")

	driver, restored, err := getDriver(ctx, driverName, hostName)
	if err != nil {
		glog.Warningln("Err=", err)
		server.HandleError(ctx, http.StatusNotFound, "err-not-found:"+driverName)
		return "", nil, err
	}

	// If this driver instance is not restored from persistent store, then initialize
	// it with the defaults from the flags and from the http post input.
	if !restored {
		input := jsonFlags{}
		// Set default values from the flag definitions
		for _, flag := range driver.GetCreateFlags() {
			input[flag.String()] = flag.Default()
		}
		// Unmarshal over the defaults
		err = server.Unmarshal(resp, req, &input)
		if err != nil {
			server.HandleError(ctx, http.StatusBadRequest, err.Error())
			return "", nil, err
		}
		err = driver.SetConfigFromFlags(input)
	}

	glog.Infoln("DRIVER=", driverToJSON(driver))

	if err != nil {
		server.HandleError(ctx, http.StatusBadRequest, err.Error())
		return "", nil, err
	}
	return hostName, driver, nil
}

func CreateInstance(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	hostName, driver, err := loadDriver(ctx, resp, req)
	if err != nil {
		return
	}

	err = driver.Create()
	if err != nil {
		server.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Store the state of the driver so that in future calls we can rebuild the driver
	// and make changes accordingly.  For example the driver can have specific instance id
	// required by the provider's api for start / stop / terminate, etc.
	err = saveDriver(ctx, driver, "create", hostName)
	if err != nil {
		server.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
}

func GetInstanceState(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	hostName, driver, err := loadDriver(ctx, resp, req)
	if err != nil {
		return
	}

	state, err := driver.GetState()
	if err != nil {
		server.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]interface{}{
		"name":  hostName,
		"state": state.String(),
	}
	server.Marshal(resp, req, result)
}

func PutInstanceState(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	hostName, driver, err := loadDriver(ctx, resp, req)
	if err != nil {
		return
	}

	action := server.GetUrlParameter(req, "action")
	switch action {
	case "start":
		err = driver.Start()
	case "stop":
		err = driver.Stop()
	case "restart":
		err = driver.Restart()
	case "kill":
		err = driver.Kill()
	default:
		server.HandleError(ctx, http.StatusBadRequest, "err-unknown-action:"+action)
		return
	}
	if err != nil {
		server.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = saveDriver(ctx, driver, action, hostName)
	if err != nil {
		server.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	newState, err := driver.GetState()
	if err != nil {
		server.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result := map[string]interface{}{
		"name":  hostName,
		"state": newState.String(),
	}
	server.Marshal(resp, req, result)
}

func RemoveInstance(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	hostName, driver, err := loadDriver(ctx, resp, req)
	if err != nil {
		return
	}

	err = driver.Remove()
	if err != nil {
		server.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = saveDriver(ctx, driver, "remove", hostName)
	if err != nil {
		server.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	newState, err := driver.GetState()
	if err != nil {
		server.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	result := map[string]interface{}{
		"name":  hostName,
		"state": newState.String(),
	}
	server.Marshal(resp, req, result)
}
