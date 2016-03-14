package volhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	cf_http_handlers "github.com/cloudfoundry-incubator/cf_http/handlers"
	"github.com/cloudfoundry-incubator/volman"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"
)

func respondWithError(logger lager.Logger, info string, err error, w http.ResponseWriter) {
	logger.Error(info, err)
	cf_http_handlers.WriteJSONResponse(w, http.StatusInternalServerError, volman.NewError(err))
}

func NewHandler(logger lager.Logger, client volman.Manager) (http.Handler, error) {
	logger = logger.Session("server")
	logger.Info("start")
	defer logger.Info("end")
	var handlers = rata.Handlers{
		"drivers": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			drivers, _ := client.ListDrivers(logger)
			cf_http_handlers.WriteJSONResponse(w, http.StatusOK, drivers)
		}),
		"mount": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger.Info("mount")
			defer logger.Info("mount end")
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				respondWithError(logger, "Error reading mount request body", err, w)
				return
			}

			var mountRequest volman.MountRequest
			if err = json.Unmarshal(body, &mountRequest); err != nil {
				respondWithError(logger, fmt.Sprintf("Error reading mount request body: %#v", body), err, w)
				return
			}

			mountPoint, err := client.Mount(logger, mountRequest.DriverId, mountRequest.VolumeId, "fix me")
			if err != nil {
				respondWithError(logger, fmt.Sprintf("Error mounting volume %s with driver %s", mountRequest.VolumeId, mountRequest.DriverId), err, w)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, http.StatusOK, mountPoint)
		}),
		"unmount": http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger.Info("unmount")
			defer logger.Info("unmount end")
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				respondWithError(logger, "Error reading unmount request body", err, w)
				return
			}

			var unmountRequest volman.UnmountRequest
			if err = json.Unmarshal(body, &unmountRequest); err != nil {
				respondWithError(logger, fmt.Sprintf("Error reading unmount request body: %#v", body), err, w)
				return
			}

			err = client.Unmount(logger, unmountRequest.DriverId, unmountRequest.VolumeId)
			if err != nil {
				respondWithError(logger, fmt.Sprintf("Error unmounting volume %s with driver %s", unmountRequest.VolumeId, unmountRequest.DriverId), err, w)
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

			var createRequest volman.CreateRequest
			if err = json.Unmarshal(body, &createRequest); err != nil {
				respondWithError(logger, fmt.Sprintf("Error reading create request body: %#v", body), err, w)
				return
			}

			err = client.Create(logger, createRequest.DriverId, createRequest.VolumeId, createRequest.Opts)
			if err != nil {
				respondWithError(logger, fmt.Sprintf("Error creating volume %s with driver %s", createRequest.VolumeId, createRequest.DriverId), err, w)
				return
			}

			cf_http_handlers.WriteJSONResponse(w, http.StatusOK, struct{}{})
		}),
	}

	return rata.NewRouter(volman.Routes, handlers)
}
