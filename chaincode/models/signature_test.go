package models

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSignatureModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Election Model Suite")
}

// 3045022100c08a66ef6aef6ccf6b806068a383683e6023cc488d9c541987c46aab0a0760f902203a07f1fd1ba0711572635b76e9c0cadae673c10e91e84b3e6010820218595436
