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

type Amweb struct {
	Id               string `json:"identifier"`
	Name             string `json:"name"`
	Organisation     string `json:"organization"`
	ConnectionType   string `json:"connectionType"`
	ConnectionState  string `json:"connectionState"`
	ConnectionsCount int    `json:"nbrOfWebSocketConnections"`
}

type Config struct {
	Hostname string `yaml:"hostname"`
	Token    string `yaml:"token"`
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

func main() {
	// Konfiguration aus YAML-Datei lesen
	config := readConfig(getConfigFilePath())

	apiURL := config.Protocol + "://" + config.Hostname + ":" + config.Port + "/rest/monitoring/amweb"
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
		os.Exit(1)
	}
	defer resp.Body.Close()
	// Überprüfen Sie den HTTP-Statuscode
	if resp.StatusCode != http.StatusOK {
		log.Println("Error:", resp.Status)
		os.Exit(1)
	}
	// Dekodieren Sie die JSON-Antwort
	var amwebs []Amweb
	if err := json.NewDecoder(resp.Body).Decode(&amwebs); err != nil {
		log.Println("Error decoding JSON:", err)
		os.Exit(1)
	}
	if len(amwebs) == 0 {
		os.Exit(1)
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

func readConfig(filename string) Config {
	// YAML-Datei öffnen
	file, err := os.Open(filename)
	if err != nil {
		log.WithError(err).Fatal("Error opening config file")
	}
	defer file.Close()

	// YAML-Datei parsen
	decoder := yaml.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		log.WithError(err).Fatal("Error decoding config file")
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
