package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type options struct {
	Bind       string
	IpType     ipType
	Short      ipType
	VerboseErr bool
	APIs       []string
	TimeoutSec int
	PrintVer   bool
	PrintUsage bool
}

type ipType int

const (
	ipv4 ipType = iota
	ipv6
	ipUnset
)

func getOptions() (options, error) {
	var opts options

	// Define flags
	flag.StringVar(&opts.Bind, "bind", "", "Specify the binding (interface or IP)")
	flag.Var((*ipTypeFlag)(&opts.IpType), "ipv6", "Set IP type to IPv6")
	flag.Var((*ipTypeFlag)(&opts.IpType), "ipv4", "Set IP type to IPv4")
	flag.Var((*ipTypeFlag)(&opts.Short), "short", "Specify short IP type (ipv4 or ipv6)")
	flag.BoolVar(&opts.VerboseErr, "verbose", false, "Enable verbose error messages")
	flag.Var((*stringSliceFlag)(&opts.APIs), "api", "Add an API")
	flag.Var((*intFlag)(&opts.TimeoutSec), "timeout", "Specify a timeout in seconds")
	flag.BoolVar(&opts.PrintVer, "version", false, "Print version and exit")
	flag.BoolVar(&opts.PrintUsage, "help", false, "Print usage information and exit")

	// Parse the command-line arguments
	flag.Parse()

	// Check for constraints
	err := checkOptionsConstraints(opts)
	if err != nil {
		return options{}, err
	}

	return opts, nil
}

func checkOptionsConstraints(opts options) error {
	// Example of constraints checking
	if opts.Short != ipUnset && (opts.IpType != ipUnset || opts.VerboseErr || opts.TimeoutSec != 0) {
		return fmt.Errorf("Short IP type can't be used with other flags")
	}
	if opts.IpType != ipUnset && (opts.Short != ipUnset || opts.VerboseErr || opts.TimeoutSec != 0) {
		return fmt.Errorf("IP type can't be used with other flags")
	}

	// Add more constraints as needed

	return nil
}

// Custom flag type for ipType
type ipTypeFlag ipType

func (f *ipTypeFlag) Set(value string) error {
	switch strings.ToLower(value) {
	case "ipv4", "4":
		*f = ipTypeFlag(ipv4)
	case "ipv6", "6":
		*f = ipTypeFlag(ipv6)
	default:
		return fmt.Errorf("invalid IP type: %s", value)
	}
	return nil
}

func (f *ipTypeFlag) String() string {
	return fmt.Sprintf("%d", *f)
}

// Custom flag type for string slice
type stringSliceFlag []string

func (f *stringSliceFlag) String() string {
	return strings.Join(*f, ", ")
}

func (f *stringSliceFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

// Custom flag type for int
type intFlag int

func (f *intFlag) String() string {
	return fmt.Sprintf("%d", *f)
}

func (f *intFlag) Set(value string) error {
	val, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = intFlag(val)
	return nil
}

func main() {
	opts, err := getOptions()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	// Use the parsed options as needed
	fmt.Printf("Bind: %s\n", opts.Bind)
	fmt.Printf("IP Type: %v\n", opts.IpType)
	fmt.Printf("Short IP Type: %v\n", opts.Short)
	fmt.Printf("Verbose Error: %v\n", opts.VerboseErr)
	fmt.Printf("APIs: %v\n", opts.APIs)
	fmt.Printf("Timeout (seconds): %v\n", opts.TimeoutSec)
	fmt.Printf("Print Version: %v\n", opts.PrintVer)
	fmt.Printf("Print Usage: %v\n", opts.PrintUsage)
}
