package etcd

import (
	etcdClient "github.com/coreos/etcd/client"
)

var (
	// Kapi is an etcd client tu use
	Kapi etcdClient.KeysAPI
)
