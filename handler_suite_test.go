package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestParsing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parsing Suite")
}
