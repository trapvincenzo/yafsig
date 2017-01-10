package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// Routes struct to hold all the routes
type Routes struct {
	Routes []Route `yaml:"routes,flow"`
}

// Route struct to hold single route schema
type Route struct {
	Path         string   `yaml:"path"`
	Response     string   `yaml:"response"`
	Methods      []string `yaml:",flow"`
	ResponseType string   `yaml:"responseType,omitempty"`
}

// ServerConfig struct to hold server configurations
type ServerConfig struct {
	Port string `yaml:"port"`
}

func main() {
	content, _ := ioutil.ReadFile("config.yaml")

	routes := Routes{}
	serverConfig := ServerConfig{}

	yaml.Unmarshal(content, &routes)
	yaml.Unmarshal(content, &serverConfig)

	router := mux.NewRouter()
	for _, route := range routes.Routes {
		router.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			switch route.ResponseType {
			case "json":
				json.NewEncoder(w).Encode(route.Response)
			default:
				w.Write([]byte(route.Response))
			}

		}).Methods(route.Methods...)
	}

	log.Fatal(http.ListenAndServe(serverConfig.Port, router))
}
