package phase

import (
	"strings"

	"github.com/k0sproject/k0sctl/config/cluster"
	log "github.com/sirupsen/logrus"
)

// PrepareHosts installs required packages and so on on the hosts.
type PrepareHosts struct {
	GenericPhase
}

// Title for the phase
func (p *PrepareHosts) Title() string {
	return "Prepare hosts"
}

// Run the phase
func (p *PrepareHosts) Run() error {
	return p.Config.Spec.Hosts.ParallelEach(p.prepareHost)
}

func (p *PrepareHosts) prepareHost(h *cluster.Host) error {
	if len(h.Environment) > 0 {
		log.Infof("%s: updating environment", h)
		if err := h.Configurer.UpdateEnvironment(h, h.Environment); err != nil {
			return err
		}
	}

	var pkgs []string

	if h.NeedCurl() {
		pkgs = append(pkgs, "curl")
	}

	if h.NeedIPTables() {
		pkgs = append(pkgs, "iptables")
	}

	if len(pkgs) > 0 {
		log.Infof("%s: installing packages (%s)", h, strings.Join(pkgs, ", "))
		if err := h.Configurer.InstallPackage(h, pkgs...); err != nil {
			return err
		}
	}

	if h.IsController() && !h.Configurer.CommandExist(h, "kubectl") {
		log.Infof("%s: installing kubectl", h)
		if err := h.Configurer.InstallKubectl(h); err != nil {
			return err
		}
	}

	if h.Configurer.IsContainer(h) {
		log.Infof("%s: is a container, applying a fix", h)
		if err := h.Configurer.FixContainer(h); err != nil {
			return err
		}
	}

	return nil
}
