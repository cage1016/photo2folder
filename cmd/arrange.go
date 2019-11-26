package cmd

import (
	"github.com/cage1016/photo2folder/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	yyyyMM   = "yyyyMM"
	yyyyMMdd = "yyyyMMdd"
)

var folderPatten map[string]string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "arrange",
	Aliases: []string{"a"},
	Short:   "Use to arrange photo by date pattern yyyyMM(default), yyyyMMdd",
	Example: `arrange by folder pattern by yyyyMM
-a /Users/cage/Downloads/aa -p yyyyMM

arrange by folder pattern by yyyyMMdd
-a /Users/cage/Downloads/aa -p yyyyMMdd
`,
	Run: func(cmd *cobra.Command, args []string) {
		folder, _ := cmd.Flags().GetString("folder")
		if folder == "" {
			logrus.Error("You must provide the folder name")
			return
		}
		pattern, _ := cmd.Flags().GetString("pattern")
		if pattern != yyyyMM && pattern != yyyyMMdd {
			logrus.Error("folder patten must be one of ", yyyyMM, "(default), ", yyyyMMdd)
			return
		}
		utils.Arrange(folder, folderPatten[pattern])
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("folder", "f", "", "Set your execute folder")
	addCmd.Flags().StringP("pattern", "p", yyyyMM, "Set your folder pattern")

	folderPatten = map[string]string{}
	folderPatten[yyyyMM] = "200601"
	folderPatten[yyyyMMdd] = "20060102"
}
