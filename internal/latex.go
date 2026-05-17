package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// GetMainTexFiles retourne la liste des fichiers .tex qui semblent être des fichiers principaux (contiennent \begin{document})
func GetMainTexFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".tex" {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil // Ignorer les fichiers illisibles
			}
			contentStr := string(content)
			if strings.Contains(contentStr, "\\documentclass") &&
				strings.Contains(contentStr, "\\begin{document}") &&
				strings.Contains(contentStr, "\\end{document}") {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

// CleanAux supprime les fichiers et dossiers auxiliaires
func CleanAux(dir string, extensions []string, paths []string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		base := info.Name()

		// Vérifier si c'est un dossier à supprimer
		if info.IsDir() {
			for _, p := range paths {
				// Gestion des wildcards simples comme _minted-*
				match, _ := filepath.Match(p, base)
				if match || strings.Contains(base, p) {
					os.RemoveAll(path)
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Vérifier si c'est un fichier à supprimer par extension
		ext := filepath.Ext(path)
		for _, auxExt := range extensions {
			if ext == auxExt {
				os.Remove(path)
				break
			}
		}

		return nil
	})
}

// ExtractErrors extrait les blocs d'erreurs d'un fichier .log (approche plus robuste par blocs)
func ExtractErrors(logPath string) []string {
	content, err := os.ReadFile(logPath)
	if err != nil {
		return nil
	}

	// Regex pour capturer le bloc d'erreur commençant par ! et finissant par l.ligne
	// Similaire à la logique Python: !(.*?)\n(?:l\.\d+.*?)(?=\n\n|\Z)
	re := regexp.MustCompile(`(?s)! (.*?)\nl\.\d+.*?(?:\n\n|\z)`)
	matches := re.FindAllString(string(content), -1)

	if len(matches) == 0 {
		// Fallback sur une détection ligne par ligne si le bloc n'est pas trouvé
		var errors []string
		scanner := bufio.NewScanner(strings.NewReader(string(content)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "! ") {
				errors = append(errors, line)
			}
		}
		return errors
	}

	var cleanedMatches []string
	for _, m := range matches {
		cleanedMatches = append(cleanedMatches, strings.TrimSpace(m))
	}

	return cleanedMatches
}

// RunTwoPasses exécute deux passes de compilation pour un fichier
func RunTwoPasses(buildCmdFunc func(string) *exec.Cmd, texFile string) error {
	for i := 1; i <= 2; i++ {
		fmt.Printf("⏳ Passe %d/2 pour %s...\n", i, filepath.Base(texFile))
		cmd := buildCmdFunc(texFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("erreur lors de la passe %d : %w", i, err)
		}
	}
	return nil
}
