package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
)

var SlugifyCmd = &cobra.Command{
	Use:   "slugify [path]",
	Short: "Renomme fichiers et dossiers en slugs ASCII",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		RunSlugify(path)
	},
}

func RunSlugify(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("❌ Chemin invalide : %v\n", err)
		return
	}

	fmt.Printf("🐌 Slugification dans %s...\n", absPath)

	// On doit traiter les fichiers de manière récursive, mais du plus profond au plus superficiel
	// pour éviter de renommer un dossier avant ses fichiers.
	var entries []string
	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path != absPath {
			entries = append(entries, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Erreur lors du parcours : %v\n", err)
		return
	}

	// Renommer dans l'ordre inverse des profondeurs
	for i := len(entries) - 1; i >= 0; i-- {
		oldPath := entries[i]
		dir := filepath.Dir(oldPath)
		base := filepath.Base(oldPath)
		ext := filepath.Ext(base)
		nameWithoutExt := strings.TrimSuffix(base, ext)

		newBase := slug.Make(nameWithoutExt) + ext
		if ext == "" && os.IsPathSeparator(oldPath[len(oldPath)-1]) {
			newBase = slug.Make(base)
		} else if ext == "" {
			newBase = slug.Make(base)
		}

		newPath := filepath.Join(dir, newBase)

		if oldPath != newPath {
			fmt.Printf("  %s -> %s\n", base, newBase)
			err := os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Printf("  ❌ Erreur : %v\n", err)
			}
		}
	}

	fmt.Println("✅ Slugification terminée.")
}
