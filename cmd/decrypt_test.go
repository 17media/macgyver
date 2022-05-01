package cmd

import (
	"github.com/17media/macgyver/cmd/keys"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type decryptTestSuite struct {
	suite.Suite
}

func (s *decryptTestSuite) SetupSuite()    {}
func (s *decryptTestSuite) TearDownSuite() {}
func (s *decryptTestSuite) SetupTest()     {}
func (s *decryptTestSuite) TearDownTest()  {}

func (s *decryptTestSuite) TestDecrypt() {
	s.Run("decrypt text base64", func() {
		realStdout := os.Stdout
		defer func() { os.Stdout = realStdout }()
		r, fakeStdout, _ := os.Pipe()

		cryptoProvider = "base64"
		keysType = keys.TypeText
		flags = "-test=<SECRET_TAG>dGVzdA==</SECRET_TAG>"
		os.Stdout = fakeStdout

		decrypt(encryptCmd, []string{})
		_ = fakeStdout.Close()
		newOutBytes, _ := ioutil.ReadAll(r)
		s.Equal([]byte("-test=test\n"), newOutBytes)
	})
}

func TestDecryptSuite(t *testing.T) {
	suite.Run(t, new(decryptTestSuite))
}
