package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var defaultAPIs = [...]string{
	"ip.erix.dev",
	"icanhazip.com",
	"ipecho.net/plain",
	"ifconfig.me/ip",
	"api64.ipify.org",
}

type ipType int

const (
	ipv4 ipType = iota
	ipv6
	ipInvalid
)

func main() {
	startTime := time.Now()

	output := ""

	ipv4Str, err := getIPString(ipv4, "")
	if err != nil {
		output += "IPv4: failed to get IPv4 address: " + err.Error() + "\n"
	} else {
		output += "IPv4: " + ipv4Str + "\n"
	}

	ipv6Str, err := getIPString(ipv6, "")
	if err != nil {
		output += "IPv6: failed to get IPv6 address: " + err.Error() + "\n"
	} else {
		output += "IPv6: " + ipv6Str + "\n"
	}

	endTime := time.Since(startTime)

	output += "Query time: " + strconv.Itoa(int(endTime.Milliseconds())) + " ms"

	fmt.Println(output)
}

func getIPString(ipType ipType, bind string) (string, error) {
	clients, err := getHTTPClients(ipType, bind)
	if err != nil {
		return "", err
	}

	resp, err := clients[0].Get("https://" + defaultAPIs[0])
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ipStr := strings.TrimSpace(string(body))

	if getIpType(ipStr) != ipType {
		return "", fmt.Errorf("response IP not correct")
	}

	return ipStr, nil
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

func getHTTPClients(ipType ipType, bind string) ([]*http.Client, error) {
	var clients []*http.Client

	for range defaultAPIs {
		client, err := getHTTPClient(ipType, bind)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func getHTTPClient(ipType ipType, bind string) (*http.Client, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	var transport *http.Transport

	if bind == "" {
		// Bind to default address with ipType

		var network string
		switch ipType {
		case ipv4:
			network = "tcp4"
		case ipv6:
			network = "tcp6"
		default:
			return nil, fmt.Errorf("invalid IP type")
		}

		transport = http.DefaultTransport.(*http.Transport).Clone() // Remove the 'var' keyword here

		var dialer net.Dialer
		transport.DialContext = func(ctx context.Context, _, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		}
	} else {
		// Bind to a specific address

		addr, err := net.ResolveTCPAddr("tcp", bind)
		if err != nil {
			return nil, err
		}

		dialer := &net.Dialer{LocalAddr: addr}
		transport = &http.Transport{} // Initialize the transport
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := dialer.Dial(network, addr)
			return conn, err
		}
	}

	client.Transport = transport

	return client, nil
}
