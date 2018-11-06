package ssh

import (
	"time"

	"github.com/kubermatic/kubeone/pkg/manifest"
)

type Connector struct {
	connections map[string]Connection
}

func NewConnector() *Connector {
	return &Connector{
		connections: make(map[string]Connection),
	}
}

func (c *Connector) Connect(node manifest.HostManifest) (Connection, error) {
	var err error

	conn, exists := c.connections[node.PublicAddress]
	if !exists || conn.Closed() {
		opts := Opts{
			Username:    node.SSHUsername,
			Port:        node.SSHPort,
			Hostname:    node.PublicAddress,
			KeyFile:     node.SSHPrivateKeyFile,
			AgentSocket: node.SSHAgentSocket,
			Timeout:     10 * time.Second,
		}

		conn, err = NewConnection(opts)
		if err != nil {
			return nil, err
		}

		c.connections[node.PublicAddress] = conn
	}

	return conn, nil
}

func (c *Connector) CloseAll() {
	for _, conn := range c.connections {
		conn.Close()
	}
}
