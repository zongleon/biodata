package internal

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ESearchResult struct {
	XMLName  xml.Name `xml:"eSearchResult"`
	QueryKey string   `xml:"QueryKey"`
	WebEnv   string   `xml:"WebEnv"`
	Count    int      `xml:"Count"`
	RetMax   int      `xml:"RetMax"`
	RetStart int      `xml:"RetStart"`
	Query    string   `xml:"QueryTranslation"`
	IdList   struct {
		Ids []string `xml:"Id"`
	} `xml:"IdList"`
}

func SearchDBForQuery(database, filter, query string) ([]string, string, error) {
	res, err := ESearch(database, filter, query)
	if err != nil {
		return nil, "", err
	}
	return res.IdList.Ids, res.Query, err
}

func ESearch(database, filter, query string) (*ESearchResult, error) {
	baseURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi"

	// Build query URL with parameters
	params := url.Values{}
	params.Add("db", database)
	params.Add("term", filter+"[filter] "+query)
	params.Add("retmode", "xml")
	params.Add("sort", "relevance")

	fullURL := baseURL + "?" + params.Encode()

	// Make the request
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse XML response
	var result ESearchResult
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %v", err)
	}

	return &result, nil
}

type GBSet struct {
	XMLName   xml.Name `xml:"GBSet"`
	Sequences []GBSeq  `xml:"GBSeq"`
}

type GBSeq struct {
	Locus            string   `xml:"GBSeq_locus"`
	Length           int      `xml:"GBSeq_length"`
	StrandType       string   `xml:"GBSeq_strandedness"`
	MolType          string   `xml:"GBSeq_moltype"`
	Topology         string   `xml:"GBSeq_topology"`
	Division         string   `xml:"GBSeq_division"`
	Definition       string   `xml:"GBSeq_definition"`
	PrimaryAccession string   `xml:"GBSeq_primary-accession"`
	CreationDate     string   `xml:"GBSeq_create-date"`
	UpdateDate       string   `xml:"GBSeq_update-date"`
	Organism         string   `xml:"GBSeq_organism"`
	Taxonomy         string   `xml:"GBSeq_taxonomy"`
	References       []GBRef  `xml:"GBSeq_references>GBReference"`
	Keywords         []string `xml:"GBSeq_keywords>GBKeyword"`
	Sequence         string   `xml:"GBSeq_sequence"`
}

type GBRef struct {
	Title     string   `xml:"GBReference_title"`
	Authors   []string `xml:"GBReference_authors>GBAuthor"`
	Journal   string   `xml:"GBReference_journal"`
	PubMed    string   `xml:"GBReference_pubmed"`
	RefNumber int      `xml:"GBReference_reference"`
}

func EFetch(database string, ids []string, wholeSeq bool) ([]GBSeq, error) {
	var baseURL, params string
	baseURL = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi"

	if wholeSeq {
		params = "%s?db=%s&id=%s&retmode=xml&rettype=gb"
	} else {
		params = "%s?db=%s&id=%s&retmode=xml&rettype=gb&seq_start=1&seq_stop=1"
	}

	url := fmt.Sprintf(params,
		baseURL, database, strings.Join(ids, ","))

	// url := fmt.Sprintf("%s?db=%s&id=%s&retmode=xml&rettype=gp",
	// 	baseURL, database, strings.Join(ids, ","))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result GBSet
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Sequences, nil
}

var LabelPadding = 20

func (seq GBSeq) PrettyPrint() string {
	var sb strings.Builder
	padding := fmt.Sprintf("%%-%ds", LabelPadding) // creates the %-20s format string dynamically

	// Header with basic information
	sb.WriteString("\n=== SEQUENCE RECORD ===\n")
	sb.WriteString(fmt.Sprintf(padding+" %s\n", "Locus:", seq.Locus))
	sb.WriteString(fmt.Sprintf(padding+" %s\n", "Accession:", seq.PrimaryAccession))
	// sb.WriteString(fmt.Sprintf(padding+" %d bp\n", "Length:", seq.Length))

	// Sequence characteristics
	sb.WriteString("\n--- SEQUENCE CHARACTERISTICS ---\n")
	sb.WriteString(fmt.Sprintf(padding+" %s\n", "Molecule Type:", seq.MolType))
	if seq.StrandType != "" {
		sb.WriteString(fmt.Sprintf(padding+" %s\n", "Strand Type:", seq.StrandType))
	}
	if seq.Topology != "" {
		sb.WriteString(fmt.Sprintf(padding+" %s\n", "Topology:", seq.Topology))
	}
	if seq.Division != "" {
		sb.WriteString(fmt.Sprintf(padding+" %s\n", "Division:", seq.Division))
	}

	// Definition and taxonomy
	sb.WriteString("\n--- DESCRIPTION ---\n")
	sb.WriteString(fmt.Sprintf(padding+" %s\n", "Definition:", seq.Definition))
	sb.WriteString(fmt.Sprintf(padding+" %s\n", "Organism:", seq.Organism))
	// if seq.Taxonomy != "" {
	// 	sb.WriteString(fmt.Sprintf(padding+" %s\n", "Taxonomy:", seq.Taxonomy))
	// }

	// Keywords
	if len(seq.Keywords) > 0 {
		sb.WriteString("\n--- KEYWORDS ---\n")
		for _, keyword := range seq.Keywords {
			sb.WriteString(fmt.Sprintf("  • %s\n", keyword))
		}
	}

	// References
	if len(seq.References) > 0 {
		sb.WriteString("\n--- REFERENCES ---\n")
		for _, ref := range seq.References {
			sb.WriteString(fmt.Sprintf("Reference %d:\n", ref.RefNumber))
			if ref.Title != "" {
				sb.WriteString(fmt.Sprintf("  "+padding+" %s\n", "Title:", ref.Title))
			}
			if len(ref.Authors) > 0 {
				if len(ref.Authors) == 1 {
					sb.WriteString(fmt.Sprintf("  "+padding+" %s\n", "Author:", ref.Authors[0]))
				} else {
					sb.WriteString(fmt.Sprintf("  "+padding+" %s et al.\n", "Authors:", ref.Authors[0]))
				}
			}
			if ref.Journal != "" {
				sb.WriteString(fmt.Sprintf("  "+padding+" %s\n", "Journal:", ref.Journal))
			}
			if ref.PubMed != "" {
				sb.WriteString(fmt.Sprintf("  "+padding+" %s\n", "PubMed:", ref.PubMed))
			}
			sb.WriteString("\n")
		}
	}

	// Dates
	sb.WriteString("--- RECORD INFO ---\n")
	sb.WriteString(fmt.Sprintf(padding+" %s\n", "Created:", seq.CreationDate))
	sb.WriteString(fmt.Sprintf(padding+" %s\n", "Last Updated:", seq.UpdateDate))

	return strings.TrimSpace(sb.String())
}
