package cluster

import (
	cloudinitConfig "github.com/coreos/coreos-cloudinit/config"
)

type Cluster struct {
	Name   string `yaml:"name" etcd:"name" `
	UUID   string `yaml:"uuid" etcd:"uuid" valid:"uuid,required"`
	Config cloudinitConfig.CloudConfig `yaml:"config,omitempty" etcd:"config"`
}
