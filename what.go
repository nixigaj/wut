package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var defaultAPIs = []string{
	"https://ip.erix.dev",
	"https://icanhazip.com",
	"https://ipecho.net/plain",
	"https://ifconfig.me/ip",
	"https://api64.ipify.org",
}

const (
	whatVersion      = "0.1.0"
	clientTimeoutSec = 5
)

type bindType struct {
	IP    ipType
	iFace bool
}

type ipType int

const (
	ipv4 ipType = iota
	ipv6
	ipUnset
)

type ipStringResp struct {
	IP  string
	Err error
}

type options struct {
	Bind       string
	BindType   bindType
	Short      ipType
	VerboseErr bool
	APIs       []string
	PrintVer   bool
	PrintUsage bool
}

func main() {
	opt, err := getOptions()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	if opt.PrintUsage {
		flag.Usage()
		os.Exit(0)
	}

	if opt.PrintVer {
		fmt.Println("`what` version " + whatVersion)
		os.Exit(0)
	}

	if opt.Short != ipUnset {
		fmt.Print(getSingleOutput(opt.BindType, opt.Bind, opt.APIs, opt.VerboseErr))
	} else {
		fmt.Println(getPrettyOutput(opt.BindType, opt.Bind, opt.APIs, opt.VerboseErr))
	}

	os.Exit(0)
}

func getOptions() (options, error) {
	opt := options{
		Bind:       "",
		Short:      ipUnset,
		VerboseErr: false,
		APIs:       defaultAPIs,
		PrintVer:   false,
		PrintUsage: false,
	}
	opt.BindType = getBindType(opt.Bind)

	opt.BindType.IP = ipUnset

	return opt, nil
}

func getPrettyOutput(bType bindType, bind string, apis []string, verboseErr bool) string {
	var output string
	var queryTime time.Duration
	var wg sync.WaitGroup

	queryStartTime := time.Now()

	if bType.IP != ipUnset {
		singleOut := getSingleOutput(bType, bind, apis, verboseErr)
		queryTime = time.Since(queryStartTime)

		output += getIPTypeStr(bType) + ": " + singleOut + "\n"
	} else {
		respChanV4 := make(chan ipStringResp)
		respChanV6 := make(chan ipStringResp)

		wg.Add(2)
		go getIPString(respChanV4, apis, bindType{IP: ipv4, iFace: bType.iFace}, bind, &wg)
		go getIPString(respChanV6, apis, bindType{IP: ipv6, iFace: bType.iFace}, bind, &wg)

		respV4 := <-respChanV4
		respV6 := <-respChanV6

		queryTime = time.Since(queryStartTime)

		close(respChanV4)
		close(respChanV6)

		output += "IPv4: "
		if respV4.Err != nil {
			output += "failed to get IPv4 address"
			if verboseErr {
				output += ": " + respV4.Err.Error()
			}
		} else {
			output += respV4.IP
		}
		output += "\n"

		output += "IPv6: "
		if respV6.Err != nil {
			output += "failed to get IPv6 address"
			if verboseErr {
				output += ": " + respV6.Err.Error()
			}
		} else {
			output += respV6.IP
		}
		output += "\n"
	}

	output += "Query time: " + strconv.Itoa(int(queryTime.Milliseconds())) + " ms"

	wg.Wait()

	return output
}

func getIPTypeStr(bType bindType) string {
	switch bType.IP {
	case ipv4:
		return "IPv4"
	case ipv6:
		return "IPv6"
	default:
		return "IP"
	}
}

func trimSubnet(ipStr string) string {
	slashPos := strings.Index(ipStr, "/")

	if slashPos != -1 {
		return (ipStr)[:slashPos]
	}

	return ipStr
}

func getInterfaceIP(bType bindType, ifaceName string) (string, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return "", err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipAddr := addr.(*net.IPNet).IP

		switch bType.IP {
		case ipv4:
			if ipAddr.To4() != nil {
				return trimSubnet(addr.String()), nil
			}
		case ipv6:
			if ipAddr.To4() == nil {
				return trimSubnet(addr.String()), nil
			}
		default:
			return "", fmt.Errorf("invalid bind IP type")
		}
	}

	return "", errors.New(
		fmt.Sprintf("interface %s does not have an %s address\n", ifaceName, getIPTypeStr(bType)))
}

