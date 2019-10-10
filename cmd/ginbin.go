package cmd

import (
	"github.com/spf13/cobra"
	"path"
	"sshfortress/util"
)

var ginbinCmd = &cobra.Command{
	Use:   "ginbin",
	Short: "",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {
		util.RunGinStatic(flagSrc, flagDest, flagTags, flagPkg, flagPkgCmt, flagNoMtime, flagNoCompress, flagForce)
	},
}
var (
	flagSrc, flagDest, flagTags, flagPkg, flagPkgCmt string
	flagNoMtime, flagNoCompress, flagForce           bool
)

func init() {
	rootCmd.AddCommand(ginbinCmd)

	ginbinCmd.Flags().StringVarP(&flagSrc, "src", "s", path.Join(".", "dist"), "The path of the source directory.")
	ginbinCmd.Flags().StringVarP(&flagDest, "dest", "d", ".", "The destination path of the generated package.")
	ginbinCmd.Flags().StringVarP(&flagTags, "tags", "t", "", "The golang tags.")
	ginbinCmd.Flags().StringVarP(&flagPkg, "package", "p", "felixbin", "The destination path of the generated package.")
	ginbinCmd.Flags().StringVarP(&flagPkgCmt, "comment", "c", "", "The package comment. An empty value disables this comment.")
	ginbinCmd.Flags().BoolVarP(&flagNoCompress, "zip", "z", false, "Do not use compression to shrink the files.")
	ginbinCmd.Flags().BoolVarP(&flagNoMtime, "mtime", "m", false, "Ignore modification times on files.")
	ginbinCmd.Flags().BoolVarP(&flagForce, "force", "f", true, "Overwrite destination file if it already exists.")
}
