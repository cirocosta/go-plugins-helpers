package logging

import (
	"log"
	"net/http"

	"github.com/docker/docker/daemon/logger"
	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	manifest         = `{"Implements": ["LogDriver"]}`
	startLoggingPath = "/LogDriver.StartLogging"
	stopLoggingPath  = "/LogDriver.StopLogging"
	readLogsPath     = "/LogDriver.ReadLogs"
	capabilitiesPath = "/LogDriver.Capabilities"
)

// CreateRequest is the structure that docker's requests are deserialized to.
type StartLoggingRequest struct {
	File string
	Info logger.Info
}

// StopLoggingRequest is the structure that docker's requests are deserialized to.
type StopLoggingRequest struct {
	File string
}

// ReadLogsRequest is the structure that docker's requests are deserialized to.
// TODO
type ReadLogsRequest struct{
  ReadConfig logger.ReadConfig
  Info logger.Info
}

// CapabilitiesResponse structure for logging capability response
type CapabilitiesResponse struct {
	ReadLogs bool
}

// ErrorResponse is a formatted error message that docker can understand
type ErrorResponse struct {
	Err string
}

// NewErrorResponse creates an ErrorResponse with the provided message
func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Err: msg}
}

// Driver represent the interface a driver must fulfill.
type Driver interface {
	StartLogging(*StartLoggingRequest) error
	StopLogging(*StopLoggingRequest) error
	ReadLogs(*ReadLogsRequest) error
	Capabilities() *CapabilitiesResponse
}

// Handler forwards requests and responses between the docker daemon and the plugin.
type Handler struct {
	driver Driver
	sdk.Handler
}

// NewHandler initializes the request handler with a driver implementation.
func NewHandler(driver Driver) *Handler {
	h := &Handler{driver, sdk.NewHandler(manifest)}
	h.initMux()
	return h
}

func (h *Handler) initMux() {
	h.HandleFunc(startLoggingPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers startLoggingPath")

		req := &StartLoggingRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.StartLogging(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})

	h.HandleFunc(stopLoggingPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers stopLoggingPath")
		req := &StopLoggingRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.StopLogging(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, struct{}{}, false)
	})

	h.HandleFunc(readLogsPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers readLogsPath")
		req := &ReadLogsRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.ReadLogs(req)
		if err != nil {
			sdk.EncodeResponse(w, NewErrorResponse(err.Error()), true)
			return
		}
		sdk.EncodeResponse(w, res, false)
	})

	h.HandleFunc(capabilitiesPath, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Entering go-plugins-helpers capabilitiesPath")
		sdk.EncodeResponse(w, h.driver.Capabilities(), false)
	})
}
