package keys

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/17media/macgyver/cmd/crypto"
)

var (
	keys = map[Type]Keys{
		TypeText: &flagsKeys{},
		TypeEnv:  &envsKeys{},
		TypeFile: &fileKeys{},
	}
)

type envsKeys struct {
}

func (e *envsKeys) Import(input []string, secretTag string) []Key {
	var ks []Key
	envRegexp := `^(\w+)=(.*)$`
	reEnv := regexp.MustCompile(envRegexp)
	reSecret := getSecretRegexp(secretTag)
	for _, env := range input {
		key, value, err := getKVfromInput(env, reEnv)
		if err != nil {
			log.Printf("WARN: %s\n", err)
		}
		ks = append(ks, Key{
			Key:     key,
			Value:   value,
			Secrets: reSecret.parseValueToSecrets(value),
		})
	}
	return ks
}

func (e *envsKeys) Encrypt(input string, secretTag string, cp crypto.Crypto) map[string]interface{} {
	// Currently, it's only for file type Key
	log.Panic("not implement")
	return map[string]interface{}{}
}

func (e *envsKeys) Decrypt(input string, secretTag string, cp crypto.Crypto) map[string]interface{} {
	// Currently, it's only for file type Key
	log.Panic("not implement")
	return map[string]interface{}{}
}

func (e *envsKeys) ReplaceOriginFile(input string, values map[string]interface{}) error {
	// Currently, it's only for file type Key
	log.Panic("not implement")
	return nil
}

func (e *envsKeys) Export(keys []Key, secretTag string, writeCloser io.WriteCloser) error {
	exportStrs := ""
	reSecret := getSecretRegexp(secretTag)
	for _, k := range keys {
		newValue := reSecret.replaceSecrets(k.Value, k.Secrets)
		exportStrs += fmt.Sprintf("export %s='%s'\n", k.Key, newValue)
	}
	if _, err := writeCloser.Write([]byte(exportStrs)); err != nil {
		return err
	}
	return writeCloser.Close()
}

type flagsKeys struct {
}

func (f *flagsKeys) Encrypt(input string, secretTag string, cp crypto.Crypto) map[string]interface{} {
	// Currently, it's only for file type Key
	log.Panic("not implement")
	return map[string]interface{}{}
}

func (f *flagsKeys) Decrypt(input string, secretTag string, cp crypto.Crypto) map[string]interface{} {
	// Currently, it's only for file type Key
	log.Panic("not implement")
	return map[string]interface{}{}
}

func (f *flagsKeys) ReplaceOriginFile(input string, values map[string]interface{}) error {
	// Currently, it's only for file type Key
	log.Panic("not implement")
	return nil
}

func (f *flagsKeys) Import(input []string, secretTag string) []Key {
	var ks []Key
	flagRegexp := `^\-(\w+)=(.+)$`
	reFlag := regexp.MustCompile(flagRegexp)
	reSecret := getSecretRegexp(secretTag)
	for _, flag := range input {
		key, value, err := getKVfromInput(flag, reFlag)
		if err != nil {
			log.Printf("WARN: %s\n", err)
		}
		ks = append(ks, Key{
			Key:     key,
			Value:   value,
			Secrets: reSecret.parseValueToSecrets(value),
		})
	}
	return ks
}

func (f *flagsKeys) Export(keys []Key, secretTag string, writeCloser io.WriteCloser) error {
	exportFlags := ""
	reSecret := getSecretRegexp(secretTag)
	for _, k := range keys {
		newValue := reSecret.replaceSecrets(k.Value, k.Secrets)
		exportFlags += fmt.Sprintf(" -%s=%s", k.Key, newValue)
	}
	if _, err := writeCloser.Write([]byte(strings.TrimLeft(exportFlags, " ") + "\n")); err != nil {
		return err
	}
	return writeCloser.Close()
}

type fileKeys struct {
}

func (f *fileKeys) eval(input string, function func(s string) (interface{}, error)) map[string]interface{} {
	rawData, err := ioutil.ReadFile(input)
	if err != nil {
		log.Panicf("ReadFile failed %s", err)
	}
	data := convertYamlToMap(rawData)

	tmp, err := operationInMap(data, function)
	if err != nil {
		log.Panicf("convertToKeys failed %s", err)
	}
	return tmp.(map[string]interface{})
}

func (f *fileKeys) Encrypt(input string, secretTag string, cp crypto.Crypto) map[string]interface{} {
	reSecret := getSecretRegexp(secretTag)
	stringEncrypt := func(s string) (interface{}, error) {
		value := Key{
			Value:   s,
			Secrets: reSecret.parseValueToSecrets(s),
		}
		for _, s := range value.Secrets {
			if s.IsEncrypted {
				continue
			}
			encryptText, err := cp.Encrypt([]byte(s.Text))
			if err != nil {
				log.Panic(err)
			}
			s.Text = string(encryptText)
			s.IsEncrypted = true
		}
		return reSecret.replaceSecrets(value.Value, value.Secrets), nil
	}

	return f.eval(input, stringEncrypt)
}

