package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/nixigaj/go-default-route"
)

var apis = [...]string{
	"ip.erix.dev",
	"icanhazip.com",
	"ipecho.net/plain",
	"ifconfig.me/ip",
	"api64.ipify.org",
}

const (
	defaultBind = "0.0.0.0:0"
)

type ipType int

const (
	ipv4 ipType = iota
	ipv6
	ipInvalid
)

func main() {
	fmt.Println("Hello!")
}

func getIpType(str string) ipType {
	ip := net.ParseIP(str)
	if ip == nil {
		return ipInvalid
	}
	if ip.To4() == nil {
		return ipv6
	}
	return ipv4
}

func getAddr() (*net.TCPAddr, error) {
	defaultRoute, err := defaultroute.DefaultRouteInterface()
}

func getInterfaceAddr(interf *net.Interface, ipType ipType) (*net.TCPAddr, error) {

	addrs, err := interf.Addrs()
	if err != nil {
		return nil, err
	}

	/*
		if ipType
		for addr := range addrs {

		}
	*/

	return addrs[0], nil
}

func getHttpClient(addr *net.TCPAddr) (*http.Client, error) {
	dialer := &net.Dialer{
		LocalAddr: addr,
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := dialer.Dial(network, addr)
		return conn, err
	}

	transport := &http.Transport{DialContext: dialContext}
	client := &http.Client{
		Transport: transport,
	}

	return client, nil
}
