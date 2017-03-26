package etcd

import (
	"context"
	"path"
	"reflect"

	"github.com/guilhem/kocorico/cluster"

	"github.com/mitchellh/mapstructure"
	etcdClient "github.com/coreos/etcd/client"
	"github.com/mickep76/etcdmap"
)

const clusterpath = "/clusters"

func GetCluster(clusterUUID string) (cluster.Cluster, error) {
	resp, err := Kapi.Get(context.Background(), path.Join(clusterpath, clusterUUID), &etcdClient.GetOptions{Recursive: true})
	var c cluster.Cluster
	if err != nil {
		return c, err
	}
	m := etcdmap.Map(resp.Node)
	err = mapstructure.Decode(m, &c)
	return c, err
}

func CreateCluster(c cluster.Cluster) error {
	return etcdmap.Create(Kapi, path.Join(clusterpath, c.UUID), reflect.ValueOf(c))
}
