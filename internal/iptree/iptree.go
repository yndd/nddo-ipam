package iptree

import "github.com/k-sone/critbitgo"

type IpTree struct {
	NetTree *critbitgo.Net
}

func NewIpTree() *IpTree {
	return &IpTree{
		NetTree: critbitgo.NewNet(),
	}
}
