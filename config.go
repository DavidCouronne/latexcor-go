package main

const (
	ContainerImage = "infocornouaille/tools:perso"
	DefaultEngine  = "xelatex"
)

var AuxExtensions = []string{
	".bbl", ".blg", ".synctex", ".bar", ".cor", ".lua", ".lub", ".tab",
	".log", ".gz", ".aux", ".out", ".fdb_latexmk", ".fls", ".xdv", ".dvi",
	".bara", ".barb", ".synctex.gz", ".toc",
}

var AuxPaths = []string{"_minted-", "mermaid"}
