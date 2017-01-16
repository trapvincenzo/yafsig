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

// FakeServerParams config
type FakeServerParams struct {
	Router         RouterInterface
	ConfigFilename string
	HTTP           HTTPServer
}

// RouterInterface struct to be implemented by the routers
type RouterInterface interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request)) *mux.Route
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// HTTPServer struct to be implemented by the servers
type HTTPServer interface {
	ListenAndServe(string, http.Handler) error
}

var router RouterInterface

const defaultConfig = "config.yaml.dist"

// NewFakeServer creates a new fake router
func NewFakeServer(p *FakeServerParams) {

	router = mux.NewRouter()

	if p.Router != nil {
		router = p.Router
	}

	if p.ConfigFilename == "" {
		p.ConfigFilename = defaultConfig
	}

	content, err := ioutil.ReadFile(p.ConfigFilename)

	if err == nil {
		routes := Routes{}
		serverConfig := ServerConfig{}

		yaml.Unmarshal(content, &routes)
		yaml.Unmarshal(content, &serverConfig)

		for _, route := range routes.Routes {
			createRouteHandler(route)
		}

		if p.HTTP != nil {
			p.HTTP.ListenAndServe(serverConfig.Port, router)
			return
		}

		log.Fatal(http.ListenAndServe(serverConfig.Port, router))
	}
}

func createRouteHandler(route Route) {
	router.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
		switch route.ResponseType {
		case "json":
			w.Header().Set("Content-Type", "application/json")
		default:
			w.Header().Set("Content-type", "text/plain")
		}

		w.Write([]byte(route.Response))
	}).Methods(route.Methods...)
}

func main() {
	NewFakeServer(&FakeServerParams{})
}
