// internal/systeminfo/common.go
package systeminfo

type Info struct {
	Hostname string
	Platform string
}

type Fetcher interface {
	GetInfo() Info
}
