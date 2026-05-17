package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"latexcor/internal"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var WatchCmd = &cobra.Command{
	Use:   "watch [path]",
	Short: "Surveille le dossier, recompile à chaque sauvegarde .tex",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		RunWatch(path)
	},
}

func RunWatch(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		var lastCompile time.Time
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) && filepath.Ext(event.Name) == ".tex" {
					// Anti-bounce : ne pas compiler trop souvent
					if time.Since(lastCompile) < 2*time.Second {
						continue
					}
					
					fmt.Printf("🔄 Modification détectée : %s\n", filepath.Base(event.Name))
					
					// Vérifier si c'est un fichier principal
					files, _ := internal.GetMainTexFiles(absPath)
					isMain := false
					for _, f := range files {
						if f == event.Name {
							isMain = true
							break
						}
					}
					
					if isMain {
						buildFunc := func(f string) *exec.Cmd {
							return internal.BuildCommand(CfgContainerImage, CfgDefaultEngine, f)
						}
						internal.RunTwoPasses(buildFunc, event.Name)
						
						// Nettoyage systématique après compilation
						fmt.Println("🧹 Nettoyage automatique...")
						internal.CleanAux(absPath, CfgAuxExtensions, CfgAuxPaths)
						
						lastCompile = time.Now()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(absPath)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("👀 Surveillance active dans %s...\n", absPath)
	<-done
}
