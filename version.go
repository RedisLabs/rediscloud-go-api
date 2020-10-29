package rediscloud_api

import (
	"fmt"
	"runtime"
	"strings"
)

const Version = "0.1.0"

var userAgent = buildUserAgent("rediscloud-go-api", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)

func buildUserAgent(name string, version string, info ...string) string {
	product := fmt.Sprintf("%s/%s", name, version)
	systemInfo := strings.Join(info, "; ")
	return fmt.Sprintf("%s (%s)", product, systemInfo)
}
