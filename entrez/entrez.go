package entrez

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

func SearchDBForQuery(database, query string) ([]string, string, error) {
	res, err := ESearch(database, query)
	if err != nil {
		return nil, "", err
	}
	return res.IdList.Ids, res.Query, err
}

func ESearch(database, query string) (*ESearchResult, error) {
	baseURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi"

	// Build query URL with parameters
	params := url.Values{}
	params.Add("db", database)
	params.Add("term", query)
	params.Add("retmode", "xml")
	params.Add("sort", "date released")

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
	Locus      string  `xml:"GBSeq_locus"`
	Length     int     `xml:"GBSeq_length"`
	Definition string  `xml:"GBSeq_definition"`
	Organism   string  `xml:"GBSeq_organism"`
	References []GBRef `xml:"GBSeq_references>GBReference"`
	MolType    string  `xml:"GBSeq_moltype"`
}

type GBRef struct {
	Title   string   `xml:"GBReference_title"`
	Authors []string `xml:"GBReference_authors>GBAuthor"`
}

func EFetch(database string, ids []string) ([]GBSeq, error) {
	baseURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi"

	url := fmt.Sprintf("%s?db=%s&id=%s&retmode=xml",
		baseURL, database, strings.Join(ids, ","))

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
