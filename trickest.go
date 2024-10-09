package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Target represents the original JSON structure
type Target struct {
	Name           string   `json:"name"`
	Domains        []string `json:"domains"`
	URL            string   `json:"url"`
	WildcardFilters []string `json:"wildcard_filters"`
}

// Output represents the desired JSON structure
type Output struct {
	Domain       string `json:"domain"`
	Name         string `json:"name"`
	PlatformURL  string `json:"platform_url"`
	GitHubURL    string `json:"github_url"`
	Hostnames    string `json:"hostnames"`
	DNSReport    string `json:"dns_report"`
	ServerReport string `json:"server_report"`
	Servers      string `json:"servers"`
}

// Function to fetch the JSON from the URL
func fetchJSON(url string) ([]Target, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var data struct {
		Targets []Target `json:"targets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Targets, nil
}

// Function to convert targets to the desired output format
func convertTargets(targets []Target) []Output {
	var outputs []Output
	for _, target := range targets {
		for _, domain := range target.Domains {
			nameWithEncodedSpaces := strings.ReplaceAll(target.Name, " ", "%20")
			outputs = append(outputs, Output{
				Domain:       domain,
				Name:         target.Name,
				PlatformURL:  target.URL,
				GitHubURL:    fmt.Sprintf("https://github.com/trickest/inventory/tree/main/%s", nameWithEncodedSpaces),
				Hostnames:    fmt.Sprintf("https://raw.githubusercontent.com/trickest/inventory/main/%s/hostnames.txt", nameWithEncodedSpaces),
				DNSReport:    fmt.Sprintf("https://raw.githubusercontent.com/trickest/inventory/main/%s/dns-report.csv", nameWithEncodedSpaces),
				ServerReport: fmt.Sprintf("https://raw.githubusercontent.com/trickest/inventory/main/%s/server-report.csv", nameWithEncodedSpaces),
				Servers:      fmt.Sprintf("https://raw.githubusercontent.com/trickest/inventory/main/%s/servers.txt", nameWithEncodedSpaces),
			})
		}
	}
	return outputs
}

func main() {
	url := "https://raw.githubusercontent.com/trickest/inventory/main/targets.json"

	// Fetch the targets from the URL
	targets, err := fetchJSON(url)
	if err != nil {
		log.Fatalf("Error fetching JSON: %v", err)
	}

	// Convert targets to the desired output format
	outputs := convertTargets(targets)

	// Convert output to JSON
	jsonData, err := json.MarshalIndent(outputs, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Write the JSON to a file
	err = ioutil.WriteFile("trickest-targets.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}

	fmt.Println("JSON data has been written to trickest-targets.json")
}
