package vollocal_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/clock/fakeclock"
	"code.cloudfoundry.org/lager/lagertest"

	"code.cloudfoundry.org/clock"
	"github.com/cloudfoundry-incubator/voldriver"
	"github.com/cloudfoundry-incubator/voldriver/voldriverfakes"
	"github.com/cloudfoundry-incubator/volman/vollocal"
	"github.com/cloudfoundry-incubator/volman/volmanfakes"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var _ = Describe("Driver Syncer", func() {
	var (
		logger *lagertest.TestLogger

		scanInterval time.Duration

		fakeClock         *fakeclock.FakeClock
		fakeDriverFactory *volmanfakes.FakeDriverFactory
		fakeDriverSyncer  *volmanfakes.FakeDriverSyncer

		registry vollocal.DriverRegistry
		syncer   vollocal.DriverSyncer
		process  ifrit.Process
	)

	BeforeEach(func() {

		logger = lagertest.NewTestLogger("driver-syncer-test")

		fakeClock = fakeclock.NewFakeClock(time.Unix(123, 456))
		fakeDriverFactory = new(volmanfakes.FakeDriverFactory)
		fakeDriverSyncer = new(volmanfakes.FakeDriverSyncer)

		scanInterval = 10 * time.Second

		registry = vollocal.NewDriverRegistry()
		syncer = vollocal.NewDriverSyncerWithDriverFactory(logger, registry, []string{"/some/path"}, scanInterval, fakeClock, fakeDriverFactory)
	})

	Describe("#Runner", func() {
		It("has a non-nil runner", func() {
			Expect(syncer.Runner()).NotTo(BeNil())
		})

		It("has a non-nil and empty driver registry", func() {
			Expect(registry).NotTo(BeNil())
			Expect(len(registry.Drivers())).To(Equal(0))
		})
	})

	Describe("#Run", func() {
		Context("when there are no drivers", func() {
			It("should have no drivers in registry map", func() {
				drivers := registry.Drivers()
				Expect(len(drivers)).To(Equal(0))
				Expect(fakeDriverSyncer.DiscoverCallCount()).To(Equal(0))
				Expect(fakeDriverFactory.DriverCallCount()).To(Equal(0))
			})
		})

		Context("when there are drivers", func() {
			var (
				fakeDriver              *voldriverfakes.FakeDriver
				driverName              string
				fakeRemoteClientFactory *voldriverfakes.FakeRemoteClientFactory
				driverFactory           vollocal.DriverFactory
				syncer                  vollocal.DriverSyncer
			)

			BeforeEach(func() {
				driverName = "fakedriver"
				err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "spec", []byte("http://0.0.0.0:8080"))
				Expect(err).NotTo(HaveOccurred())
				fakeRemoteClientFactory = new(voldriverfakes.FakeRemoteClientFactory)
				driverFactory = vollocal.NewDriverFactoryWithRemoteClientFactory(fakeRemoteClientFactory)

				syncer = vollocal.NewDriverSyncerWithDriverFactory(logger, registry, []string{defaultPluginsDirectory}, scanInterval, fakeClock, driverFactory)

				fakeDriver = new(voldriverfakes.FakeDriver)
				fakeDriver.ActivateReturns(voldriver.ActivateResponse{
					Implements: []string{"VolumeDriver"},
				})

				fakeRemoteClientFactory.NewRemoteClientReturns(fakeDriver, nil)

				process = ginkgomon.Invoke(syncer.Runner())
			})

			AfterEach(func() {
				ginkgomon.Kill(process)
			})

			It("should have fake driver in registry map", func() {
				drivers := registry.Drivers()
				Expect(len(drivers)).To(Equal(1))
				Expect(fakeRemoteClientFactory.NewRemoteClientCallCount()).To(Equal(1))
				Expect(fakeDriver.ActivateCallCount()).To(Equal(1))
			})

			Context("when drivers are added", func() {
				BeforeEach(func() {
					err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, "anotherfakedriver", "spec", []byte("http://0.0.0.0:8080"))
					Expect(err).NotTo(HaveOccurred())
				})

				It("should find them!", func() {
					fakeClock.Increment(scanInterval * 2)
					Eventually(registry.Drivers).Should(HaveLen(2))
					Expect(fakeRemoteClientFactory.NewRemoteClientCallCount()).To(Equal(3))
					Expect(fakeDriver.ActivateCallCount()).To(Equal(3))
				})
			})
			Context("when drivers are not responding", func() {
				BeforeEach(func() {
					fakeDriver.ActivateReturns(voldriver.ActivateResponse{
						Err: "some err",
					})
				})

				It("should find no drivers", func() {
					fakeClock.Increment(scanInterval * 2)
					Eventually(registry.Drivers).Should(HaveLen(0))
				})
			})
		})

	})

	Describe("#Discover", func() {
		Context("when given driverspath with no drivers", func() {
			It("no drivers are found", func() {
				fakeRemoteClientFactory := new(voldriverfakes.FakeRemoteClientFactory)
				driverFactory := vollocal.NewDriverFactoryWithRemoteClientFactory(fakeRemoteClientFactory)
				driverSyncer := vollocal.NewDriverSyncerWithDriverFactory(logger, nil, []string{"some-invalid-drivers-path"}, time.Second*60, clock.NewClock(), driverFactory)
				drivers, err := driverSyncer.Discover(logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(drivers)).To(Equal(0))
			})
		})

		Context("with a single driver", func() {
			var (
				fakeDriver              *voldriverfakes.FakeDriver
				driverName              string
				fakeRemoteClientFactory *voldriverfakes.FakeRemoteClientFactory
				driverFactory           vollocal.DriverFactory
				driverSyncer            vollocal.DriverSyncer
			)

			BeforeEach(func() {
				driverName = "some-driver-name"
				err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "spec", []byte("http://0.0.0.0:8080"))
				Expect(err).NotTo(HaveOccurred())
				fakeRemoteClientFactory = new(voldriverfakes.FakeRemoteClientFactory)
				driverFactory = vollocal.NewDriverFactoryWithRemoteClientFactory(fakeRemoteClientFactory)
				driverSyncer = vollocal.NewDriverSyncerWithDriverFactory(logger, nil, []string{defaultPluginsDirectory}, time.Second*60, clock.NewClock(), driverFactory)

				fakeDriver = new(voldriverfakes.FakeDriver)
				fakeDriver.ActivateReturns(voldriver.ActivateResponse{
					Implements: []string{"VolumeDriver"},
				})

				fakeRemoteClientFactory.NewRemoteClientReturns(fakeDriver, nil)
			})

			It("should not find drivers that are unresponsive", func() {
				fakeDriver.ActivateReturns(voldriver.ActivateResponse{Err: "Error"})
				drivers, err := driverSyncer.Discover(logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(drivers)).To(Equal(0))
				Expect(fakeRemoteClientFactory.NewRemoteClientCallCount()).To(Equal(1))
			})

			It("should find drivers", func() {
				drivers, err := driverSyncer.Discover(logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(drivers)).To(Equal(1))
				Expect(fakeRemoteClientFactory.NewRemoteClientCallCount()).To(Equal(1))
			})
		})

		Context("when given a simple driverspath", func() {
			var (
				fakeRemoteClientFactory *voldriverfakes.FakeRemoteClientFactory
				driverFactory           vollocal.DriverFactory
				driverName              string
				driverSyncer            vollocal.DriverSyncer
			)

			Context("with hetergeneous driver specifications", func() {
				var fakeDriver *voldriverfakes.FakeDriver

				BeforeEach(func() {
					driverName = "some-driver-name"
					err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "json", []byte("{\"Addr\":\"http://0.0.0.0:8080\"}"))
					Expect(err).NotTo(HaveOccurred())
					err = voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "spec", []byte("http://0.0.0.0:9090"))
					Expect(err).NotTo(HaveOccurred())

					fakeRemoteClientFactory = new(voldriverfakes.FakeRemoteClientFactory)
					driverFactory = vollocal.NewDriverFactoryWithRemoteClientFactory(fakeRemoteClientFactory)

					fakeDriver = new(voldriverfakes.FakeDriver)
					fakeDriver.ActivateReturns(voldriver.ActivateResponse{
						Implements: []string{"VolumeDriver"},
					})

					fakeRemoteClientFactory.NewRemoteClientReturns(fakeDriver, nil)
					driverSyncer = vollocal.NewDriverSyncerWithDriverFactory(logger, nil, []string{defaultPluginsDirectory}, time.Second*60, clock.NewClock(), driverFactory)
				})

				It("should preferentially select spec over json specification", func() {
					drivers, err := driverSyncer.Discover(logger)
					Expect(err).ToNot(HaveOccurred())
					Expect(len(drivers)).To(Equal(1))
					actualAddress, _ := fakeRemoteClientFactory.NewRemoteClientArgsForCall(0)
					Expect(actualAddress).To(Equal("http://0.0.0.0:9090"))
				})
			})
		})

		Context("when given a compound driverspath", func() {
			var (
				fakeRemoteClientFactory *voldriverfakes.FakeRemoteClientFactory
				driverFactory           vollocal.DriverFactory
				fakeDriver              *voldriverfakes.FakeDriver
				driverSyncer            vollocal.DriverSyncer
				driverName              string
			)

			JustBeforeEach(func() {
				fakeRemoteClientFactory = new(voldriverfakes.FakeRemoteClientFactory)
				driverFactory = vollocal.NewDriverFactoryWithRemoteClientFactory(fakeRemoteClientFactory)
				fakeDriver = new(voldriverfakes.FakeDriver)
				fakeDriver.ActivateReturns(voldriver.ActivateResponse{
					Implements: []string{"VolumeDriver"},
				})

				fakeRemoteClientFactory.NewRemoteClientReturns(fakeDriver, nil)
				driverSyncer = vollocal.NewDriverSyncerWithDriverFactory(logger, nil, []string{defaultPluginsDirectory, secondPluginsDirectory}, time.Second*60, clock.NewClock(), driverFactory)
			})

			Context("with a single driver", func() {
				BeforeEach(func() {
					driverName = "some-driver-name"
					err := voldriver.WriteDriverSpec(logger, secondPluginsDirectory, driverName, "spec", []byte("http://0.0.0.0:8080"))
					Expect(err).NotTo(HaveOccurred())
				})

				It("should find drivers", func() {
					drivers, err := driverSyncer.Discover(logger)
					Expect(err).ToNot(HaveOccurred())
					Expect(len(drivers)).To(Equal(1))
					Expect(fakeRemoteClientFactory.NewRemoteClientCallCount()).To(Equal(1))
				})

			})

			Context("with multiple drivers in multiple directories", func() {

				BeforeEach(func() {
					driverName = "some-driver-name"
					err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "json", []byte("{\"Addr\":\"http://0.0.0.0:8080\"}"))
					Expect(err).NotTo(HaveOccurred())
					err = voldriver.WriteDriverSpec(logger, secondPluginsDirectory, "some-other-driver-name", "json", []byte("{\"Addr\":\"http://0.0.0.0:9090\"}"))
					Expect(err).NotTo(HaveOccurred())
				})

				It("should find both drivers", func() {
					drivers, err := driverSyncer.Discover(logger)
					Expect(err).ToNot(HaveOccurred())
					Expect(len(drivers)).To(Equal(2))
				})

			})

			Context("with the same driver but in multiple directories", func() {

				BeforeEach(func() {
					driverName = "some-driver-name"
					err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "json", []byte("{\"Addr\":\"http://0.0.0.0:8080\"}"))
					Expect(err).NotTo(HaveOccurred())
					err = voldriver.WriteDriverSpec(logger, secondPluginsDirectory, driverName, "spec", []byte("http://0.0.0.0:9090"))
					Expect(err).NotTo(HaveOccurred())
				})

				It("should preferentially select the driver in the first directory", func() {
					_, err := driverSyncer.Discover(logger)
					Expect(err).ToNot(HaveOccurred())
					actualAddress, _ := fakeRemoteClientFactory.NewRemoteClientArgsForCall(0)
					Expect(actualAddress).To(Equal("http://0.0.0.0:8080"))
				})
			})
		})

		Context("when given a driver spec not in canonical form", func() {
			var (
				fakeRemoteClientFactory *voldriverfakes.FakeRemoteClientFactory
				driverFactory           vollocal.DriverFactory
				fakeDriver              *voldriverfakes.FakeDriver
				driverName              string
				driverSyncer            vollocal.DriverSyncer
			)

			JustBeforeEach(func() {
				fakeRemoteClientFactory = new(voldriverfakes.FakeRemoteClientFactory)
				driverFactory = vollocal.NewDriverFactoryWithRemoteClientFactory(fakeRemoteClientFactory)
				driverSyncer = vollocal.NewDriverSyncerWithDriverFactory(logger, nil, []string{defaultPluginsDirectory}, time.Second*60, clock.NewClock(), driverFactory)
			})

			TestCanonicalization := func(context, actual, it, expected string) {
				Context(context, func() {
					BeforeEach(func() {
						driverName = "some-driver-name"
						err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "spec", []byte(actual))
						Expect(err).NotTo(HaveOccurred())
					})

					JustBeforeEach(func() {
						fakeDriver = new(voldriverfakes.FakeDriver)
						fakeDriver.ActivateReturns(voldriver.ActivateResponse{
							Implements: []string{"VolumeDriver"},
						})

						fakeRemoteClientFactory.NewRemoteClientReturns(fakeDriver, nil)
					})

					It(it, func() {
						drivers, err := driverSyncer.Discover(logger)
						Expect(err).ToNot(HaveOccurred())
						Expect(len(drivers)).To(Equal(1))
						Expect(fakeRemoteClientFactory.NewRemoteClientCallCount()).To(Equal(1))
						Expect(fakeRemoteClientFactory.NewRemoteClientArgsForCall(0)).To(Equal(expected))
					})
				})
			}

			TestCanonicalization("with an ip (and no port)", "127.0.0.1", "should return a canonicalized address", "http://127.0.0.1")
			TestCanonicalization("with an ip and port", "127.0.0.1:8080", "should return a canonicalized address", "http://127.0.0.1:8080")
			TestCanonicalization("with a tcp protocol uri with port", "tcp://127.0.0.1:8080", "should return a canonicalized address", "http://127.0.0.1:8080")
			TestCanonicalization("with a tcp protocol uri without port", "tcp://127.0.0.1", "should return a canonicalized address", "http://127.0.0.1")
			TestCanonicalization("with a unix address including protocol", "unix:///other.sock", "should return a canonicalized address", "unix:///other.sock")
			TestCanonicalization("with a unix address missing its protocol", "/other.sock", "should return a canonicalized address", "/other.sock")

			Context("with an invalid url", func() {
				BeforeEach(func() {
					driverName = "some-driver-name"
					err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "spec", []byte("htt%p:\\\\"))
					Expect(err).NotTo(HaveOccurred())
				})

				It("doesn't make a driver", func() {
					_, err := driverSyncer.Discover(logger)
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeRemoteClientFactory.NewRemoteClientCallCount()).To(Equal(0))
				})
			})
		})

		Context("when given a driver spec with a bad driver", func() {
			var (
				fakeRemoteClientFactory *voldriverfakes.FakeRemoteClientFactory
				driverFactory           vollocal.DriverFactory
				fakeDriver              *voldriverfakes.FakeDriver
				driverName              string
				driverSyncer            vollocal.DriverSyncer
			)

			JustBeforeEach(func() {
				fakeRemoteClientFactory = new(voldriverfakes.FakeRemoteClientFactory)
				driverFactory = vollocal.NewDriverFactoryWithRemoteClientFactory(fakeRemoteClientFactory)
				driverSyncer = vollocal.NewDriverSyncerWithDriverFactory(logger, nil, []string{defaultPluginsDirectory}, time.Second*60, clock.NewClock(), driverFactory)
			})

			BeforeEach(func() {
				driverName = "some-driver-name"
				err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "spec", []byte("127.0.0.1:8080"))
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return no drivers if the driver doesn't implement VolumeDriver", func() {
				fakeDriver = new(voldriverfakes.FakeDriver)
				fakeDriver.ActivateReturns(voldriver.ActivateResponse{
					Implements: []string{"something-else"},
				})

				fakeRemoteClientFactory.NewRemoteClientReturns(fakeDriver, nil)
				drivers, err := driverSyncer.Discover(logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(drivers)).To(Equal(0))
			})

			It("should return no drivers if the driver doesn't respond", func() {
				fakeDriver = new(voldriverfakes.FakeDriver)
				fakeDriver.ActivateReturns(voldriver.ActivateResponse{
					Err: "some-error",
				})

				fakeRemoteClientFactory.NewRemoteClientReturns(fakeDriver, nil)
				drivers, err := driverSyncer.Discover(logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(drivers)).To(Equal(0))

			})
		})
	})
})
