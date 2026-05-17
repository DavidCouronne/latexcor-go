🐹 Aide-mémoire Go (pour un développeur Python/NSI)

🚀 Gestion du Projet & DépendancesGo utilise des modules. 

Le fichier go.mod est l'équivalent de votre pyproject.toml ou requirements.txt.CommandeActiongo mod init <nom>Initialise un nouveau projet (ex: go mod init latexcor).go get <url>Télécharge une dépendance et l'ajoute au projet (ex: go get [github.com/saintfish/chardet](https://github.com/saintfish/chardet)).go mod tidyIndispensable : Nettoie le go.mod (ajoute les imports manquants, supprime les inutiles).go clean -modcacheVide le cache global des modules (pour libérer de l'espace disque).🛠 Développement & CompilationContrairement à Python, Go est un langage compilé.CommandeActiongo run main.goCompile en RAM et exécute. Parfait pour le développement.go build -o <nom>Compile un binaire exécutable optimisé pour votre machine../<nom>Exécute le binaire produit (le ./ est obligatoire sur Mac/Linux).go installCompile et place le binaire dans $HOME/go/bin (accessible partout via le PATH).🎨 Formatage & Qualité du CodeLe style est imposé par le langage pour garantir une lecture universelle.CommandeActiongo fmt ./...Formate récursivement tous les fichiers du projet.go vet ./...Analyse statique pour détecter les erreurs potentielles (variables inutilisées, etc.).Cmd + SDans VS Code, le formatage et la gestion des imports sont automatiques.💡 Points clés de syntaxe (VS Python)Gestion des erreurs : Pas de try/except. On vérifie explicitement chaque retour :Gocontent, err := os.ReadFile("file.tex")
if err != nil {
    // Gérer l'erreur ici
}
Visibilité :Une fonction qui commence par une Majuscule (DetectEncoding) est publique (exportée).Une fonction qui commence par une minuscule (detectEncoding) est privée au package.Types : Go est statiquement typé. var x int = 10 ou plus court x := 10.Dépendances : Pas de dossier node_modules ou venv local. Tout est centralisé dans $HOME/go/pkg/mod et partagé entre vos projets.📂 Structure de fichiers recommandéePour votre projet latexcor, voici une structure simple et propre :Plaintextlatexcor/
├── go.mod          # Définition du module
├── go.sum          # Checksums des dépendances
├── main.go         # Point d'entrée (CLI / Typer equivalent)
├── encoding.go     # Logique de détection (votre module Python converti)
└── compiler.go     # Logique Podman / Compilation
🍎 Astuce Mac (Homebrew)Si vous voulez que vos outils installés via go install fonctionnent partout, ajoutez ceci à votre fichier ~/.zshrc :Bashexport PATH=$PATH:$(go env GOPATH)/bin