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
	"https://api64.ipify.org",
	"https://icanhazip.com",
	"https://ifconfig.me/ip",
	"https://ip.erix.dev",
	"https://ipecho.net/plain",
}

const (
	whatVersion = "0.1.0"

	// With multiple APIs, it is unlikely that the query will take longer than one second
	defaultClientTimeoutSec = 1
)

type ipType int

const (
	ipv4 ipType = iota
	ipv6
	ipUnset
)

type bindType struct {
	IP    ipType
	iFace bool
}

type ipStringResp struct {
	IP        string
	QueryTime time.Duration
	Err       error
}

type options struct {
	Bind       string
	BindType   bindType
	Short      ipType
	VerboseErr bool
	APIs       []string
	Timeout    time.Duration
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
		out, _ := getSingleOutput(opt.BindType, opt.Bind, opt.APIs, opt.VerboseErr, opt.Timeout)
		fmt.Print(out)
	} else {
		out := getPrettyOutput(opt.BindType, opt.Bind, opt.APIs, opt.VerboseErr, opt.Timeout)
		fmt.Println(out)
	}

	os.Exit(0)
}

// TODO: Implement flag parsing to options struct
func getOptions() (options, error) {
	opt := options{
		Bind:       "2a0e:1c80:1337:1:10:0:52:85",
		Short:      ipUnset,
		VerboseErr: false,
		APIs:       defaultAPIs,
		Timeout:    defaultClientTimeoutSec * time.Second,
		PrintVer:   false,
		PrintUsage: false,
	}
	opt.BindType = getBindType(opt.Bind)

	//opt.BindType.IP = ipUnset

	return opt, nil
}

func getBindType(str string) bindType {
	if str == "" {
		return bindType{ipUnset, false}
	}

	// TODO: This also accepts IPv4s with square brackets
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

func getSingleOutput(bType bindType, bind string, apis []string, verboseErr bool, timeout time.Duration) (string, time.Duration) {
	respChan := make(chan ipStringResp)
	var wg sync.WaitGroup

	wg.Add(1)
	go getIPString(respChan, apis, bType, bind, timeout, &wg)
	resp := <-respChan
	close(respChan)

	if resp.Err != nil {
		output := "failed to get " + getIPTypeStr(bType) + " address"
		if verboseErr {
			output += ": " + resp.Err.Error()
		}
		return output, resp.QueryTime
	}

	wg.Wait()
	return resp.IP, resp.QueryTime
}

func getPrettyOutput(bType bindType, bind string, apis []string, verboseErr bool, timeout time.Duration) string {
	var output string
	var longestQueryTime time.Duration
	var wg sync.WaitGroup

	if bType.IP != ipUnset {
		var singleOut string
		singleOut, longestQueryTime = getSingleOutput(bType, bind, apis, verboseErr, timeout)

		output += getIPTypeStr(bType) + ": " + singleOut + "\n"
	} else {
		respChanV4 := make(chan ipStringResp)
		respChanV6 := make(chan ipStringResp)

		wg.Add(2)
		go getIPString(respChanV4, apis, bindType{IP: ipv4, iFace: bType.iFace}, bind, timeout, &wg)
		go getIPString(respChanV6, apis, bindType{IP: ipv6, iFace: bType.iFace}, bind, timeout, &wg)

		respV4 := <-respChanV4
		respV6 := <-respChanV6

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

		if respV4.QueryTime > respV6.QueryTime {
			longestQueryTime = respV4.QueryTime
		} else {
			longestQueryTime = respV6.QueryTime
		}
	}

	output += "Query time: " + strconv.Itoa(int(longestQueryTime.Milliseconds())) + " ms"

	wg.Wait()

	return output
}

func getIPString(respChan chan ipStringResp, apis []string, bType bindType, bind string, timeout time.Duration, wg *sync.WaitGroup) {
	clients, err := getHTTPClients(len(apis), bType, bind, timeout)
	if err != nil {
		respChan <- ipStringResp{
			IP:        "",
			QueryTime: 0,
			Err:       err,
		}
		wg.Done()
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
	longestQueryTime := time.Duration(0)
	for range clients {
		resp := <-clientRespChan

		if success {
			continue
		}

		if resp.Err == nil {
			success = true
			respChan <- ipStringResp{
				IP:        resp.IP,
				QueryTime: resp.QueryTime,
				Err:       nil,
			}
			cancel()
		}
		if resp.QueryTime > longestQueryTime {
			longestQueryTime = resp.QueryTime
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
		IP:        "",
		QueryTime: longestQueryTime,
		Err:       errors.New(errsString),
	}

	wg.Done()
}

func getHTTPClients(noClients int, bType bindType, bind string, timeout time.Duration) ([]*http.Client, error) {
	var clients []*http.Client

	for i := 0; i < noClients; i++ {
		client, err := getHTTPClient(bType, bind, timeout)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func getHTTPClient(bType bindType, bind string, timeout time.Duration) (*http.Client, error) {
	client := &http.Client{
		Timeout: timeout,
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

func trimSubnet(ipStr string) string {
	slashPos := strings.Index(ipStr, "/")

	if slashPos != -1 {
		return (ipStr)[:slashPos]
	}

	return ipStr
}

func fetchIP(respChan chan ipStringResp, client *http.Client, ctx context.Context, api string, ipType ipType) {
	req, err := http.NewRequestWithContext(ctx, "GET", api, nil)
	if err != nil {
		respChan <- ipStringResp{"", 0, err}
		return
	}

	queryStartTime := time.Now()
	resp, err := client.Do(req)
	queryTime := time.Since(queryStartTime)
	if err != nil {
		respChan <- ipStringResp{"", queryTime, err}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respChan <- ipStringResp{"", queryTime, err}
		return
	}

	ipStr := strings.TrimSpace(string(body))

	if getBindType(ipStr).IP != ipType {
		respChan <- ipStringResp{
			IP:        "",
			QueryTime: queryTime,
			Err:       fmt.Errorf("response IP not correct"),
		}
		return
	}

	respChan <- ipStringResp{
		IP:        ipStr,
		QueryTime: queryTime,
		Err:       nil,
	}
}
