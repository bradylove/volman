package volman_test

import (
	. "github.com/cloudfoundry-incubator/volman"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var client Client

func TestLocalClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Local Client Suite")
}