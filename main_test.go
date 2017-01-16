package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRoute struct {
	Path     string
	Callback func(http.ResponseWriter, *http.Request)
}

type mockRouter struct {
	mock.Mock
	Registered []mockRoute
}

type mockHTTP struct {
	mock.Mock
}

func (m *mockHTTP) ListenAndServe(addr string, handler http.Handler) error {
	m.Called(addr, handler)
	return nil
}

func (m *mockRouter) HandleFunc(path string, handleFunc func(http.ResponseWriter, *http.Request)) *mux.Route {
	m.Called(path, handleFunc)
	m.Registered = append(m.Registered, mockRoute{Path: path, Callback: handleFunc})
	return &mux.Route{}
}

func (m *mockRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func TestDefaultConfigIsLoadedWhenNotSpecified(t *testing.T) {
	router := &mockRouter{}
	http := &mockHTTP{}

	router.On("HandleFunc", "/test", mock.AnythingOfType("func(http.ResponseWriter, *http.Request)")).Return(nil)
	router.On("HandleFunc", "/test/{id}", mock.AnythingOfType("func(http.ResponseWriter, *http.Request)")).Return(nil)
	http.On("ListenAndServe", ":12345", router).Return(nil)

	params := &FakeServerParams{Router: router, HTTP: http}
	NewFakeServer(params)

	assert.Equal(t, "config.yaml.dist", params.ConfigFilename)
}

func TestNewFakeServer(t *testing.T) {
	router := &mockRouter{}
	http := &mockHTTP{}

	router.On("HandleFunc", "/test", mock.AnythingOfType("func(http.ResponseWriter, *http.Request)")).Return(nil)
	router.On("HandleFunc", "/test/{id}", mock.AnythingOfType("func(http.ResponseWriter, *http.Request)")).Return(nil)
	http.On("ListenAndServe", ":12345", router).Return(nil)

	NewFakeServer(&FakeServerParams{Router: router, HTTP: http, ConfigFilename: "fixtures/config.yaml"})
	router.AssertExpectations(t)
	http.AssertExpectations(t)

	assert.Len(t, router.Registered, 2, "Lenght of the registerd routes is not correct")
	assert.Equal(t, "/test/{id}", router.Registered[0].Path)
	assert.Equal(t, "/test", router.Registered[1].Path)

	w := httptest.NewRecorder()
	router.Registered[0].Callback(w, nil)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	router.Registered[1].Callback(w, nil)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))

}
