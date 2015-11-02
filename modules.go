package main

import (
	_ "github.com/dtmistry/registrator/consul"
	_ "github.com/dtmistry/registrator/consulkv"
	_ "github.com/dtmistry/registrator/etcd"
	_ "github.com/dtmistry/registrator/skydns2"
	_ "github.com/dtmistry/registrator/bigip"
)
