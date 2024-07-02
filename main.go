package main

import (
	"fmt"
	"math"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/mackerelio/go-osstat/memory"
	"golang.org/x/sys/unix"
)

func main() {
	path, _ := BatExists("BAT0")
	perc, _ := GetBatteryLevel(path)

	fmt.Printf("%v%% | %s | %sG free | %s | %s\n", perc, GetMemoryUsage(), GetFreeSpace(), GetLocalIP(), GetLocalTime())
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

func BatExists(name string) (string, bool) {
	filePath := path.Join("/sys", "class", "power_supply", name)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", false
	}

	return filePath, true
}

func GetBatteryLevel(filePath string) (float64, error) {
	energyNow, err := getFloat(path.Join(filePath, "energy_now"))
	if err != nil {
		return 0, err
	}

	energyFull, err := getFloat(path.Join(filePath, "energy_full"))
	if err != nil {
		return 0, err
	}

	return math.Round((energyNow / energyFull) * 100), nil
}

func getFloat(filepath string) (value float64, err error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return value, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(file)), 64)
}
