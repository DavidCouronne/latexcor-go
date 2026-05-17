# latexcor — Instructions pour agent IA (Claude Code / Gemini CLI)

## Contexte projet

Outil CLI en Go pour compiler des fichiers LaTeX via un conteneur Podman.
Cible : enseignants non-développeurs (profs de maths lycée). Distribution = **binaire unique**, pas de runtime.

## Hypothèses fixes — ne jamais contourner

- **Podman est installé** sur la machine. Ne pas gérer le cas "pas de runtime conteneur".
- **L'image `infocornouaille/tools:perso` est présente** (ou à puller une seule fois via `latexcor init`).
- **Pas de Docker** : uniquement Podman.
- **Priorité Linux/macOS**. Windows = support secondaire, géré via `runtime.GOOS` uniquement là où c'est nécessaire (chemins de volume).
- **Pas de vérification défensive à chaque commande** : on suppose l'environnement prêt après `latexcor init`.

---

## Fonctionnalités attendues

| Commande | Description |
|---|---|
| `latexcor init` | Vérifie/démarre la VM Podman (macOS/Windows), pull l'image si absente |
| `latexcor check` | Diagnostic : Podman OK ? VM active ? Image présente ? XeLaTeX fonctionnel ? |
| `latexcor compile [path]` | Compile tous les `.tex` principaux du dossier (2 passes) |
| `latexcor watch [path]` | Surveille le dossier, recompile à chaque sauvegarde `.tex` |
| `latexcor clean [path]` | Supprime les fichiers auxiliaires (`.aux`, `.log`, `.toc`, `.out`, `.synctex.gz`…) |
| `latexcor slugify [path]` | Renomme fichiers et dossiers en slugs ASCII |
| `latexcor convert-utf8 [path]` | Convertit les `.tex` en UTF-8 |

`[path]` est toujours optionnel et défaut au répertoire courant (`os.Getwd()`).

---

## Architecture — structure des fichiers

```
latexcor/
├── main.go            # Entrée : rootCmd cobra + sous-commandes
├── config.go          # Constantes globales uniquement
├── cmd/
│   ├── init.go        # RunInit()
│   ├── check.go       # RunCheck()
│   ├── compile.go     # RunCompile()
│   ├── watch.go       # RunWatch()
│   ├── clean.go       # RunClean()
│   ├── slugify.go     # RunSlugify()
│   └── convert.go     # RunConvertUTF8()
└── internal/
    ├── podman.go      # IsPodmanInstalled, IsPodmanVMRunning, IsImagePresent, BuildCommand
    └── latex.go       # GetMainTexFiles, RunTwoPasses, CleanAux, ExtractErrors
```

**Règle** : `cmd/` contient uniquement la logique CLI (flags, affichage). La logique métier va dans `internal/`.

---

## Conventions Go à respecter

- **Gestion d'erreurs** : toujours retourner `error`, jamais ignorer silencieusement. Utiliser `fmt.Errorf("contexte : %w", err)`.
- **Pas de `os.Chdir()`** : utiliser `cmd.Dir` sur `exec.Cmd` pour définir le répertoire de travail.
- **Pas de parsing de texte brut** pour les sorties de commandes Podman : préférer `--format` avec templates Go (`{{.Running}}`).
- **Affichage** : utiliser `lipgloss` pour les styles (couleurs, bold). Utiliser `schollz/progressbar` pour les barres de progression.
- **Toutes les constantes** (image, engine par défaut, extensions auxiliaires) dans `config.go`, jamais en dur dans le code.

---

## Constantes (config.go)

```go
const (
    ContainerImage    = "infocornouaille/tools:perso"
    DefaultEngine     = "xelatex"
)

var AuxExtensions = []string{".aux", ".log", ".toc", ".out", ".synctex.gz", ".fls", ".fdb_latexmk"}
```

---

## Commande Podman (compilation)

Format attendu de la commande, construit dans `internal/podman.go` :

```
podman run -i --rm -v <chemin_absolu_résolu>:/data infocornouaille/tools:perso xelatex -interaction=nonstopmode -shell-escape <fichier.tex>
```

- `<chemin_absolu_résolu>` = `filepath.Clean(filepath.Abs(...))` — jamais `$(pwd)`, jamais de chemin relatif.
- Sur Linux : ajouter `--security-opt=label=disable` au lieu de `:Z` sur le volume.
- Sur Windows : convertir le chemin (`C:\foo` → `/c/foo`).
- Compilation en **2 passes** systématiquement (nécessaire pour TOC, références croisées).

---

## Dépendances Go (go.mod)

```
github.com/spf13/cobra          # CLI
github.com/charmbracelet/lipgloss  # styles terminal
github.com/schollz/progressbar/v3  # barre de progression
github.com/fsnotify/fsnotify    # watch fichiers
gopkg.in/yaml.v3                # YAML si besoin
github.com/gosimple/slug        # slugify
golang.org/x/text               # détection/conversion encodage UTF-8
```

TOML : stdlib Go 1.21+ (`encoding/toml`), pas de dépendance externe.

---

## Ce qu'il ne faut pas faire

- Ne pas réintroduire de vérification Docker (`exec.LookPath("docker")`).
- Ne pas faire `os.Chdir()` — utiliser `cmd.Dir`.
- Ne pas parser la sortie texte de `podman machine ls` sans `--format`.
- Ne pas dupliquer la logique de construction de commande Podman hors de `internal/podman.go`.
- Ne pas gérer `podman machine` sur Linux (inutile, Podman est natif).