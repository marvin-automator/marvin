package cmd

import (
	"fmt"
	"github.com/marvin-automator/marvin/internal"
	"github.com/marvin-automator/marvin/internal/chores"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

var importCmt = &cobra.Command{
	Use: "import <filepath> [-n=name]",
	Short: "Import a javascript file as a template",

	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		info, err := os.Stat(filename)
		internal.ErrorAndExit(err)

		f, err := os.Open(filename)
		internal.ErrorAndExit(err)

		bytes, err := ioutil.ReadAll(f)
		internal.ErrorAndExit(err)
		data := string(bytes)

		name := cmd.Flags().Lookup("name").Value.String()
		if name == "" {
			name = info.Name()
			dot := strings.LastIndex(name, ".")
			if dot > 0 { // We also want to ignore the case of a hidden file (starting with a dot)
				name = name[:dot]
			}
		}

		ct, err := chores.NewChoreTemplate(name, data)
		internal.ErrorAndExit(err)

		internal.ErrorAndExit(ct.Save())
		fmt.Println("Successfully saved template: ", name)
	},
}

func init(){
	importCmt.Flags().StringP("name", "n", "", "The name of the template. Will be taken from the filename if you leave this blank.")
}
