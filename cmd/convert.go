package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/charmap"
)

var ConvertCmd = &cobra.Command{
	Use:   "convert-utf8 [path]",
	Short: "Convertit les .tex en UTF-8",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		RunConvertUTF8(path)
	},
}

func RunConvertUTF8(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("❌ Chemin invalide : %v\n", err)
		return
	}

	fmt.Printf("UTF-8 🔄 Conversion dans %s...\n", absPath)

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".tex" {
			// On tente de détecter si c'est du Latin1 (ISO-8859-1) et on convertit
			// Pour simplifier, on convertit tout ce qui n'est pas déjà UTF-8 valide.
			if !isUTF8(path) {
				fmt.Printf("  Conversion de %s...\n", filepath.Base(path))
				return convertToUTF8(path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Erreur lors de la conversion : %v\n", err)
	} else {
		fmt.Println("✅ Conversion terminée.")
	}
}

func isUTF8(path string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return isUtf8(content)
}

func isUtf8(buf []byte) bool {
	for i := 0; i < len(buf); {
		if buf[i] < 128 {
			i++
		} else if buf[i] >= 192 && buf[i] <= 223 && i+1 < len(buf) && buf[i+1] >= 128 && buf[i+1] <= 191 {
			i += 2
		} else if buf[i] >= 224 && buf[i] <= 239 && i+2 < len(buf) && buf[i+1] >= 128 && buf[i+1] <= 191 && buf[i+2] >= 128 && buf[i+2] <= 191 {
			i += 3
		} else if buf[i] >= 240 && buf[i] <= 247 && i+3 < len(buf) && buf[i+1] >= 128 && buf[i+1] <= 191 && buf[i+2] >= 128 && buf[i+2] <= 191 && buf[i+3] >= 128 && buf[i+3] <= 191 {
			i += 4
		} else {
			return false
		}
	}
	return true
}

func convertToUTF8(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// On assume ISO-8859-1 (Latin1) car c'est le plus commun pour les profs en France
	reader := charmap.Windows1252.NewDecoder().Reader(f)
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	return os.WriteFile(path, content, 0644)
}
