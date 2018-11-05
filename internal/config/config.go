package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)


var DataDir string

func init() {
	viper.SetDefault("template_dir", "chore_templates")
	viper.SetDefault("data_dir", "marvin_data")
}

func Setup() {
	DataDir = resolvePath(viper.GetString("data_dir"))
	ensureDir(DataDir)
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


