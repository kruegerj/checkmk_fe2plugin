package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type InputService struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	State string `json:"state"`
}
type InputServiceDetail struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	State   string `json:"state"`
}

type Config struct {
	Hostname string `yaml:"hostname"`
	Token    string `yaml:"token"`
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`
	ApiURL   string
}
type Amweb struct {
	Id               string `json:"identifier"`
	Name             string `json:"name"`
	Organisation     string `json:"organization"`
	ConnectionType   string `json:"connectionType"`
	ConnectionState  string `json:"connectionState"`
	ConnectionsCount int    `json:"nbrOfWebSocketConnections"`
}

type CloudService struct {
	Name  string `json:"service"`
	State string `json:"state"`
}

type RedundancyState struct {
	State      string `json:"state"`
	Current    string `json:"current"`
	Configured string `json:"configured"`
}

type Status struct {
	State             string          `json:"state"`
	Message           string          `json:"message"`
	NbrOfLoggedErrors int             `json:"nbrOfLoggedErrors"`
	RedundancyState   RedundancyState `json:"redundancyState"`
}

type Mqtt struct {
	Defaultbroker string `json:"defaultBroker"`
	Kubernetes    string `json:"kubernetes"`
}

func main() {
	// Konfiguration aus YAML-Datei lesen
	config := readConfig(getConfigFilePath())
	config.ApiURL = config.Protocol + "://" + config.Hostname + ":" + config.Port + "/rest/monitoring/"
	getinput(config)
	getAmWeb(config)
	getcloud(config)
	getstatus(config)
	getmqtt(config)
}

func readConfig(filename string) Config {
	// YAML-Datei öffnen
	file, err := os.Open(filename)
	if err != nil {
		log.WithError(err).Fatal("Error opening config file")
		os.Exit(1)
	}
	defer file.Close()

	// YAML-Datei parsen
	decoder := yaml.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		log.WithError(err).Fatal("Error decoding config file")
		os.Exit(1)
	}

	return config
}
func getConfigFilePath() string {
	// Pfad zum Ordner %ProgramData%\checkmk\agent\local
	agentLocalFolder := filepath.Join(os.Getenv("ProgramData"), "checkmk", "agent", "local")

	// Pfad zur Konfigurationsdatei im angegebenen Ordner
	configFilePath := filepath.Join(agentLocalFolder, "config.yaml")

	return configFilePath
}

// getMonitorIDs führt eine Anfrage durch, um die IDs zu erhalten
func getMonitorIDs(config Config) []string {
	apiURL := config.ApiURL + "input"
	// Erstellen Sie eine HTTP-Anfrage mit dem Authorization-Header
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
	}

	req.Header.Set("Authorization", config.Token)
	// Führen Sie die Anfrage durch
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
	}
	defer resp.Body.Close()
	// Überprüfen Sie den HTTP-Statuscode
	if resp.StatusCode != http.StatusOK {
		log.Println("Error:", resp.Status)
	}
	// Dekodieren Sie die JSON-Antwort
	var services []InputService
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		log.Println("Error decoding JSON:", err)
	}
	// Extrahiere die IDs aus dem Slice von Services
	var ids []string
	for _, service := range services {
		ids = append(ids, service.ID)
	}

	return ids
}

func getDetailedMonitorInfo(config Config, id string) InputServiceDetail {
	url := config.ApiURL + "input/" + id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
	}

	req.Header.Set("Authorization", config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error making HTTP request")
	}
	defer resp.Body.Close()

	var detailedInfo InputServiceDetail
	if err := json.NewDecoder(resp.Body).Decode(&detailedInfo); err != nil {
		log.WithError(err).Fatal("Error decoding JSON")
	}

	return detailedInfo
}

func getAmWeb(config Config) {
	apiURL := config.ApiURL + "amweb"
	// Erstellen Sie eine HTTP-Anfrage mit dem Authorization-Header
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
		os.Exit(1)
	}

	req.Header.Set("Authorization", config.Token)
	// Führen Sie die Anfrage durch
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
		return
	}
	defer resp.Body.Close()
	// Überprüfen Sie den HTTP-Statuscode
	if resp.StatusCode != http.StatusOK {
		log.Println("Error:", resp.Status)
		return
	}
	// Dekodieren Sie die JSON-Antwort
	var amwebs []Amweb
	if err := json.NewDecoder(resp.Body).Decode(&amwebs); err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	if len(amwebs) == 0 {
		return
	}
	// Detaillierte Informationen für jede ID abrufen
	// CheckMK-Ausgabe
	for _, amweb := range amwebs {
		amwebStatus := 0
		if amweb.ConnectionState != "OK" {
			amwebStatus = 1
		}
		fmt.Printf("%d \"AmWeb: %s\" connection=%d Organisation: %s ConnectionType: %s\n", amwebStatus, amweb.Name, amweb.ConnectionsCount, amweb.Organisation, amweb.ConnectionType)
	}
}

func getcloud(config Config) {
	apiURL := config.ApiURL + "cloud"
	// Erstellen Sie eine HTTP-Anfrage mit dem Authorization-Header
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
		return
	}

	req.Header.Set("Authorization", config.Token)
	// Führen Sie die Anfrage durch
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
		return
	}
	defer resp.Body.Close()
	// Überprüfen Sie den HTTP-Statuscode
	if resp.StatusCode != http.StatusOK {
		log.Println("Error:", resp.Status)
		return
	}
	// Dekodieren Sie die JSON-Antwort
	var services []CloudService
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	// Detaillierte Informationen für jede ID abrufen
	// CheckMK-Ausgabe
	for _, service := range services {
		serviceStatus := 0
		if service.State != "OK" {
			serviceStatus = 1
		}
		fmt.Printf("%d \"FE2 Cloud: %s\" myvalue=-\n", serviceStatus, service.Name)
	}
}

func getstatus(config Config) {
	apiURL := config.ApiURL + "status"
	// Erstellen Sie eine HTTP-Anfrage mit dem Authorization-Header
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
		return
	}

	req.Header.Set("Authorization", config.Token)
	// Führen Sie die Anfrage durch
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
		return
	}
	defer resp.Body.Close()
	// Überprüfen Sie den HTTP-Statuscode
	if resp.StatusCode != http.StatusOK {
		log.Println("Error:", resp.Status)
		return
	}
	// Dekodieren Sie die JSON-Antwort
	var services Status
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	// CheckMK-Ausgabe
	fmt.Printf("P \"FE2 Selfstatus\" errors=%d;1;5 %s\n", services.NbrOfLoggedErrors, services.Message)
}

func getmqtt(config Config) {
	apiURL := config.ApiURL + "mqtt"
	// Erstellen Sie eine HTTP-Anfrage mit dem Authorization-Header
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
		return
	}

	req.Header.Set("Authorization", config.Token)
	// Führen Sie die Anfrage durch
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error creating HTTP request")
		return
	}
	defer resp.Body.Close()
	// Überprüfen Sie den HTTP-Statuscode
	if resp.StatusCode != http.StatusOK {
		log.Println("Error:", resp.Status)
		return
	}
	// Dekodieren Sie die JSON-Antwort
	var services Mqtt
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	defaultstate := 0
	if services.Defaultbroker == "ERROR" {
		defaultstate = 1
	} else if services.Defaultbroker == "NOT_USED" {
		defaultstate = 3
	}
	kubernetesstate := 0
	if services.Defaultbroker == "ERROR" {
		kubernetesstate = 1
	} else if services.Defaultbroker == "NOT_USED" {
		kubernetesstate = 3
	}

	// CheckMK-Ausgabe
	fmt.Printf("%d \"FE2 MQTT Defaultbroker\" novalue=-\n", defaultstate)
	fmt.Printf("%d \"FE2 MQTT Kubernetes\" novalue=-\n", kubernetesstate)
}

func getinput(config Config) {
	ids := getMonitorIDs(config)
	// Detaillierte Informationen für jede ID abrufen
	// CheckMK-Ausgabe
	for _, id := range ids {
		detailedInfo := getDetailedMonitorInfo(config, id)
		serviceStatus := 0
		if detailedInfo.State != "OK" {
			serviceStatus = 1
		}
		if detailedInfo.Message == "" {
			detailedInfo.Message = "No Message available"
		}
		fmt.Printf("%d \"FE2 Input: %s\" myvalue=- %s\n", serviceStatus, detailedInfo.Name, detailedInfo.Message)
	}
}
