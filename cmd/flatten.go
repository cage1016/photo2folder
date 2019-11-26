package cmd

import (
	"github.com/cage1016/photo2folder/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var flattenCmd = &cobra.Command{
	Use:     "flatten",
	Aliases: []string{"f"},
	Short:   "Flatten photos from sub directories",
	Example: `Flatten sub directories photos by given directory
-f /Users/cage/Downloads/aa
`,
	Run: func(cmd *cobra.Command, args []string) {
		folder, _ := cmd.Flags().GetString("folder")
		if folder == "" {
			logrus.Error("You must provide the folder name")
			return
		}
		err := utils.Flatten(folder)
		if err != nil {
			logrus.Error(err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(flattenCmd)
	flattenCmd.Flags().StringP("folder", "f", "", "Set your execute folder")
}
