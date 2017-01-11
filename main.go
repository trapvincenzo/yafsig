package main

import (
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

var router *mux.Router

func createRouteHandler(route Route) {
	router.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
		switch route.ResponseType {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(route.Response))
		default:
			w.Header().Set("Content-type", "text/plain")
			w.Write([]byte(route.Response))
		}

	}).Methods(route.Methods...)
}

func main() {
	content, _ := ioutil.ReadFile("config.yaml")
	router = mux.NewRouter()

	routes := Routes{}
	serverConfig := ServerConfig{}

	yaml.Unmarshal(content, &routes)
	yaml.Unmarshal(content, &serverConfig)

	for _, route := range routes.Routes {
		createRouteHandler(route)
	}

	log.Fatal(http.ListenAndServe(serverConfig.Port, router))
}
