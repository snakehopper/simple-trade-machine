package function

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// Uses directory "serverless_function_source_code" as defined in the Go
// Functions Framework Buildpack.
// See https://github.com/GoogleCloudPlatform/buildpacks/blob/56eaad4dfe6c7bd0ecc4a175de030d2cfab9ae1c/cmd/go/functions_framework/main.go#L38.
const fnSourceDir = "./serverless_function_source_code"

type Whitelist map[string]bool

func NewWhitelist() Whitelist {
	var wl Whitelist = make(map[string]bool)

	for _, fp := range []string{filepath.Join(fnSourceDir, "ip-whitelist"), "ip-whitelist"} {
		bs, err := ioutil.ReadFile(fp)
		if err != nil {
			continue
		}
		for _, ln := range strings.Fields(string(bs)) {
			wl[ln] = true
		}
	}

	return wl
}

func (wl Whitelist) Allow(req *http.Request) bool {
	if wl[req.RemoteAddr] {
		return true
	}

	if wl[req.Header.Get("X-Forwarded-For")] {
		return true
	}

	fmt.Printf("incoming illegal ip request: %+v\n", req)
	return false
}
