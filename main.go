package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/mackerelio/go-osstat/memory"
	"golang.org/x/sys/unix"
)

func main() {
	fmt.Printf("%s | %sG free | %s | %s\n", GetMemoryUsage(), GetFreeSpace(), GetLocalIP(), GetLocalTime())
}

func GetMemoryUsage() string {
	memory, err := memory.Get()
	if err != nil {
		return "error getting memory"
	}

	used := float64(memory.Used) / 1_000_000_000
	total := float64(memory.Total) / 1_000_000_000

	return fmt.Sprintf("%.1fG/%.1fG", used, total)
}

func GetFreeSpace() string {
	var stat unix.Statfs_t

	wd, err := os.Getwd()
	if err != nil {
		return "error getting free space"
	}

	unix.Statfs(wd, &stat)

	// Available blocks * size per block = available space in bytes
	space := stat.Bavail * uint64(stat.Bsize) / 1_000_000_000

	return fmt.Sprintf("%d", space)

}

func GetLocalTime() string {
	t := time.Now()
	return t.Format("Mon 02-01-2006 15:04")
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "no connection"
}
