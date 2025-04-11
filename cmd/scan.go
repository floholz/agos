package cmd

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/grandcat/zeroconf"
)

type DnsDiscoveryResult struct {
	Instance string
	HostName string
	Port     int
	IPv4     string
}

// DiscoverAdbPort scans for open ADB ports on the given IP
// nmap 192.168.0.110 -p 32000-44000 -oG - | grep -P '\d+(?=\/open\/tcp)' --only-matching
func DiscoverAdbPort(ip string, portRange PortRange) (int, error) {
	fmt.Printf("Discovering devices at %s within port range %d-%d...\n", ip, portRange.MinPort, portRange.MaxPort)
	ports := []int{5555}
	for port := portRange.MinPort; port <= portRange.MaxPort; port++ {
		ports = append(ports, port)
	}

	results := make(chan int, 1)
	done := make(chan struct{})
	defer close(done)

	concurrency := 100 // Number of concurrent goroutines
	sem := make(chan struct{}, concurrency)

	for _, port := range ports {
		go func(port int) {
			sem <- struct{}{}        // Acquire a slot
			defer func() { <-sem }() // Release the slot

			address := fmt.Sprintf("%s:%d", ip, port)
			conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
			if err == nil {
				conn.Close()
				select {
				case results <- port:
				case <-done:
				}
			}
		}(port)
	}

	select {
	case port := <-results:
		return port, nil
	case <-time.After(5 * time.Second): // Total timeout for the scan operation
		return 0, fmt.Errorf("no open ADB port found on %s", ip)
	}
}

func DiscoverZeroconf() (*DnsDiscoveryResult, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resolver: %v", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	fmt.Println("Browsing for _adb-tls-pairing._tcp.local...")
	err = resolver.Browse(ctx, "_adb-tls-pairing._tcp", "local", entries)
	if err != nil {
		return nil, fmt.Errorf("failed to browse: %v", err)
	}

	for entry := range entries {
		fmt.Printf("Found servic: %s\n", entry.Instance)
		fmt.Printf("  Host: %s\n", entry.HostName)
		fmt.Printf("  Port: %d\n", entry.Port)
		fmt.Printf("  IPs: %v\n", entry.AddrIPv4) // Or entry.AddrIPv6
		fmt.Printf("  TXT Records: %v\n", entry.Text)

		return &DnsDiscoveryResult{
			Instance: entry.Instance,
			HostName: entry.HostName,
			Port:     entry.Port,
			IPv4:     entry.AddrIPv4[0].String(),
		}, nil
	}

	<-ctx.Done() // Wait for the timeout or cancellation
	return nil, fmt.Errorf("timed out waiting for service entries")
}
