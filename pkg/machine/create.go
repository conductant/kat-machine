package machine

import (
	"github.com/conductant/gohm/pkg/server"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"net/http"
)

func CreateInstance(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	driverName := server.GetUrlParameter(req, "driver")
	hostName := server.GetUrlParameter(req, "name")

	driver, err := getDriver(ctx, driverName, hostName)
	if err != nil {
		server.HandleError(ctx, http.StatusNotFound, "err-not-found:"+driverName)
		return
	}

	input := jsonFlags{}
	err = server.Unmarshal(resp, req, &input)
	if err != nil {
		return
	}
	// Set default values from the flag definitions
	for _, flag := range driver.GetCreateFlags() {
		if _, has := input[flag.String()]; !has {
			input[flag.String()] = flag.Default()
		}
	}
	err = driver.SetConfigFromFlags(input)
	if err != nil {
		glog.Warningln("Err=", err)
		server.HandleError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	glog.Infoln("Driver=", driver)
	err = driver.Create()
	if err != nil {
		glog.Warningln("Err=", err)
		server.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
}
