package cmd

import (
	"fmt"
	"net"
	"time"
)

// DiscoverAdbPort scans for open ADB ports on the given IP
// nmap 192.168.0.110 -p 32000-44000 -oG - | grep -P '\d+(?=\/open\/tcp)' --only-matching
func DiscoverAdbPort(ip string, minPort, maxPort int) (int, error) {
	fmt.Printf("Discovering devices at %s within port range %d-%d...\n", ip, minPort, maxPort)
	ports := []int{5555}
	for port := minPort; port <= maxPort; port++ {
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
		return 0, fmt.Errorf("No open ADB port found on %s", ip)
	}
}
