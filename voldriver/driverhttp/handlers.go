package driverhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cf_http_handlers "github.com/cloudfoundry-incubator/cf_http/handlers"
	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"
)

func respondWithError(logger lager.Logger, info string, err error, w http.ResponseWriter) {
	logger.Error(info, err)
	cf_http_handlers.WriteJSONResponse(w, http.StatusInternalServerError, voldriver.ErrorResponse{Err: err.Error()})
}

func NewHandler(logger lager.Logger, client voldriver.Driver) (http.Handler, error) {
	logger = logger.Session("server")
	logger.Info("start")
	defer logger.Info("end")
	var handlers = rata.Handlers{
		"mount": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			volumeID := ""
			logger.Info("mount")
			defer logger.Info("mount end")
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				respondWithError(logger, "Error reading mount request body", err, w)
				return
			}

			var mountRequest voldriver.MountRequest
			if err = json.Unmarshal(body, &mountRequest); err != nil {
				respondWithError(logger, fmt.Sprintf("Error reading mount request body: %#v", body), err, w)
				return
			}

			mountResponse := client.Mount(logger, mountRequest)
			if err != nil {
				respondWithError(logger, fmt.Sprintf("Error mounting volume %s", volumeID), err, w)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, http.StatusOK, mountResponse)
		}),

		"unmount": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger.Info("unmount")
			defer logger.Info("unmount end")
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				respondWithError(logger, "Error reading unmount request body", err, w)
				return
			}

			var unmountRequest voldriver.UnmountRequest
			if err = json.Unmarshal(body, &unmountRequest); err != nil {
				respondWithError(logger, fmt.Sprintf("Error reading unmount request body: %#v", body), err, w)
				return
			}

			err = client.Unmount(logger, unmountRequest)
			if err != nil {
				respondWithError(logger, fmt.Sprintf("Error unmounting volume %s", unmountRequest.VolumeId), err, w)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, http.StatusOK, struct{}{})
		}),

		"create": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger.Info("create")
			defer logger.Info("create end")
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				respondWithError(logger, "Error reading create request body", err, w)
				return
			}

			var createRequest voldriver.CreateRequest
			if err = json.Unmarshal(body, &createRequest); err != nil {
				respondWithError(logger, fmt.Sprintf("Error reading create request body: %#v", body), err, w)
				return
			}

			createResponse := client.Create(logger, createRequest)
			if createResponse.Err != "" {
				respondWithError(logger, fmt.Sprintf("Error creating volume %s", createRequest.Name), err, w)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, http.StatusOK, struct{}{})
		}),
	}

	return rata.NewRouter(voldriver.Routes, handlers)
}
