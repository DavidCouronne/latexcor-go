package main

import (
	"fmt"
	"os"

	"latexcor/cmd"
)

func main() {
	// Initialisation des constantes du package cmd
	cmd.Setup(ContainerImage, DefaultEngine, AuxExtensions, AuxPaths)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
