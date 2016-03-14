package volman

import (
	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/tedsuo/rata"
)

const (
	ListDriversRoute = "drivers"
	MountRoute       = "mount"
	UnmountRoute     = "unmount"
	CreateRoute      = "create"
)

var Routes = rata.Routes{
	{Path: "/drivers", Method: "GET", Name: ListDriversRoute},
	{Path: "/drivers/mount", Method: "POST", Name: MountRoute},
	{Path: "/drivers/unmount", Method: "POST", Name: UnmountRoute},
	{Path: "/drivers/create", Method: "POST", Name: CreateRoute},
}

type ListDriversResponse struct {
	Drivers []voldriver.InfoResponse `json:"drivers"`
}

type MountRequest struct {
	DriverId string `json:"driverId"`
	VolumeId string `json:"volumeId"`
}

type CreateRequest struct {
	DriverId string                 `json:"driverId"`
	VolumeId string                 `json:"volumeId"`
	Opts     map[string]interface{} `json:"opts"`
}

type MountResponse struct {
	Path string `json:"path"`
}

type UnmountRequest struct {
	DriverId string `json:"driverId"`
	VolumeId string `json:"volumeId"`
}

func NewError(err error) Error {
	return Error{err.Error()}
}

type Error struct {
	Description string `json:"description"`
}

func (e Error) Error() string {
	return e.Description
}
