# latexcor — Compilation LaTeX simplifiée pour les profs

`latexcor` est un petit outil en ligne de commande (CLI) conçu pour les enseignants (notamment de mathématiques) qui souhaitent compiler leurs documents LaTeX sans avoir à installer et maintenir une distribution complète (comme TeX Live ou MiKTeX) sur leur machine.

## 🚀 Pourquoi utiliser latexcor ?

- **Binaire unique** : Pas besoin de runtime (Python, Node, etc.).
- **Environnement isolé** : Tout se passe dans un conteneur **Podman**. Votre système reste propre.
- **Prêt pour le lycée** : Utilise une image spécialement dédiée (`infocornouaille/tools:perso`).
- **Automatique** : Surveillez vos fichiers et compilez-les à chaque sauvegarde avec nettoyage automatique.

---

## ⚡ Installation rapide (Copier-Coller)

Choisissez la méthode correspondant à votre système pour installer `latexcor` dans `/usr/local/bin`.

### macOS (Apple Silicon - M1/M2/M3)
```bash
sudo curl -L https://github.com/DavidCouronne/latexcor-go/releases/latest/download/latexcor-macos-arm64 -o /usr/local/bin/latexcor && sudo chmod +x /usr/local/bin/latexcor
```

### macOS (Intel)
```bash
sudo curl -L https://github.com/DavidCouronne/latexcor-go/releases/latest/download/latexcor-macos-intel -o /usr/local/bin/latexcor && sudo chmod +x /usr/local/bin/latexcor
```

### Linux (x86_64)
```bash
sudo curl -L https://github.com/DavidCouronne/latexcor-go/releases/latest/download/latexcor-linux -o /usr/local/bin/latexcor && sudo chmod +x /usr/local/bin/latexcor
```

> **Note** : Vous devez avoir **Podman** installé. Sur macOS : `brew install podman`.

---

## 🛠️ Premier lancement

Une fois installé, préparez l'environnement avec :

```bash
latexcor init
```

Cette commande vérifie Podman et télécharge l'image de compilation (environ 2-3 Go, une seule fois).

---

## 📖 Utilisation courante

### 1. Mode automatique (Recommandé)
Surveillez votre dossier et recompilez dès que vous sauvegardez un fichier `.tex`. Les fichiers temporaires sont nettoyés automatiquement après chaque passe.

```bash
latexcor watch
```
*Appuyez sur `Ctrl+C` pour quitter.*

### 2. Compilation ponctuelle
Pour compiler tous les fichiers `.tex` principaux du dossier :

```bash
latexcor compile
```

### 3. Nettoyage manuel
Si besoin, pour forcer la suppression des fichiers auxiliaires (`.aux`, `.log`, dossiers `_minted`, etc.) :

```bash
latexcor clean
```

### 4. Fonctions avancées
- **`latexcor slugify`** : Renomme fichiers/dossiers en ASCII simple (ex: `Cours de Maths.tex` → `cours-de-maths.tex`).
- **`latexcor convert-utf8`** : Convertit vos anciens fichiers `.tex` vers l'encodage moderne UTF-8.

---

## 👨‍🏫 Auteur
Développé pour faciliter la vie des profs de maths.