func (f *fileKeys) Decrypt(input string, secretTag string, cp crypto.Crypto) map[string]interface{} {
	reSecret := getSecretRegexp(secretTag)
	stringEncrypt := func(s string) (interface{}, error) {
		value := Key{
			Value:   s,
			Secrets: reSecret.parseValueToSecrets(s),
		}
		for _, s := range value.Secrets {
			if !s.IsEncrypted {
				continue
			}
			encryptText, err := cp.Decrypt([]byte(s.Text))
			if err != nil {
				log.Panic(err)
			}
			s.Text = string(encryptText)
			s.IsEncrypted = false
		}
		return reSecret.replaceSecrets(value.Value, value.Secrets), nil
	}
	return f.eval(input, stringEncrypt)
}

func (f *fileKeys) ReplaceOriginFile(input string, values map[string]interface{}) error {
	yamlData, err := yaml.Marshal(&values)
	if err != nil {
		log.Panicf("Marshal yaml failed %s", err)
	}
	file, err := os.OpenFile(input, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	defer file.Close()
	if err != nil {
		log.Panic(err)
	}
	if _, err := file.Write(yamlData); err != nil {
		log.Panic(err)
	}
	return nil
}

func (f *fileKeys) Import(input []string, secretTag string) []Key {
	// It's for flag and env type
	log.Panic("Not implement")
	return []Key{}
}

func (f *fileKeys) Export(keys []Key, secretTag string, writeCloser io.WriteCloser) error {
	exportStrs := ""
	reSecret := getSecretRegexp(secretTag)
	for _, k := range keys {
		newValue := reSecret.replaceSecrets(k.Value, k.Secrets)
		exportStrs += fmt.Sprintf("export %s='%s'\n", k.Key, newValue)
	}
	if _, err := writeCloser.Write([]byte(exportStrs)); err != nil {
		return err
	}
	return writeCloser.Close()
}

func getKVfromInput(input string, re *regexp.Regexp) (key string, value string, err error) {
	kv := re.FindStringSubmatch(input)
	if len(kv) != 3 {
		emptyRegexp := `^\-(\w+)=$`
		emptyFlag := regexp.MustCompile(emptyRegexp)
		k := emptyFlag.FindStringSubmatch(input)
		return k[1], "", fmt.Errorf(`Cannot find value for key "%s"`, k[1])
	}
	key = kv[1]
	value = kv[2]
	return key, value, nil
}

func getSecretRegexp(secretTag string) *secretRegexp {
	secretTagRegexpTemplate := `<%s>(.*?)</%s>|<%s>(.*?)</%s>`
	secretTagRegexp := fmt.Sprintf(secretTagRegexpTemplate,
		strings.ToLower(secretTag), strings.ToLower(secretTag),
		strings.ToUpper(secretTag), strings.ToUpper(secretTag),
	)
	return &secretRegexp{
		secretTag: secretTag,
		re:        regexp.MustCompile(secretTagRegexp),
	}
}

type secretRegexp struct {
	secretTag string
	re        *regexp.Regexp
}

func (s *secretRegexp) parseValueToSecrets(value string) []*Secret {
	matchedSecrets := s.re.FindAllStringSubmatch(value, -1)
	// If we don't find any tagged secrets, we regards the entire value as a plaintext secret.
	if len(matchedSecrets) == 0 {
		return []*Secret{
			{
				Text:        value,
				IsEncrypted: false,
			},
		}
	}

	// Otherwise, returns secret converted from tagged secrets.
	var secrets []*Secret
	for _, ms := range matchedSecrets {
		plaintext := ms[1]
		ciphertext := ms[2]

		text := plaintext
		isEncrypted := false
		if ciphertext != "" {
			text = ciphertext
			isEncrypted = true
		}
		secrets = append(secrets, &Secret{
			Text:        text,
			IsEncrypted: isEncrypted,
		})
	}
	return secrets
}

func (s *secretRegexp) replaceSecrets(value string, secrets []*Secret) string {
	if len(secrets) == 0 {
		return ""
	}

	if len(s.re.FindAllStringSubmatch(value, -1)) == 0 {
		return s.outputSecret(secrets[0])
	}

	i := 0
	f := func(_ string) string {
		defer func() { i++ }()
		return s.outputSecret(secrets[i])
	}
	return s.re.ReplaceAllStringFunc(value, f)
}

func (s *secretRegexp) outputSecret(secret *Secret) string {
	text := secret.Text
	if secret.IsEncrypted == false {
		return text
	}
	return fmt.Sprintf("<%s>%s</%s>", strings.ToUpper(s.secretTag), text, strings.ToUpper(s.secretTag))
}
