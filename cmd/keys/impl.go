package keys

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	keys = map[Type]Keys{
		TypeText: &flagsKeys{},
		TypeEnv:  &envsKeys{},
	}
)

type envsKeys struct {
}

func (e *envsKeys) Import(input []string, secretTag string) []Key {
	var ks []Key
	envRegexp := `^(\w*)=(.*)$`
	reEnv := regexp.MustCompile(envRegexp)
	reSecret := getSecretRegexp(secretTag)
	for _, env := range input {
		kv := reEnv.FindStringSubmatch(env)
		ks = append(ks, Key{
			Key:     kv[1],
			Value:   kv[2],
			Secrets: reSecret.parseValueToSecrets(kv[2]),
		})
	}
	return ks
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

func (f *flagsKeys) Import(input []string, secretTag string) []Key {
	var ks []Key
	flagRegexp := `^\-(\w*)=(.*)$`
	reFlag := regexp.MustCompile(flagRegexp)
	reSecret := getSecretRegexp(secretTag)
	for _, flag := range input {
		kv := reFlag.FindStringSubmatch(flag)
		ks = append(ks, Key{
			Key:     kv[1],
			Value:   kv[2],
			Secrets: reSecret.parseValueToSecrets(kv[2]),
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
