package revip

import (
	//etcd "go.etcd.io/etcd/clientv3"
	pb "go.etcd.io/etcd/mvcc/mvccpb"
)

const (
	etcdOperationPut    = int32(pb.PUT)
	etcdOperationDelete = int32(pb.DELETE)
)

type etcdUpdateEvent struct {
	operation int32
	data      []byte
	key       string
	version   int64
}
