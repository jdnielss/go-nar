package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Condition struct {
	MetricKey      string `json:"metricKey"`
	ActualValue    string `json:"actualValue"`
	Comparator     string `json:"comparator"`
	ErrorThreshold string `json:"errorThreshold"`
	Status         string `json:"status"`
}

type ProjectStatus struct {
	Status     string      `json:"status"`
	Conditions []Condition `json:"conditions"`
}

type Response struct {
	ProjectStatus ProjectStatus `json:"projectStatus"`
}

func main() {
	// Define command-line arguments
	projectName := flag.String("n", "", "CI Project Name")
	sonarURL := flag.String("h", "", "Sonar URL")
	sonarLogin := flag.String("k", "", "Sonar Login")
	flag.Parse()

	// Validate arguments
	if *projectName == "" || *sonarURL == "" || *sonarLogin == "" {
		fmt.Println("Missing required arguments")
		flag.Usage()
		os.Exit(1)
	}

	// Construct the project URL
	projectURL := fmt.Sprintf("%s/api/qualitygates/project_status?projectKey=%s", *sonarURL, *projectName)

	// Make the HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("GET", projectURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}
	req.SetBasicAuth(*sonarLogin, "")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	// Parse the JSON response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error parsing JSON response:", err)
		os.Exit(1)
	}

	// Check the quality gate status
	status := response.ProjectStatus.Status
	if status == "OK" {
		fmt.Println("✅ Quality Gate passed")
	} else {
		fmt.Println("⛔️ Quality Gate failed")
		for _, condition := range response.ProjectStatus.Conditions {
			if condition.Status == "OK" {
				fmt.Printf("✅ %s: %s %s %s => OK\n", condition.MetricKey, condition.ActualValue, condition.Comparator, condition.ErrorThreshold)
			} else {
				fmt.Printf("⛔️ %s: %s %s %s => ERROR\n", condition.MetricKey, condition.ActualValue, condition.Comparator, condition.ErrorThreshold)
			}
		}
	}
}
