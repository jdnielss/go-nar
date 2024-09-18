package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

type GitLabNote struct {
	Body string `json:"body"`
}

func main() {
	// Define command-line arguments
	projectName := flag.String("n", "", "CI Project Name")
	sonarURL := flag.String("h", "", "Sonar URL")
	sonarLogin := flag.String("k", "", "Sonar Login")
	gitlabURL := flag.String("g", "", "GitLab URL")
	gitlabToken := flag.String("t", "", "GitLab Token")
	projectID := flag.String("p", "", "GitLab Project ID")
	mergeRequestIID := flag.String("m", "", "GitLab Merge Request IID")

	flag.Parse()

	// Validate arguments
	if *projectName == "" || *sonarURL == "" || *sonarLogin == "" || *gitlabURL == "" || *gitlabToken == "" || *projectID == "" || *mergeRequestIID == "" {
		fmt.Println("Missing required arguments")
		flag.Usage()
		os.Exit(1)
	}

	// Construct the project URL
	projectURL := fmt.Sprintf("%s/api/qualitygates/project_status?projectKey=%s", *sonarURL, *projectName)

	// Make the HTTP request to SonarQube
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

	// Check the quality gate status and prepare the message
	status := response.ProjectStatus.Status
	var message strings.Builder

	if status == "OK" {
		message.WriteString("✅ Quality Gate passed\n")
	} else {
		message.WriteString("⛔️ Quality Gate failed\n")
		for _, condition := range response.ProjectStatus.Conditions {
			if condition.Status == "OK" {
				message.WriteString(fmt.Sprintf("✅ %s: %s %s %s => OK\n", condition.MetricKey, condition.ActualValue, condition.Comparator, condition.ErrorThreshold))
			} else {
				message.WriteString(fmt.Sprintf("⛔️ %s: %s %s %s => ERROR\n", condition.MetricKey, condition.ActualValue, condition.Comparator, condition.ErrorThreshold))
			}
		}
	}

	fmt.Println(message.String()) // Print the message to console

	// Create a note on the GitLab merge request
	gitlabNoteURL := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests/%s/notes", *gitlabURL, *projectID, *mergeRequestIID)
	note := GitLabNote{Body: message.String()}
	noteJSON, err := json.Marshal(note)
	if err != nil {
		fmt.Println("Error marshaling GitLab note:", err)
		os.Exit(1)
	}

	gitlabReq, err := http.NewRequest("POST", gitlabNoteURL, bytes.NewBuffer(noteJSON))
	if err != nil {
		fmt.Println("Error creating GitLab request:", err)
		os.Exit(1)
	}
	gitlabReq.Header.Set("Content-Type", "application/json")
	gitlabReq.Header.Set("PRIVATE-TOKEN", *gitlabToken)

	gitlabResp, err := client.Do(gitlabReq)
	if err != nil {
		fmt.Println("Error making GitLab request:", err)
		os.Exit(1)
	}
	defer gitlabResp.Body.Close()

	if gitlabResp.StatusCode >= 200 && gitlabResp.StatusCode < 300 {
		fmt.Println("Successfully added note to GitLab merge request")
	} else {
		fmt.Printf("Failed to add note to GitLab merge request. Status code: %d\n", gitlabResp.StatusCode)
		gitlabBody, _ := ioutil.ReadAll(gitlabResp.Body)
		fmt.Println("GitLab response:", string(gitlabBody))
	}
}