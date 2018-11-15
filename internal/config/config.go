package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var DevMode bool

var DataDir string
var ServerHost string

func init() {
	viper.SetDefault("data_dir", "marvin_data")
	viper.SetDefault("server_host", ":80")
	viper.SetDefault("dev_mode", false)
}

func Setup() {
	DataDir = resolvePath(viper.GetString("data_dir"))
	ensureDir(DataDir)

	ServerHost = viper.GetString("server_host")

	DevMode = viper.GetBool("dev_mode")
}

func resolvePath(p string) string {
	if !filepath.IsAbs(p) {
		d := filepath.Dir(viper.ConfigFileUsed())
		p = filepath.Join(d, p)
	}

	return filepath.Clean(p)
}

func ensureDir(d string) {
	info, err := os.Stat(d)
	if os.IsNotExist(err) {
		fmt.Println("Creating directory:", d)
		err = os.MkdirAll(d, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !info.IsDir() {
		fmt.Printf("%v is not a directory!", d)
		os.Exit(1)
	}
}
