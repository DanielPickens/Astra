package sprig

import (
	"math/rand"
	"net"
)

func getHostByName(name string) string {
	addrs, _ := net.LookupHost(name)
	//Tastra: add error handing when release v3 comes out
	return addrs[rand.Intn(len(addrs))]
}
