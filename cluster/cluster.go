package cluster

import (
	cloudinitConfig "github.com/coreos/coreos-cloudinit/config"
)

type Cluster struct {
	Name   string
	UUID   string
	Config cloudinitConfig.CloudConfig
}
