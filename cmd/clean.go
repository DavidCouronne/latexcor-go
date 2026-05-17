package cmd

import (
	"fmt"
	"path/filepath"

	"latexcor/internal"

	"github.com/spf13/cobra"
)

var CleanCmd = &cobra.Command{
	Use:   "clean [path]",
	Short: "Supprime les fichiers auxiliaires et dossiers temporaires",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		RunClean(path)
	},
}

func RunClean(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("❌ Chemin invalide : %v\n", err)
		return
	}

	fmt.Printf("🧹 Nettoyage dans %s...\n", absPath)
	err = internal.CleanAux(absPath, CfgAuxExtensions, CfgAuxPaths)
	if err != nil {
		fmt.Printf("❌ Erreur lors du nettoyage : %v\n", err)
	} else {
		fmt.Println("✅ Nettoyage terminé.")
	}
}
