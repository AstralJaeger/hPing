package main

import (
	"fmt"
	probing "github.com/prometheus-community/pro-bing"
	"strings"
	"time"
)

var endpoints = map[string]string{
	"nbg1": "nbg1-speed.hetzner.com",
	"fsn1": "fsn1-speed.hetzner.com",
	"hel1": "hel1-speed.hetzner.com",
	"ash":  "ash-speed.hetzner.com",
	"hil":  "hil-speed.hetzner.com",
	"sin":  "sin-speed.hetzner.com",
}

var datacenters = map[string]string{
	"nbg1": "Nuremberg 1",
	"fsn1": "Falkenstein 1",
	"hel1": "Helsinki 1",
	"ash":  "Ashburn",
	"hil":  "Hillsboro",
	"sin":  "Singapore",
}

func main() {

	err := PingCloudFlare()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	dcStats := make(map[string]*probing.Statistics)

	for key, value := range endpoints {
		stats, err := PingAndPrint(key, value)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		dcStats[key] = stats
	}

	bestDC := FindBestDC(dcStats)
	fmt.Printf("Best datacenter: %s\n", datacenters[bestDC])
}

func FindBestDC(dcStats map[string]*probing.Statistics) string {
	bestDC := ""
	bestAvg := time.Duration(0)
	for key, value := range dcStats {
		if bestAvg == 0 || value.AvgRtt < bestAvg {
			bestAvg = value.AvgRtt
			bestDC = key
		}
	}
	return bestDC
}

func PingCloudFlare() error {
	fmt.Printf("===== CloudFlare (for reference)\n")
	stats, err := Ping("1.1.1.1", 10*time.Second)
	if err != nil {
		fmt.Println(err)
		return err
	}
	PrintProbeStatistics(stats)
	fmt.Printf("%s\n\n", strings.Repeat("=", 32))
	return nil
}

func PingAndPrint(id string, address string) (*probing.Statistics, error) {
	fmt.Printf("===== Datacenter: %s\n", datacenters[id])
	stats, err := Ping(address, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	PrintProbeStatistics(stats)
	fmt.Printf("%s\n", strings.Repeat("=", 32))
	return stats, nil
}

func PrintProbeStatistics(stats *probing.Statistics) {
	fmt.Printf(" > Packets transmitted: %d\n", stats.PacketsSent)
	fmt.Printf(" > Packets received: %d\n", stats.PacketsRecv)
	fmt.Printf(" > Packet loss: %.2f%%\n", stats.PacketLoss)
	fmt.Printf(" > RTT min: %v\n", stats.MinRtt)
	fmt.Printf(" > RTT max: %v\n", stats.MaxRtt)
	fmt.Printf(" > RTT avg: %v\n", stats.AvgRtt)
	fmt.Printf(" > RTT mdev: %v\n", stats.StdDevRtt)
	fmt.Printf("\n")
}

func Ping(address string, timeout time.Duration) (*probing.Statistics, error) {
	pinger, err := probing.NewPinger(address)
	if err != nil {
		return nil, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = 10
	pinger.Timeout = timeout
	err = pinger.Run()
	if err != nil {
		return nil, err
	}
	return pinger.Statistics(), nil
}
