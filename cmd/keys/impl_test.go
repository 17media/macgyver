package keys

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type keysTestSuite struct {
	suite.Suite
}

func (s *keysTestSuite) SetupSuite()    {}
func (s *keysTestSuite) TearDownSuite() {}
func (s *keysTestSuite) SetupTest()     {}
func (s *keysTestSuite) TearDownTest()  {}

func (s *keysTestSuite) TestEnvsKeys() {
	tests := []struct {
		desc            string
		input           []string
		secretTag       string
		expImportedKeys []Key
		expExportedStr  string
	}{
		{
			desc: "test one secret plaintexts",
			input: []string{
				"ENV_1=AAA",
				"ENV_2=NotSecretPrefix<secret_tag>BBB</secret_tag>NotSecretSuffix",
				"ENV_3=<secret_tag>CCC</secret_tag>",
			},
			secretTag: "secret_tag",
			expImportedKeys: []Key{
				{
					Key:   "ENV_1",
					Value: "AAA",
					Secrets: []*Secret{
						{
							Text:        "AAA",
							IsEncrypted: false,
						},
					},
				},
				{
					Key:   "ENV_2",
					Value: "NotSecretPrefix<secret_tag>BBB</secret_tag>NotSecretSuffix",
					Secrets: []*Secret{
						{
							Text:        "BBB",
							IsEncrypted: false,
						},
					},
				},
				{
					Key:   "ENV_3",
					Value: "<secret_tag>CCC</secret_tag>",
					Secrets: []*Secret{
						{
							Text:        "CCC",
							IsEncrypted: false,
						},
					},
				},
			},
			expExportedStr: "export ENV_1='AAA'\nexport ENV_2='NotSecretPrefixBBBNotSecretSuffix'\nexport ENV_3='CCC'\n",
		},
		{
			desc: "test one secret ciphertexts",
			input: []string{
				"ENV_1=NotSecretPrefix<SECRET_TAG>AAA</SECRET_TAG>NotSecretSuffix",
				"ENV_2=<SECRET_TAG>BBB</SECRET_TAG>",
			},
			secretTag: "secret_tag",
			expImportedKeys: []Key{
				{
					Key:   "ENV_1",
					Value: "NotSecretPrefix<SECRET_TAG>AAA</SECRET_TAG>NotSecretSuffix",
					Secrets: []*Secret{
						{
							Text:        "AAA",
							IsEncrypted: true,
						},
					},
				},
				{
					Key:   "ENV_2",
					Value: "<SECRET_TAG>BBB</SECRET_TAG>",
					Secrets: []*Secret{
						{
							Text:        "BBB",
							IsEncrypted: true,
						},
					},
				},
			},
			expExportedStr: "export ENV_1='NotSecretPrefix<SECRET_TAG>AAA</SECRET_TAG>NotSecretSuffix'\nexport ENV_2='<SECRET_TAG>BBB</SECRET_TAG>'\n",
		},
		{
			desc: "test mixed secrets",
			input: []string{
				"MONGO_URI=mongo://<secret_tag>plaintext_userName</secret_tag>:<SECRET_TAG>ciphertext_password</SECRET_TAG>@1.2.3.4:<secret_tag>8080</secret_tag>/production",
			},
			secretTag: "secret_tag",
			expImportedKeys: []Key{
				{
					Key:   "MONGO_URI",
					Value: "mongo://<secret_tag>plaintext_userName</secret_tag>:<SECRET_TAG>ciphertext_password</SECRET_TAG>@1.2.3.4:<secret_tag>8080</secret_tag>/production",
					Secrets: []*Secret{
						{
							Text:        "plaintext_userName",
							IsEncrypted: false,
						},
						{
							Text:        "ciphertext_password",
							IsEncrypted: true,
						},
						{
							Text:        "8080",
							IsEncrypted: false,
						},
					},
				},
			},
			expExportedStr: "export MONGO_URI='mongo://plaintext_userName:<SECRET_TAG>ciphertext_password</SECRET_TAG>@1.2.3.4:8080/production'\n",
		},
	}

	envsKeys, err := Get(TypeEnv)
	s.NoError(err)
	for _, t := range tests {
		// test import
		keys := envsKeys.Import(t.input, t.secretTag)
		s.Equal(t.expImportedKeys, keys, "check imported keys", t.desc)

		// test export
		buf := bytes.NewBuffer([]byte{})
		s.NoError(envsKeys.Export(keys, t.secretTag, &nopWriterCloser{Writer: buf}), "check export error", t.desc)
		s.Equal(t.expExportedStr, buf.String(), "check exported string", t.desc)
	}
}

