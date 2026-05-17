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
			if strings.Contains(string(content), "\\begin{document}") {
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
				if strings.Contains(base, p) {
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

// ExtractErrors extrait les erreurs d'un fichier .log
func ExtractErrors(logPath string) []string {
	file, err := os.Open(logPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var errors []string
	scanner := bufio.NewScanner(file)
	// Regex simple pour les erreurs LaTeX
	re := regexp.MustCompile(`^! (.*)`)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			errors = append(errors, matches[1])
		}
	}
	return errors
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
