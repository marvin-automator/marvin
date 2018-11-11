package cmd

import (
	"bufio"
	"fmt"
	"github.com/marvin-automator/marvin/internal/config"
	"github.com/marvin-automator/marvin/internal/web"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "marvin",
	Short: "Marvin is a tool for automating all sorts of things.",
	Long: `Here I am, brain the size of a planet. and they ask me to check the weather, update their calendar, send emails, ... Call that job satisfaction? I don't.

Marvin allows you to automate all sorts of things.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Starting app...")
		web.RunApp()
		fmt.Println("Goodbye!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./marvin.yaml, and if that doesn't exist, $HOME/.marvin.yaml)")
	rootCmd.AddCommand(importCmt)
	rootCmd.AddCommand(setPassword)
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("Couldn't determine the home directory:", err, "\nSkipping looking for config files there.")
		} else {
			viper.AddConfigPath(home)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("marvin")
	}

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			if cfgFile != "" {
				fmt.Println("No config file found at the given location:", cfgFile)
				os.Exit(1)
			}
			err2 := createConfigFile()
			if err2 != nil {
				fmt.Println(err2)
			}

		default:
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Loaded config file from: %s\n", viper.ConfigFileUsed())
	}

	config.Setup()
}

func createConfigFile() error {
	fmt.Print("No configuration found. We'll create a default one for you. Where should we create it? (leave blank to use ./marvin.yaml): ")

	r := bufio.NewReader(os.Stdin)
	l, _, err := r.ReadLine()
	if err != nil {
		return err
	}

	name := string(l)
	if name == "" {
		name = "./marvin.yaml"
	}

	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("file exists at %v", name)
	}

	return viper.WriteConfigAs(name)
}
