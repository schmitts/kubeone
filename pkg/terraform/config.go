package terraform

import (
	"encoding/json"
	"strconv"

	"github.com/kubermatic/kubeone/pkg/manifest"
)

// Config represents configuration in the terraform output format
type Config struct {
	KubeOneAPI struct {
		Value struct {
			Endpoint string `json:"endpoint"`
		} `json:"value"`
	} `json:"kubeone_api"`

	KubeOneHosts struct {
		Value struct {
			ControlPlane []struct {
				PublicAddress     []string `json:"public_address"`
				PrivateAddress    []string `json:"private_address"`
				SSHUser           string   `json:"ssh_user"`
				SSHPort           string   `json:"ssh_port"`
				SSHPrivateKeyFile string   `json:"ssh_private_key_file"`
				SSHAgentSocket    string   `json:"ssh_agent_socket"`
			} `json:"control_plane"`
		} `json:"value"`
	} `json:"kubeone_hosts"`
}

// NewConfigFromJSON creates a new config object from json
func NewConfigFromJSON(j []byte) (c *Config, err error) {
	c = &Config{}
	return c, json.Unmarshal(j, c)
}

// Apply adds the terraform configuration options to the given manifest
func (c Config) Apply(m *manifest.Manifest) {
	if c.KubeOneAPI.Value.Endpoint != "" {
		m.APIServer.Address = c.KubeOneAPI.Value.Endpoint
	}

	var hosts []manifest.HostManifest
	cp := c.KubeOneHosts.Value.ControlPlane[0]
	sshPort, _ := strconv.Atoi(cp.SSHPort)

	privateIPs := cp.PrivateAddress

	for i, publicIP := range cp.PublicAddress {
		privateIP := publicIP
		if i < len(privateIPs) {
			privateIP = privateIPs[i]
		}

		hosts = append(hosts, manifest.HostManifest{
			PublicAddress:     publicIP,
			PrivateAddress:    privateIP,
			SSHUsername:       cp.SSHUser,
			SSHPort:           sshPort,
			SSHPrivateKeyFile: cp.SSHPrivateKeyFile,
			SSHAgentSocket:    cp.SSHAgentSocket,
		})
	}

	if len(hosts) > 0 {
		m.Hosts = hosts
	}
}
