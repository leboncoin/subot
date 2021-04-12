package handler_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestReplierHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReplierHandler Suite")
}
