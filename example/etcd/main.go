package main

import (
	"github.com/corpix/revip"

	"github.com/davecgh/go-spew/spew"
	etcd "go.etcd.io/etcd/clientv3"
)

type X struct {
	Foo string
	Bar map[string]int
	Baz []string
	Qux struct{ QuxFoo string }
	XXX map[string]struct{ XXXFoo string }
	YYY struct {
		ZZZ string
		XXX *string
		YYY *struct{ QQQ string }
	}
}

func (x *X) Update(xx interface{}) {
	spew.Dump("UPDATE", x, xx)
}

func main() {
	prefix := "test"
	x := X{
		Foo: "hello",
		Bar: map[string]int{"doom666": 666},
		Baz: []string{"has", "come"},
		Qux: struct{ QuxFoo string }{QuxFoo: "I was awaiting you"},
		XXX: map[string]struct{ XXXFoo string }{"test": {XXXFoo: "test"}},
	}

	c, err := etcd.New(etcd.Config{Endpoints: []string{"127.0.0.1:2379"}})
	if err != nil {
		panic(err)
	}

	err = revip.ToEtcd(c, prefix, revip.JsonMarshaler)(x)
	if err != nil {
		panic(err)
	}

	xx := X{}
	err = revip.FromEtcd(c, prefix, revip.JsonUnmarshaler)(&xx)
	if err != nil {
		panic(err)
	}

	watch := revip.WithUpdatesFromEtcd(c, prefix, revip.JsonUnmarshaler)
	err = watch(&x)
	if err != nil {
		panic(err)
	}

	spew.Dump("TO ETCD", x)
	spew.Dump("FROM ETCD", xx)

	select {}
}
