package syslog_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLearningRagel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Syslog Parser Suite")
}
