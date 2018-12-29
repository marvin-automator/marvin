// +build mage

package main

import (
	"fmt"
	"github.com/gobuffalo/envy"
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

const V8_VERSION = "6.3.292.48.1"

func Setup() {
	println("Installing deps...")
	runCommand(true, "", "go", "get", "-u", "-v", "./...")
	runCommand(true, "", "go", "-get", "-u", "-v", "github.com/gobuffalo/packr/packr")

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
	filename := fmt.Sprintf("libv8-%v-%v.gem", V8_VERSION, os_arch)

	wd, err := os.Getwd()
	errorAndExit(err)

	os.Mkdir("v8binary", 0700)
	runCommand(false, wd+"/v8binary", "curl", "https://rubygems.org/downloads/"+filename, "-o", filename)
	runCommand(false, wd+"/v8binary", "tar", "-xvf", filename)
	runCommand(false, wd+"/v8binary", "tar", "-xzvf", "data.tar.gz")
	runCommand(false, "", "ln", "-f", "-s", wd+"/v8binary/vendor/v8/out/x64.release", envy.GoPath()+"/src/github.com/augustoroman/v8/libv8")
	runCommand(false, "", "ln", "-f", "-s", wd+"/v8binary/vendor/v8/include", envy.GoPath()+"/src/github.com/augustoroman/v8/include")

	println("Installing frontend deps...")
	runCommand(false, wd+"", "npm", "install", "-g", "yarn")
	runCommand(false, wd+"/frontend", "yarn")
	println("Building frontend...")
	runCommand(false, wd+"/frontend", "yarn", "run", "build")

	println("Building and installing marvin...")
	runCommand(false, "", "packr", "install")
}
