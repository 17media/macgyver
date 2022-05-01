package cmd

import (
	"github.com/17media/macgyver/cmd/keys"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type encryptTestSuite struct {
	suite.Suite
}

func (s *encryptTestSuite) SetupSuite()    {}
func (s *encryptTestSuite) TearDownSuite() {}
func (s *encryptTestSuite) SetupTest()     {}
func (s *encryptTestSuite) TearDownTest()  {}

func (s *encryptTestSuite) TestEncrypt() {
	s.Run("encrypt text base64", func() {
		realStdout := os.Stdout
		defer func() { os.Stdout = realStdout }()
		r, fakeStdout, _ := os.Pipe()

		cryptoProvider = "base64"
		keysType = keys.TypeText
		flags = "-test=test"
		os.Stdout = fakeStdout

		encrypt(encryptCmd, []string{})
		_ = fakeStdout.Close()
		newOutBytes, _ := ioutil.ReadAll(r)
		s.Equal([]byte("-test=<SECRET_TAG>dGVzdA==</SECRET_TAG>\n"), newOutBytes)
	})
}

func TestEncryptSuite(t *testing.T) {
	suite.Run(t, new(encryptTestSuite))
}
