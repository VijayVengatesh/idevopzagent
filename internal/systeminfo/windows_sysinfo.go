//go:build windows
// +build windows

package systeminfo

import "os"

type WindowsFetcher struct{}

func (w WindowsFetcher) GetInfo() Info {
	hostname, _ := os.Hostname()
	return Info{Hostname: hostname, Platform: "Windows"}
}

func GetFetcher() Fetcher {
	return WindowsFetcher{}
}
