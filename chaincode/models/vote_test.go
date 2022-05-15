package models

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVoteModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Election Model Suite")
}
