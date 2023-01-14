package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Get the file name and API key from command line arguments
	fileName := os.Args[1]
	apiKey := os.Args[2]

	// Create the results directory if it doesn't exist
	resultsDir := "results"
	if _, err := os.Stat(resultsDir); os.IsNotExist(err) {
		os.Mkdir(resultsDir, os.ModePerm)
	}

	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is an IP address or FQDN
		if net.ParseIP(line) != nil {
			result := sendToAPI(line, apiKey)
			saveResult(line, result, resultsDir)
		} else if isFQDN(line) {
			result := sendToAPI(line, apiKey)
			saveResult(line, result, resultsDir)
		} else {
			fmt.Println("Invalid input:", line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func sendToAPI(data string, apiKey string) map[string]interface{} {
	url := fmt.Sprintf("https://api.abuseipdb.com/api/v2/check?ipAddress=%s&maxAgeInDays=90", data)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Key", apiKey)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	return result
}

func saveResult(data string, result map[string]interface{}, resultsDir string) {
	// Create the JSON file
	data = strings.ReplaceAll(data, "/", "-") // normalize the file name for subnets defined using CIDR.
	fileName := filepath.Join(resultsDir, data+".json")
	jsonFile, _ := os.Create(fileName)
	defer jsonFile.Close()
	// convert the result map to json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("Result saved to:", fileName)
}

func isFQDN(data string) bool {
	// Check if the data contains a '.'
	if !strings.Contains(data, ".") {
		return false
	}
	// Check if the data contains any spaces
	if strings.Contains(data, " ") {
		return false
	}
	// If the data passes both checks, it is likely a FQDN
	return true
}
