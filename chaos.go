package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Program represents the original JSON structure
type Program struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Bounty  bool     `json:"bounty"`
	Swag    bool     `json:"swag"`
	Domains []string `json:"domains"`
}

// Output represents the desired JSON structure
type Output struct {
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	PlatformURL string `json:"platform_url"`
	Zip         string `json:"ZIP"`
}

// Function to fetch the JSON from the URL
func fetchJSON(url string) ([]Program, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var data struct {
		Programs []Program `json:"programs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Programs, nil
}

// Function to convert programs to the desired output format
func convertPrograms(programs []Program) []Output {
	var outputs []Output
	for _, program := range programs {
		// Convert name to lowercase and replace spaces with underscores for ZIP
		zipName := strings.ToLower(program.Name)
		zipName = strings.ReplaceAll(zipName, " ", "_")

		for _, domain := range program.Domains {
			outputs = append(outputs, Output{
				Domain:      domain,
				Name:        program.Name,
				PlatformURL: program.URL,
				Zip:         fmt.Sprintf("https://chaos-data.projectdiscovery.io/%s.zip", zipName),
			})
		}
	}
	return outputs
}

func main() {
	url := "https://raw.githubusercontent.com/projectdiscovery/public-bugbounty-programs/main/chaos-bugbounty-list.json"

	// Fetch the programs from the URL
	programs, err := fetchJSON(url)
	if err != nil {
		log.Fatalf("Error fetching JSON: %v", err)
	}

	// Convert programs to the desired output format
	outputs := convertPrograms(programs)

	// Convert output to JSON
	jsonData, err := json.MarshalIndent(outputs, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Write the JSON to a file
	err = ioutil.WriteFile("chaos-targets.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}

	fmt.Println("JSON data has been written to chaos-targets.json")
}
