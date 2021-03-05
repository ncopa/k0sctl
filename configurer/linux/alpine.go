package linux

import (
	"strings"

	"github.com/k0sproject/k0sctl/configurer"
	"github.com/k0sproject/rig"
	"github.com/k0sproject/rig/os"
	"github.com/k0sproject/rig/os/registry"
)

// RigLinux is a reference to rig's linux struct, renamed to not overlap with configurer.Linux
type RigLinux struct {
	os.Linux
}

// Alpine provides OS support for Alpine Linux
type Alpine struct {
	RigLinux
	configurer.Linux
}

func init() {
	registry.RegisterOSModule(
		func(os rig.OSVersion) bool {
			return os.ID == "alpine"
		},
		func() interface{} {
			return &Alpine{}
		},
	)
}

// InstallPackage installs packages via slackpkg
func (l Alpine) InstallPackage(h os.Host, pkg ...string) error {
	return h.Execf("sudo apk add -U %s", strings.Join(pkg, " "))
}

// InstallKubectl installs kubectl using the gcloud kubernetes repo
func (l Alpine) InstallKubectl(h os.Host) error {
	return l.InstallPackage(h, "--repository https://dl-cdn.alpinelinux.org/alpine/edge/testing kubectl")
}
