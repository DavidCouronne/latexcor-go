package internal

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"path/filepath"
)

// IsPodmanInstalled vérifie la présence du binaire podman
func IsPodmanInstalled() bool {
	_, err := exec.LookPath("podman")
	return err == nil
}

// IsPodmanVMRunning vérifie si au moins une VM podman est en cours d'exécution (macOS/Windows)
func IsPodmanVMRunning() bool {
	if runtime.GOOS == "linux" {
		return true // Pas de VM sur Linux
	}

	// Utilisation de --format pour éviter le parsing de texte brut
	out, err := exec.Command("podman", "machine", "ls", "--format", "{{.Running}}").Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(out), "true")
}

// IsImagePresent vérifie si l'image spécifiée est présente localement
func IsImagePresent(image string) bool {
	err := exec.Command("podman", "image", "inspect", image).Run()
	return err == nil
}

// BuildCommand construit la commande podman run pour compiler un fichier LaTeX
func BuildCommand(image, engine, texFile string) *exec.Cmd {
	absPath, err := filepath.Abs(filepath.Dir(texFile))
	if err != nil {
		absPath = "."
	}
	
	fileName := filepath.Base(texFile)
	
	volumeMapping := fmt.Sprintf("%s:/data", absPath)
	
	args := []string{
		"run", "-i", "--rm",
		"-v", volumeMapping,
	}

	if runtime.GOOS == "linux" {
		args = append(args, "--security-opt=label=disable")
	}

	args = append(args, image, engine, "-interaction=nonstopmode", "-shell-escape", fileName)
	
	cmd := exec.Command("podman", args...)
	cmd.Dir = absPath
	return cmd
}
