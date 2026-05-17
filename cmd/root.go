package cmd

import (
	"github.com/spf13/cobra"
)

var (
	CfgContainerImage string
	CfgDefaultEngine  string
	CfgAuxExtensions  []string
	CfgAuxPaths       []string
)

// Setup initialise les constantes globales pour le package cmd
func Setup(image, engine string, auxExt []string, auxPaths []string) {
	CfgContainerImage = image
	CfgDefaultEngine = engine
	CfgAuxExtensions = auxExt
	CfgAuxPaths = auxPaths
}

var RootCmd = &cobra.Command{
	Use:   "latexcor",
	Short: "Outil CLI pour compiler des fichiers LaTeX via Podman",
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(CheckCmd)
	RootCmd.AddCommand(CompileCmd)
	RootCmd.AddCommand(WatchCmd)
	RootCmd.AddCommand(CleanCmd)
	RootCmd.AddCommand(SlugifyCmd)
	RootCmd.AddCommand(ConvertCmd)
}
