package main

import (
	"fmt"
	"os"

	i "github.com/zongleon/biodata/internal"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

func newModel() i.Model {
	return i.Model{
		Pages: map[int]i.Page{
			// type, sub-type
			0: i.NewChoicePage("Data", "What type of biological data?", []string{"DNA", "RNA", "Protein", "Literature"}, []int{1, 2, 3, 4}),
			1: i.NewChoicePage("DNA", "What sort of DNA data?", []string{"Genome", "Genes", "Variation"}, []int{5, 6, 7}),
			2: i.NewChoicePage("RNA", "What sort of RNA data?", []string{"Transcript", "Expression"}, []int{8}),
			3: i.NewChoicePage("Protein", "What sort of protein data?", []string{"Sequence", "Structure", "Interactions"}, []int{9, 10, 11}),
			4: i.NewChoicePage("Literature", "What sort of literature?", []string{"Keyword", "Author"}, []int{12, 13}),

			// databases
			5: i.NewChoicePage("Genomic", "Here are some genomic databases. Choose one to see what type of queries you can make.",
				[]string{"GenBank", "RefSeq", "Ensembl", "UCSC Genome Browser"}, []int{14, 15, 16, 17}),

			// access
			14: i.NewEntrezPage("GenBank", "genbank", "GenBank is an archival NCBI dataset, containing all publicly submitted DNA sequences from individual labs and large-scale sequencing projects."),
			15: i.NewEntrezPage("RefSeq", "refseq", "RefSeq is a manually curated NCBI datasetm, aiming to provide separate and linked records for the genomic DNA, the gene transcripts, and the proteins arising from those transcripts."),
		},
		PreviousPages: []int{},
		PreviousNames: []string{},
		Page:          0,
		Quitting:      false,
		Keys:          i.Keys,
		ShowHelp:      true,
		Help:          help.New(),
	}
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}

}
