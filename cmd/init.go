package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"latexcor/internal"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Vérifie/démarre la VM Podman et pull l'image si absente",
	Run: func(cmd *cobra.Command, args []string) {
		RunInit()
	},
}

func RunInit() {
	fmt.Println("🚀 --- INITIALISATION DE L'ENVIRONNEMENT ---")

	if !internal.IsPodmanInstalled() {
		fmt.Println("❌ Podman n'est pas installé. Veuillez l'installer avant de continuer.")
		return
	}

	if !internal.IsPodmanVMRunning() {
		fmt.Println("⏳ Tentative de lancement de la machine virtuelle Podman...")

		cmdStart := exec.Command("podman", "machine", "start")
		err := cmdStart.Run()

		if err != nil {
			fmt.Println("ℹ️  Aucune machine détectée. Tentative d'initialisation (podman machine init)...")

			cmdInit := exec.Command("podman", "machine", "init")
			cmdInit.Stdout = os.Stdout
			cmdInit.Stderr = os.Stderr

			if errInit := cmdInit.Run(); errInit != nil {
				fmt.Println("❌ Erreur critique lors de l'initialisation de la machine.")
				return
			}

			fmt.Println("⏳ Initialisation terminée. Lancement en cours...")
			exec.Command("podman", "machine", "start").Run()
		}

		fmt.Println("✅ VM opérationnelle.")
	}

	if !internal.IsImagePresent(CfgContainerImage) {
		pullImage()
	}

	fmt.Println("✨ Environnement prêt !")
}

func pullImage() {
	fmt.Printf("⏳ Téléchargement de l'image %s...\n", CfgContainerImage)
	cmd := exec.Command("podman", "pull", CfgContainerImage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Erreur lors du pull de l'image : %v\n", err)
		return
	}
	fmt.Println("✅ Image récupérée.")
}
