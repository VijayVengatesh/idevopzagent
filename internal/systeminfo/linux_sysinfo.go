//go:build linux
// +build linux

package systeminfo

import "os"

type LinuxFetcher struct{}

func (l LinuxFetcher) GetInfo() Info {
	hostname, _ := os.Hostname()
	return Info{Hostname: hostname, Platform: "Linux"}
}

func GetFetcher() Fetcher {
	return LinuxFetcher{}
}
