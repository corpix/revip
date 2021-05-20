package main

import (
	"fmt"

	"github.com/corpix/revip"

	"github.com/davecgh/go-spew/spew"
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

func (x *X) Update(xx interface{}) error {
	spew.Dump("UPDATE", x, xx)
	return nil
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
	xx := X{}

	//

	url := "etcd://127.0.0.1:2379/" + prefix
	c, err := revip.NewEtcdClient(url)
	if err != nil {
		panic(err)
	}

	//

	to, err := revip.ToURL(url, revip.JsonMarshaler)
	if err != nil {
		panic(err)
	}

	from, err := revip.FromURL(url, revip.JsonUnmarshaler)
	if err != nil {
		panic(err)
	}

	watch := revip.WithUpdatesFromEtcd(c, prefix, revip.JsonUnmarshaler)

	//

	err = to(x)
	if err != nil {
		panic(err)
	}
	err = from(&xx)
	if err != nil {
		panic(err)
	}

	//

	err = watch(&x)
	if err != nil {
		panic(err)
	}

	spew.Dump("TO ETCD", x)
	spew.Dump("FROM ETCD", xx)

	fmt.Println()
	fmt.Println(`now try to run: etcdctl put test/Foo '"hello you"'`)

	select {}
}