func (s *keysTestSuite) TestFlagsKeys() {
	tests := []struct {
		desc            string
		input           []string
		secretTag       string
		expImportedKeys []Key
		expExportedStr  string
	}{
		{
			desc: "test one secret plaintexts",
			input: []string{
				"-flag1=flag-AAA",
				"-flag2=NotSecretPrefix<secret_tag>flag-BBB</secret_tag>NotSecretSuffix",
				"-flag3=<secret_tag>flag-CCC</secret_tag>",
			},
			secretTag: "secret_tag",
			expImportedKeys: []Key{
				{
					Key:   "flag1",
					Value: "flag-AAA",
					Secrets: []*Secret{
						{
							Text:        "flag-AAA",
							IsEncrypted: false,
						},
					},
				},
				{
					Key:   "flag2",
					Value: "NotSecretPrefix<secret_tag>flag-BBB</secret_tag>NotSecretSuffix",
					Secrets: []*Secret{
						{
							Text:        "flag-BBB",
							IsEncrypted: false,
						},
					},
				},
				{
					Key:   "flag3",
					Value: "<secret_tag>flag-CCC</secret_tag>",
					Secrets: []*Secret{
						{
							Text:        "flag-CCC",
							IsEncrypted: false,
						},
					},
				},
			},
			expExportedStr: "-flag1=flag-AAA -flag2=NotSecretPrefixflag-BBBNotSecretSuffix -flag3=flag-CCC",
		},
		{
			desc: "test one secret ciphertexts",
			input: []string{
				"-flag1=NotSecretPrefix<SECRET_TAG>flag-AAA</SECRET_TAG>NotSecretSuffix",
				"-flag2=<SECRET_TAG>flag-BBB</SECRET_TAG>",
			},
			secretTag: "secret_tag",
			expImportedKeys: []Key{
				{
					Key:   "flag1",
					Value: "NotSecretPrefix<SECRET_TAG>flag-AAA</SECRET_TAG>NotSecretSuffix",
					Secrets: []*Secret{
						{
							Text:        "flag-AAA",
							IsEncrypted: true,
						},
					},
				},
				{
					Key:   "flag2",
					Value: "<SECRET_TAG>flag-BBB</SECRET_TAG>",
					Secrets: []*Secret{
						{
							Text:        "flag-BBB",
							IsEncrypted: true,
						},
					},
				},
			},
			expExportedStr: "-flag1=NotSecretPrefix<SECRET_TAG>flag-AAA</SECRET_TAG>NotSecretSuffix -flag2=<SECRET_TAG>flag-BBB</SECRET_TAG>",
		},
		{
			desc: "test mixed secrets",
			input: []string{
				"-mongoURI=mongo://<secret_tag>plaintext_userName</secret_tag>:<SECRET_TAG>ciphertext_password</SECRET_TAG>@1.2.3.4:<secret_tag>8080</secret_tag>/production",
			},
			secretTag: "secret_tag",
			expImportedKeys: []Key{
				{
					Key:   "mongoURI",
					Value: "mongo://<secret_tag>plaintext_userName</secret_tag>:<SECRET_TAG>ciphertext_password</SECRET_TAG>@1.2.3.4:<secret_tag>8080</secret_tag>/production",
					Secrets: []*Secret{
						{
							Text:        "plaintext_userName",
							IsEncrypted: false,
						},
						{
							Text:        "ciphertext_password",
							IsEncrypted: true,
						},
						{
							Text:        "8080",
							IsEncrypted: false,
						},
					},
				},
			},
			expExportedStr: "-mongoURI=mongo://plaintext_userName:<SECRET_TAG>ciphertext_password</SECRET_TAG>@1.2.3.4:8080/production",
		},
	}

	textKeys, err := Get(TypeText)
	s.NoError(err)
	for _, t := range tests {
		// test import
		keys := textKeys.Import(t.input, t.secretTag)
		s.Equal(t.expImportedKeys, keys, "check imported keys", t.desc)

		// test export
		buf := bytes.NewBuffer([]byte{})
		s.NoError(textKeys.Export(keys, t.secretTag, &nopWriterCloser{Writer: buf}), "check export error", t.desc)
		s.Equal(t.expExportedStr, buf.String(), "check exported string", t.desc)
	}
}

func (s *keysTestSuite) TestSecretRegexp() {
	re := getSecretRegexp("tag")
	tests := []struct {
		desc        string
		value       string
		expSecrets  []*Secret
		expNewValue string
	}{
		{
			desc: "any charaters without tag",
			value: `.*+?...secret-part1...    $^/\	secret-part2[{()}]$`,
			expSecrets: []*Secret{
				{
					Text: `.*+?...secret-part1...    $^/\	secret-part2[{()}]$`,
					IsEncrypted: false,
				},
			},
			expNewValue: `.*+?...secret-part1...    $^/\	secret-part2[{()}]$`,
		},
		{
			desc: "any charaters with plaintext",
			value: `.*+?...<tag>secret-part1...    $^/\	secret-part2</tag>[{()}]$`,
			expSecrets: []*Secret{
				{
					Text: `secret-part1...    $^/\	secret-part2`,
					IsEncrypted: false,
				},
			},
			expNewValue: `.*+?...secret-part1...    $^/\	secret-part2[{()}]$`,
		},
		{
			desc: "any charaters with ciphertext",
			value: `.*+?...<TAG>secret-part1...    tag$^/\	secret-part2</TAG>[{()}]$`,
			expSecrets: []*Secret{
				{
					Text: `secret-part1...    tag$^/\	secret-part2`,
					IsEncrypted: true,
				},
			},
			expNewValue: `.*+?...<TAG>secret-part1...    tag$^/\	secret-part2</TAG>[{()}]$`,
		},
	}
	for _, t := range tests {
		secrets := re.parseValueToSecrets(t.value)
		s.Equal(t.expSecrets, secrets, "check parsed secrets", t.desc)
		s.Equal(t.expNewValue, re.replaceSecrets(t.value, secrets), "check new value", t.desc)
	}
}

type nopWriterCloser struct {
	io.Writer
}

func (nwc *nopWriterCloser) Close() error {
	return nil
}

func TestKeysSuite(t *testing.T) {
	suite.Run(t, new(keysTestSuite))
}
