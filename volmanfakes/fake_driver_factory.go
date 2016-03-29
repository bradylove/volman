// This file was generated by counterfeiter
package volmanfakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/cloudfoundry-incubator/volman/vollocal"
	"github.com/pivotal-golang/lager"
)

type FakeDriverFactory struct {
	DiscoverStub        func(logger lager.Logger) (map[string]string, error)
	discoverMutex       sync.RWMutex
	discoverArgsForCall []struct {
		logger lager.Logger
	}
	discoverReturns struct {
		result1 map[string]string
		result2 error
	}
	DriverStub        func(logger lager.Logger, driverId string) (voldriver.Driver, error)
	driverMutex       sync.RWMutex
	driverArgsForCall []struct {
		logger   lager.Logger
		driverId string
	}
	driverReturns struct {
		result1 voldriver.Driver
		result2 error
	}
}

func (fake *FakeDriverFactory) Discover(logger lager.Logger) (map[string]string, error) {
	fake.discoverMutex.Lock()
	fake.discoverArgsForCall = append(fake.discoverArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.discoverMutex.Unlock()
	if fake.DiscoverStub != nil {
		return fake.DiscoverStub(logger)
	} else {
		return fake.discoverReturns.result1, fake.discoverReturns.result2
	}
}

func (fake *FakeDriverFactory) DiscoverCallCount() int {
	fake.discoverMutex.RLock()
	defer fake.discoverMutex.RUnlock()
	return len(fake.discoverArgsForCall)
}

func (fake *FakeDriverFactory) DiscoverArgsForCall(i int) lager.Logger {
	fake.discoverMutex.RLock()
	defer fake.discoverMutex.RUnlock()
	return fake.discoverArgsForCall[i].logger
}

func (fake *FakeDriverFactory) DiscoverReturns(result1 map[string]string, result2 error) {
	fake.DiscoverStub = nil
	fake.discoverReturns = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeDriverFactory) Driver(logger lager.Logger, driverId string) (voldriver.Driver, error) {
	fake.driverMutex.Lock()
	fake.driverArgsForCall = append(fake.driverArgsForCall, struct {
		logger   lager.Logger
		driverId string
	}{logger, driverId})
	fake.driverMutex.Unlock()
	if fake.DriverStub != nil {
		return fake.DriverStub(logger, driverId)
	} else {
		return fake.driverReturns.result1, fake.driverReturns.result2
	}
}

func (fake *FakeDriverFactory) DriverCallCount() int {
	fake.driverMutex.RLock()
	defer fake.driverMutex.RUnlock()
	return len(fake.driverArgsForCall)
}

func (fake *FakeDriverFactory) DriverArgsForCall(i int) (lager.Logger, string) {
	fake.driverMutex.RLock()
	defer fake.driverMutex.RUnlock()
	return fake.driverArgsForCall[i].logger, fake.driverArgsForCall[i].driverId
}

func (fake *FakeDriverFactory) DriverReturns(result1 voldriver.Driver, result2 error) {
	fake.DriverStub = nil
	fake.driverReturns = struct {
		result1 voldriver.Driver
		result2 error
	}{result1, result2}
}

var _ vollocal.DriverFactory = new(FakeDriverFactory)