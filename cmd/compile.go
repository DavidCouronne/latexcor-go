package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"latexcor/internal"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var CompileCmd = &cobra.Command{
	Use:   "compile [path]",
	Short: "Compile tous les .tex principaux du dossier (2 passes)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		RunCompile(path)
	},
}

func RunCompile(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("❌ Chemin invalide : %v\n", err)
		return
	}

	files, err := internal.GetMainTexFiles(absPath)
	if err != nil {
		fmt.Printf("❌ Erreur lors de la recherche des fichiers : %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("ℹ️  Aucun fichier LaTeX principal trouvé (\\begin{document} manquant).")
		return
	}

	fmt.Printf("📂 %d fichier(s) à compiler dans %s\n", len(files), absPath)

	bar := progressbar.Default(int64(len(files)), "Compilation")

	for _, file := range files {
		buildFunc := func(f string) *exec.Cmd {
			return internal.BuildCommand(CfgContainerImage, CfgDefaultEngine, f)
		}
		
		err := internal.RunTwoPasses(buildFunc, file)
		if err != nil {
			fmt.Printf("\n❌ Erreur pour %s : %v\n", filepath.Base(file), err)
			logFile := strings.TrimSuffix(file, ".tex") + ".log"
			errors := internal.ExtractErrors(logFile)
			if len(errors) > 0 {
				fmt.Println("  Détails des erreurs :")
				for _, e := range errors {
					fmt.Printf("    - %s\n", e)
				}
			}
		}
		bar.Add(1)
	}

	// Nettoyage après compilation
	fmt.Println("\n🧹 Nettoyage automatique...")
	internal.CleanAux(absPath, CfgAuxExtensions, CfgAuxPaths)

	fmt.Println("✨ Compilation terminée.")
}
