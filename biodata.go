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
			0: i.NewChoicePage("What type of biological data?", []string{"DNA", "RNA", "Protein", "Literature"}, []int{1, 2, 3, 4}),
			1: i.NewChoicePage("What sort of DNA data?", []string{"Genome", "Genes", "Variation"}, []int{5, 6, 7}),
			2: i.NewChoicePage("What sort of RNA data?", []string{"Transcript", "Expression"}, []int{8}),
			3: i.NewChoicePage("What sort of protein data?", []string{"Sequence", "Structure", "Interactions"}, []int{9, 10, 11}),
			4: i.NewChoicePage("What sort of literature?", []string{"Keyword", "Author"}, []int{12, 13}),

			// databases
			5: i.NewChoicePage("Here are some genomic databases. Choose one to see what type of queries you can make.",
				[]string{"GenBank", "RefSeq", "Ensembl", "UCSC Genome Browser"}, []int{14, 15, 16, 17}),

			// access
			14: i.NewEntrezPage("GenBank is an archival NCBI dataset, containing all publicly submitted DNA sequences from individual labs and large-scale sequencing projects."),
		},
		PreviousPages: []int{},
		Page:          0,
		Quitting:      false,
		Keys:          i.Keys,
		ShowHelp:      true,
		Help:          help.New(),
		Height:        40,
		Width:         40,
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
