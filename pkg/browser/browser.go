// Package browser
//
//	@program:		gin-template
//	@author:		[lliuhuan](https://github.com/lliuhuan)
//	@create:		2024-07-02 21:53
//	@description:	使用浏览器打开指定的 URL
package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

var commands = map[string]string{
	"windows": "start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

func Open(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}
	if len(uri) > 4 && uri[:4] != "http" {
		uri = "http://" + uri
	}
	cmd := exec.Command(run, uri)
	return cmd.Start()
}
