package keys

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/17media/macgyver/cmd/crypto"
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
			expExportedStr: "-flag1=flag-AAA -flag2=NotSecretPrefixflag-BBBNotSecretSuffix -flag3=flag-CCC\n",
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
			expExportedStr: "-flag1=NotSecretPrefix<SECRET_TAG>flag-AAA</SECRET_TAG>NotSecretSuffix -flag2=<SECRET_TAG>flag-BBB</SECRET_TAG>\n",
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
			expExportedStr: "-mongoURI=mongo://plaintext_userName:<SECRET_TAG>ciphertext_password</SECRET_TAG>@1.2.3.4:8080/production\n",
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

func (s *keysTestSuite) TestFileKeysEncrypt() {
	tests := []struct {
		desc      string
		input     []byte
		metadata  map[string]interface{}
		secretTag string
		expOutput []byte
	}{
		{
			desc: "test single",
			input: []byte(`test: aa
`),
			metadata: map[string]interface{}{
				"test": "<SECRET_TAG>YWE=</SECRET_TAG>",
			},
			secretTag: "secret_tag",
			expOutput: []byte(
				`test: <SECRET_TAG>YWE=</SECRET_TAG>
`,
			),
		},
		{
			desc: "test nested struct and multi secrets",
			input: []byte(
				`test: 1
test2:
   abc: "aa "
   c: true
   d: 3.222
   def: mongo://plaintext_userName:ciphertext_password@1.2.3.4:8080/production
   test3:
       e:
           - aaa/
           - 232
           - 'dcsc: s'
       f: asc
       t: acc
`,
			),
			metadata: map[string]interface{}{
				"test": 1,
				"test2": map[string]interface{}{
					"abc": "<SECRET_TAG>YWEg</SECRET_TAG>",
					"c":   true,
					"d":   3.222,
					"def": "<SECRET_TAG>bW9uZ286Ly9wbGFpbnRleHRfdXNlck5hbWU6Y2lwaGVydGV4dF9wYXNzd29yZEAxLjIuMy40OjgwODAvcHJvZHVjdGlvbg==</SECRET_TAG>",
					"test3": map[string]interface{}{
						"e": []interface{}{
							"<SECRET_TAG>YWFhLw==</SECRET_TAG>",
							232,
							"<SECRET_TAG>ZGNzYzogcw==</SECRET_TAG>",
						},
						"f": "<SECRET_TAG>YXNj</SECRET_TAG>",
						"t": "<SECRET_TAG>YWNj</SECRET_TAG>",
					},
				},
			},
			secretTag: "secret_tag",
			expOutput: []byte(
				`test: 1
test2:
    abc: <SECRET_TAG>YWEg</SECRET_TAG>
    c: true
    d: 3.222
    def: <SECRET_TAG>bW9uZ286Ly9wbGFpbnRleHRfdXNlck5hbWU6Y2lwaGVydGV4dF9wYXNzd29yZEAxLjIuMy40OjgwODAvcHJvZHVjdGlvbg==</SECRET_TAG>
    test3:
        e:
            - <SECRET_TAG>YWFhLw==</SECRET_TAG>
            - 232
            - <SECRET_TAG>ZGNzYzogcw==</SECRET_TAG>
        f: <SECRET_TAG>YXNj</SECRET_TAG>
        t: <SECRET_TAG>YWNj</SECRET_TAG>
`,
			),
		},
	}

	fileKeys, err := Get(TypeFile)
	crypto.Init("base64")
	cp := crypto.Providers["base64"]
	s.NoError(err)
	for _, t := range tests {
		s.Run(t.desc, func() {
			fileName, err := RandomCreateFile(t.input)
			s.NoError(err)

			meta := fileKeys.Encrypt(fileName, t.secretTag, cp)
			s.Equal(t.metadata, meta)

			err = fileKeys.ReplaceOriginFile(fileName, meta)
			s.NoError(err)

			actOutput, err := ioutil.ReadFile(fileName)
			s.NoError(err)
			s.Equal(t.expOutput, actOutput)

			err = RemoveFile(fileName)
			s.NoError(err)
		})
	}
}

func (s *keysTestSuite) TestFileKeysDecrypt() {
	tests := []struct {
		desc      string
		input     []byte
		metadata  map[string]interface{}
		secretTag string
		expOutput []byte
	}{
		{
			desc: "test single",
			input: []byte(`test: <SECRET_TAG>YWEg</SECRET_TAG>
`),
			metadata: map[string]interface{}{
				"test": "aa ",
			},
			secretTag: "secret_tag",
			expOutput: []byte(
				`test: 'aa '
`,
			),
		},
		{
			desc: "test nested struct and multi secrets",
			input: []byte(
				`test: 1
test2:
   abc: <SECRET_TAG>YWE=</SECRET_TAG>
   c: true
   d: 3.222
   def: <SECRET_TAG>bW9uZ286Ly9wbGFpbnRleHRfdXNlck5hbWU6Y2lwaGVydGV4dF9wYXNzd29yZEAxLjIuMy40OjgwODAvcHJvZHVjdGlvbg==</SECRET_TAG>
   test3:
       e:
           - <SECRET_TAG>YWFhLw==</SECRET_TAG>
           - 232
           - <SECRET_TAG>ZGNzYzogcw==</SECRET_TAG>
       f: <SECRET_TAG>YXNj</SECRET_TAG>
       t: <SECRET_TAG>YWNj</SECRET_TAG>
`,
			),
			metadata: map[string]interface{}{
				"test": 1,
				"test2": map[string]interface{}{
					"abc": "aa",
					"c":   true,
					"d":   3.222,
					"def": "mongo://plaintext_userName:ciphertext_password@1.2.3.4:8080/production",
					"test3": map[string]interface{}{
						"e": []interface{}{
							"aaa/",
							232,
							"dcsc: s",
						},
						"f": "asc",
						"t": "acc",
					},
				},
			},
			secretTag: "secret_tag",
			expOutput: []byte(
				`test: 1
test2:
    abc: aa
    c: true
    d: 3.222
    def: mongo://plaintext_userName:ciphertext_password@1.2.3.4:8080/production
    test3:
        e:
            - aaa/
            - 232
            - 'dcsc: s'
        f: asc
        t: acc
`,
			),
		},
	}

	fileKeys, err := Get(TypeFile)
	crypto.Init("base64")
	cp := crypto.Providers["base64"]
	s.NoError(err)
	for _, t := range tests {
		s.Run(t.desc, func() {
			fileName, err := RandomCreateFile(t.input)
			s.NoError(err)

			meta := fileKeys.Decrypt(fileName, t.secretTag, cp)
			s.Equal(t.metadata, meta)

			err = fileKeys.ReplaceOriginFile(fileName, meta)
			s.NoError(err)

			actOutput, err := ioutil.ReadFile(fileName)
			s.NoError(err)
			s.Equal(t.expOutput, actOutput)

			err = RemoveFile(fileName)
			s.NoError(err)
		})
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
