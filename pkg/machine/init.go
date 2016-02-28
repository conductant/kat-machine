package machine

import (
	"github.com/docker/machine/drivers/amazonec2"
	"github.com/docker/machine/drivers/azure"
	"github.com/docker/machine/drivers/digitalocean"
	"github.com/docker/machine/drivers/exoscale"
	"github.com/docker/machine/drivers/generic"
	"github.com/docker/machine/drivers/google"
	"github.com/docker/machine/drivers/hyperv"
	"github.com/docker/machine/drivers/none"
	"github.com/docker/machine/drivers/openstack"
	"github.com/docker/machine/drivers/rackspace"
	"github.com/docker/machine/drivers/softlayer"
	"github.com/docker/machine/drivers/virtualbox"
	"github.com/docker/machine/drivers/vmwarefusion"
	"github.com/docker/machine/drivers/vmwarevcloudair"
	"github.com/docker/machine/drivers/vmwarevsphere"
	"github.com/docker/machine/libmachine/drivers"
)

type factory func(hostName, storePath string) (driverName string, driver drivers.Driver)

var driverFactories = map[string]factory{}

func init() {
	for _, factory := range []factory{
		func(h, p string) (string, drivers.Driver) {
			d := amazonec2.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := azure.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := digitalocean.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := exoscale.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := generic.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := google.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := hyperv.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := none.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := openstack.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := rackspace.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := softlayer.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := virtualbox.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := vmwarefusion.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := vmwarevcloudair.NewDriver(h, p)
			return d.DriverName(), d
		},
		func(h, p string) (string, drivers.Driver) {
			d := vmwarevsphere.NewDriver(h, p)
			return d.DriverName(), d
		},
	} {
		key, _ := factory("", "")
		driverFactories[key] = factory
	}
}
