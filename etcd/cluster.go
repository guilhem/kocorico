package etcd

import (
	"context"
	"reflect"

	"github.com/guilhem/kocorico/cluster"

	"github.com/mickep76/etcdmap"
)

const clusterpath = "clusters"

func GetCluster(clusterUUID string) (cluster.Cluster, error) {
	resp, err := Kapi.Get(context.Background(), clusterpath, nil)
	if err != nil {
		return cluster.Cluster{}, err
	}
	c := cluster.Cluster{}
	if err := etcdmap.Struct(resp.Node, reflect.ValueOf(&c)); err != nil {
		return cluster.Cluster{}, err
	}

	return c, nil
}

func CreateCluster(c cluster.Cluster) error {
	return etcdmap.Create(Kapi, clusterpath, reflect.ValueOf(c))
}
