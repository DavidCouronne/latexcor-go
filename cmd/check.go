package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"latexcor/internal"

	"github.com/spf13/cobra"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Diagnostic : Podman OK ? VM active ? Image présente ? XeLaTeX fonctionnel ?",
	Run: func(cmd *cobra.Command, args []string) {
		RunCheck()
	},
}

func RunCheck() {
	fmt.Println("🔍 --- DIAGNOSTIC SYSTÈME ---")

	if internal.IsPodmanInstalled() {
		fmt.Println("✅ Podman est installé.")
	} else {
		fmt.Println("❌ Podman n'est pas détecté.")
		return
	}

	if internal.IsPodmanVMRunning() {
		fmt.Println("✅ Machine virtuelle Podman : Lancée.")
	} else {
		fmt.Println("⚠️  Machine virtuelle Podman : Éteinte.")
	}

	if internal.IsImagePresent(CfgContainerImage) {
		fmt.Printf("✅ Image %s : Présente.\n", CfgContainerImage)
		
		// Test xelatex interne
		out, err := exec.Command("podman", "run", "--rm", CfgContainerImage, "xelatex", "--version").Output()
		if err == nil {
			fmt.Printf("ℹ️  XeLaTeX conteneur : %s\n", strings.Split(string(out), "\n")[0])
		} else {
			fmt.Println("❌ XeLaTeX inaccessible dans le conteneur.")
		}
	} else {
		fmt.Printf("⚠️  Image %s : Manquante.\n", CfgContainerImage)
	}
}
