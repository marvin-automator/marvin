// +build mage

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/envy"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func errorAndExit(err error) {
	if err != nil {
		println("Error encountered: ", err.Error())
		os.Exit(1)
	}
}

func runCommand(ignoreError bool, path, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = path

	println(name, strings.Join(args, " "))
	err := cmd.Run()
	if ignoreError {
		return
	}
	errorAndExit(err)
}

func Setup() {
	println("Installing deps...")
	runCommand(true, "", "go", "get",  "-u", "-v", "./...")
	runCommand(true, "", "go", "-get", "-u", "-v", "github.com/gobuffalo/packr/packr")

	println("Determining latest v8 version...")
	resp, err := http.Get("https://rubygems.org/api/v1/versions/libv8/latest.json")
	errorAndExit(err)
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	type version struct {
		Version string
	}
	v := version{}
	dec.Decode(&v)

	var os_arch string
	switch runtime.GOOS {
	case "darwin":
		os_arch = "x86_64-darwin-16"
	case "linux":
		os_arch = "x86_64-linux"
	default:
		println("Unsupported os", runtime.GOOS)
		os.Exit(1)
	}
	filename := fmt.Sprintf("libv8-%v-%v.gem", v.Version, os_arch)

	wd, err := os.Getwd()
	errorAndExit(err)

	os.Mkdir("v8binary", 0700)
	runCommand(false, wd+ "/v8binary", "curl", "https://rubygems.org/downloads/" + filename, "-o", filename)
	runCommand(false, wd + "/v8binary", "tar", "-xvf", filename)
	runCommand(false, wd + "/v8binary", "tar", "-xzvf", "data.tar.gz")
	runCommand(false, "", "ln", "-f", "-s", wd + "/v8binary/vendor/v8/out/x64.release",  envy.GoPath() + "/src/github.com/augustoroman/v8/libv8")
	runCommand(false, "", "ln", "-f", "-s", wd + "/v8binary/vendor/v8/include", envy.GoPath() + "/src/github.com/augustoroman/v8/include")
	runCommand(false, "", "packr", "install")
}
