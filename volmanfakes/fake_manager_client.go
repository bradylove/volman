// This file was generated by counterfeiter
package volmanfakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/volman"
	"github.com/pivotal-golang/lager"
)

type FakeManager struct {
	ListDriversStub        func(logger lager.Logger) (volman.ListDriversResponse, error)
	listDriversMutex       sync.RWMutex
	listDriversArgsForCall []struct {
		logger lager.Logger
	}
	listDriversReturns struct {
		result1 volman.ListDriversResponse
		result2 error
	}
	MountStub        func(logger lager.Logger, driverId string, volumeId string, config string) (volman.MountResponse, error)
	mountMutex       sync.RWMutex
	mountArgsForCall []struct {
		logger   lager.Logger
		driverId string
		volumeId string
		config   string
	}
	mountReturns struct {
		result1 volman.MountResponse
		result2 error
	}
	UnmountStub        func(logger lager.Logger, driverId string, volumeId string) error
	unmountMutex       sync.RWMutex
	unmountArgsForCall []struct {
		logger   lager.Logger
		driverId string
		volumeId string
	}
	unmountReturns struct {
		result1 error
	}
}

func (fake *FakeManager) ListDrivers(logger lager.Logger) (volman.ListDriversResponse, error) {
	fake.listDriversMutex.Lock()
	fake.listDriversArgsForCall = append(fake.listDriversArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.listDriversMutex.Unlock()
	if fake.ListDriversStub != nil {
		return fake.ListDriversStub(logger)
	} else {
		return fake.listDriversReturns.result1, fake.listDriversReturns.result2
	}
}

func (fake *FakeManager) ListDriversCallCount() int {
	fake.listDriversMutex.RLock()
	defer fake.listDriversMutex.RUnlock()
	return len(fake.listDriversArgsForCall)
}

func (fake *FakeManager) ListDriversArgsForCall(i int) lager.Logger {
	fake.listDriversMutex.RLock()
	defer fake.listDriversMutex.RUnlock()
	return fake.listDriversArgsForCall[i].logger
}

func (fake *FakeManager) ListDriversReturns(result1 volman.ListDriversResponse, result2 error) {
	fake.ListDriversStub = nil
	fake.listDriversReturns = struct {
		result1 volman.ListDriversResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) Mount(logger lager.Logger, driverId string, volumeId string, config string) (volman.MountResponse, error) {
	fake.mountMutex.Lock()
	fake.mountArgsForCall = append(fake.mountArgsForCall, struct {
		logger   lager.Logger
		driverId string
		volumeId string
		config   string
	}{logger, driverId, volumeId, config})
	fake.mountMutex.Unlock()
	if fake.MountStub != nil {
		return fake.MountStub(logger, driverId, volumeId, config)
	} else {
		return fake.mountReturns.result1, fake.mountReturns.result2
	}
}

func (fake *FakeManager) MountCallCount() int {
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	return len(fake.mountArgsForCall)
}

func (fake *FakeManager) MountArgsForCall(i int) (lager.Logger, string, string, string) {
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	return fake.mountArgsForCall[i].logger, fake.mountArgsForCall[i].driverId, fake.mountArgsForCall[i].volumeId, fake.mountArgsForCall[i].config
}

func (fake *FakeManager) MountReturns(result1 volman.MountResponse, result2 error) {
	fake.MountStub = nil
	fake.mountReturns = struct {
		result1 volman.MountResponse
		result2 error
	}{result1, result2}
}

func (fake *FakeManager) Unmount(logger lager.Logger, driverId string, volumeId string) error {
	fake.unmountMutex.Lock()
	fake.unmountArgsForCall = append(fake.unmountArgsForCall, struct {
		logger   lager.Logger
		driverId string
		volumeId string
	}{logger, driverId, volumeId})
	fake.unmountMutex.Unlock()
	if fake.UnmountStub != nil {
		return fake.UnmountStub(logger, driverId, volumeId)
	} else {
		return fake.unmountReturns.result1
	}
}

func (fake *FakeManager) UnmountCallCount() int {
	fake.unmountMutex.RLock()
	defer fake.unmountMutex.RUnlock()
	return len(fake.unmountArgsForCall)
}

func (fake *FakeManager) UnmountArgsForCall(i int) (lager.Logger, string, string) {
	fake.unmountMutex.RLock()
	defer fake.unmountMutex.RUnlock()
	return fake.unmountArgsForCall[i].logger, fake.unmountArgsForCall[i].driverId, fake.unmountArgsForCall[i].volumeId
}

func (fake *FakeManager) UnmountReturns(result1 error) {
	fake.UnmountStub = nil
	fake.unmountReturns = struct {
		result1 error
	}{result1}
}

var _ volman.Manager = new(FakeManager)