func getSingleOutput(bType bindType, bind string, apis []string, verboseErr bool) string {
	respChan := make(chan ipStringResp)
	var wg sync.WaitGroup

	wg.Add(1)
	go getIPString(respChan, apis, bType, bind, &wg)
	resp := <-respChan
	close(respChan)

	if resp.Err != nil {
		output := "failed to get " + getIPTypeStr(bType) + " address"
		if verboseErr {
			output += ": " + resp.Err.Error()
		}
		return output
	}

	wg.Wait()
	return resp.IP
}

func fetchIP(respChan chan ipStringResp, client *http.Client, ctx context.Context, api string, ipType ipType) {
	req, err := http.NewRequestWithContext(ctx, "GET", api, nil)
	if err != nil {
		respChan <- ipStringResp{"", err}
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		respChan <- ipStringResp{"", err}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respChan <- ipStringResp{
			IP:  "",
			Err: err,
		}
		return
	}

	ipStr := strings.TrimSpace(string(body))

	if getBindType(ipStr).IP != ipType {
		respChan <- ipStringResp{
			IP:  "",
			Err: fmt.Errorf("response IP not correct"),
		}
		return
	}

	respChan <- ipStringResp{
		IP:  ipStr,
		Err: nil,
	}
}

func getIPString(respChan chan ipStringResp, apis []string, bType bindType, bind string, wg *sync.WaitGroup) {
	clients, err := getHTTPClients(len(apis), bType, bind)
	if err != nil {
		respChan <- ipStringResp{
			IP:  "",
			Err: err,
		}
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	clientRespChan := make(chan ipStringResp)

	for i, client := range clients {
		go fetchIP(clientRespChan, client, ctx, apis[i], bType.IP)
	}

	var errs []error
	success := false
	for range clients {
		resp := <-clientRespChan

		if success {
			continue
		}

		if resp.Err == nil {
			success = true
			respChan <- ipStringResp{
				IP:  resp.IP,
				Err: nil,
			}
			cancel()
		}

		errs = append(errs, resp.Err)
	}
	close(clientRespChan)
	cancel()

	if success {
		wg.Done()
		return
	}

	var errsString string
	for i, err := range errs {
		errsString += err.Error()

		if i < len(errs)-1 {
			errsString += ", "
		}
	}

	respChan <- ipStringResp{
		IP:  "",
		Err: errors.New(errsString),
	}

	wg.Done()
}

func getBindType(str string) bindType {
	if str == "" {
		return bindType{ipUnset, false}
	}

	str = strings.Trim(str, "[]") // Trim for IPv6
	ip := net.ParseIP(str)

	if ip == nil {
		return bindType{ipUnset, true}
	}
	if ip.To4() == nil {
		return bindType{ipv6, false}
	}
	return bindType{ipv4, false}
}

func getHTTPClients(noClients int, bType bindType, bind string) ([]*http.Client, error) {
	var clients []*http.Client

	for i := 0; i < noClients; i++ {
		client, err := getHTTPClient(bType, bind)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func getHTTPClient(bType bindType, bind string) (*http.Client, error) {
	client := &http.Client{
		Timeout: clientTimeoutSec * time.Second,
	}
	var transport *http.Transport

	if bind == "" {
		// Bind to default address with bindType

		var network string
		switch bType.IP {
		case ipv4:
			network = "tcp4"
		case ipv6:
			network = "tcp6"
		default:
			return nil, fmt.Errorf("invalid IP type")
		}

		transport = http.DefaultTransport.(*http.Transport).Clone()

		var dialer net.Dialer
		transport.DialContext = func(ctx context.Context, _, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		}
	} else {
		// Bind to a specific address

		var bindIP string
		if bType.iFace {
			var err error
			bindIP, err = getInterfaceIP(bType, bind)
			if err != nil {
				return nil, err
			}
		} else {
			bindIP = strings.Trim(bind, "[]") // Trim for IPv6
		}

		// Add braces if IPv6
		if bType.IP == ipv6 {
			bindIP = "[" + bindIP + "]"
		}

		addr, err := net.ResolveTCPAddr("tcp", bindIP+":0")
		if err != nil {
			return nil, err
		}

		dialer := &net.Dialer{LocalAddr: addr}
		transport = &http.Transport{}
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := dialer.Dial(network, addr)
			return conn, err
		}
	}

	client.Transport = transport

	return client, nil
}
