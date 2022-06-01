# Macgyver

[![CircleCI](https://circleci.com/gh/17media/macgyver/tree/master.svg?style=svg)](https://circleci.com/gh/17media/macgyver/tree/master)

A tool for decrypting and encrypting strings in GCP / AWS by using key management. The tool is golang's flags friendly.

## Installation

```
$ go get -u github.com/17media/macgyver
```

## Usage
```
$ macgyver help
A tool for decrypting and encrypting strings in GCP / AWS by using KMS,
The tool is golang's flags friendly.
For example:
$ go run main.go decrypt                 \
                --cryptoProvider=gcp     \
                --GCPprojectID="demo"    \
                --GCPlocationID="global" \
                --GCPkeyRingID="foo"     \
                --GCPcryptoKeyID="bar"   \
                --flags="-a=kms_asda

Usage:
  macgyver [command]

Available Commands:
  decrypt     Decrypt entire flags
  encrypt     Encrypt entire flags
  help        Help about any command
  version     Print the version number of Macgyver

Flags:
      --AWScryptoKeyID string   the cryptoKeyID of AWS
      --AWSlocationID string    the locationID of AWS
      --AWSprofileName string   the profile name used for AWS authentication
      --GCPcryptoKeyID string   the cryptoKeyID of GCP
      --GCPkeyRingID string     the keyRingID of GCP
      --GCPlocationID string    the locationID of GCP
      --GCPprojectID string     the projectID of GCP
      --config string           config file (default is $HOME/.macgyver.yaml)
      --cryptoProvider string   Which type you using encrypto and encryto
      --file string             absolute filepath for a yaml you want to decrypt/encrypt
      --flags string            the sort code of the contact account
  -h, --help                    help for macgyver
      --keysType string         Which input type you using for encrypto and encryto (e.g. text, file and env) (default "text")
      --oAuthLocation string    location of the JSON key credentials file. If empty then use the Google Application Defaults.
      --secretTag string        the prefix of secret (default "secret_tag")

Use "macgyver [command] --help" for more information about a command.
```

---

## Example

---

### Using Base64 with text
*CAUTION:* `Base64` is only for testing, *DO NOT* use it in *production* environment.

#### Encrypt
```
$ macgyver encrypt                \
          --cryptoProvider=base64 \
          --flags="-db_URL=10.10.10.10 -db_user=root -db_password=password"
```

Output
```
-db_URL=<SECRET_TAG>MTAuMTAuMTAuMTA=</SECRET_TAG> -db_user=<SECRET_TAG>cm9vdA==</SECRET_TAG> -db_password=<SECRET_TAG>cGFzc3dvcmQ=</SECRET_TAG>
```

#### Decrypt
```
$ macgyver decrypt                \
          --cryptoProvider=base64 \
          --flags="-db_URL=<SECRET_TAG>MTAuMTAuMTAuMTA=</SECRET_TAG> -db_user=<SECRET_TAG>cm9vdA==</SECRET_TAG> -db_password=<SECRET_TAG>cGFzc3dvcmQ=</SECRET_TAG>"
```

Output
```
-db_URL=10.10.10.10 -db_user=root -db_password=password
```

---

### Using GCP KMS and service account JSON key by Google with text

#### Decrypt && Encrypt
```zsh

PROJECT=media17-stag

cipher_text=$(macgyver decrypt          \
          --cryptoProvider=gcp          \
          --GCPprojectID=${PROJECT}     \
          --GCPlocationID=global        \
          --GCPkeyRingID=app            \
          --GCPcryptoKeyID=flags        \
          --flags="-key=${1}")

echo '#######################'
echo 'cipher text >>> '$cipher_text
arr=("${(@s/=/)cipher_text}")
cipher=$arr[2]
echo 'cipher >>> '$cipher

PROVIDER=gcp
PROJECT=media17-uat

echo '==============='
echo $PROJECT
echo '==============='
macgyver encrypt                               \
	  --cryptoProvider="${PROVIDER}"             \
	  --keysType=text                            \
          --GCPprojectID="${PROJECT}"          \
          --GCPlocationID=global               \
          --GCPkeyRingID=app                   \
          --GCPcryptoKeyID=flags               \
          --flags="-pwd=${cipher}" 

# USAGE
# zsh decryp.zsh "<SECRET_TAG>cipher_text</SECRET_TAG>"
# Note that, you should change PROJECT name based on your demand
```

Output
```
#######################
cipher text >>> -key=ciphertext
===============
plaintext
===============
-pwd=<SECRET_TAG>new_ciphertext</SECRET_TAG>
```


#### Encrypt
```
$ macgyver encrypt                             \
          --cryptoProvider=gcp                 \
          --oAuthLocation=<oAuthLocation>.json \
          --GCPprojectID="<ProjectID>"         \
          --GCPlocationID="<LocationID>"       \
          --GCPkeyRingID="<KeyRingID>"         \
          --GCPcryptoKeyID="<cryptoKeyID>"     \
          --flags="-db_URL=10.10.10.10 -db_user=root -db_password=password"
```

Output
```
-db_URL=<SECRET_TAG>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=</SECRET_TAG> -db_user=<SECRET_TAG>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg==</SECRET_TAG> -db_password=<SECRET_TAG>CiQAfxfF5dBTxNZuLubqzLbilN0pzavOV7gyq7ZZHiH2oAEKm3MSMQCPVWdQhmTYSQwjIk4Xk5sgROOm4ExM0NacutDa7C2Ldp5qovv3uCJD4It/KHf5DUs=</SECRET_TAG>
```

#### Decrypt
```
$ macgyver decrypt                             \
          --cryptoProvider=gcp                 \
          --oAuthLocation=<oAuthLocation>.json \
          --GCPprojectID="<ProjectID>"         \
          --GCPlocationID="<LocationID>"       \
          --GCPkeyRingID="<KeyRingID>"         \
          --GCPcryptoKeyID="<cryptoKeyID>"     \
          --flags="-db_URL=<SECRET_TAG>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=</SECRET_TAG> -db_user=<SECRET_TAG>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg==</SECRET_TAG> -db_password=<SECRET_TAG>CiQAfxfF5dBTxNZuLubqzLbilN0pzavOV7gyq7ZZHiH2oAEKm3MSMQCPVWdQhmTYSQwjIk4Xk5sgROOm4ExM0NacutDa7C2Ldp5qovv3uCJD4It/KHf5DUs=</SECRET_TAG>"
```

Output
```
-db_URL=10.10.10.10 -db_user=root -db_password=password
```

---

### Using AWS KMS with text

#### Encrypt
Using ENVs for AWS authentication
```
# Export your account credentials to access the AWS KMS service
export AWS_ACCESS_KEY_ID='<aws_access_key_id>'
export AWS_SECRET_ACCESS_KEY='<aws_secret_access_key>'

$ macgyver encrypt                      \
          --cryptoProvider="aws"        \
          --AWSlocationID="<LocatioID>" \
          --AWScryptoKeyID="<KeyID>"    \
          --flags="-db_URL=10.10.10.10 -db_user=root -db_password=password"
```

Using AWS profile configured in ~/.aws/config
```
$cat ~/.aws/config

[profile <ProfileName>]
region = us-west-2
role_arn = arn:aws:iam:::role/*
source_profile = <source profile name configured in ~/.aws/credentials>

$ macgyver encrypt                         \
          --cryptoProvider="aws"           \
          --AWSprofileName="<ProfileName>" \
          --AWSlocationID="<LocatioID>"    \
          --AWScryptoKeyID="<KeyID>"       \
          --flags="-db_URL=10.10.10.10 -db_user=root -db_password=password"
```

Output
```
-db_URL=<SECRET_TAG>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=</SECRET_TAG> -db_user=<SECRET_TAG>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg==</SECRET_TAG> -db_password=<SECRET_TAG>CiQAfxfF5dBTxNZuLubqzLbilN0pzavOV7gyq7ZZHiH2oAEKm3MSMQCPVWdQhmTYSQwjIk4Xk5sgROOm4ExM0NacutDa7C2Ldp5qovv3uCJD4It/KHf5DUs=</SECRET_TAG>
```

#### Decrypt
Using ENVs for AWS authentication
```
# Export your account credentials to access the AWS KMS service
export AWS_ACCESS_KEY_ID='<aws_access_key_id>'
export AWS_SECRET_ACCESS_KEY='<aws_secret_access_key>'

$ macgyver decrypt                      \
          --cryptoProvider="aws"        \
          --AWSlocationID="<LocatioID>" \
          --AWScryptoKeyID="<KeyID>"    \
          --flags="-db_URL=<SECRET_TAG>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=</SECRET_TAG> -db_user=<SECRET_TAG>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg==</SECRET_TAG> -db_password=<SECRET_TAG>CiQAfxfF5dBTxNZuLubqzLbilN0pzavOV7gyq7ZZHiH2oAEKm3MSMQCPVWdQhmTYSQwjIk4Xk5sgROOm4ExM0NacutDa7C2Ldp5qovv3uCJD4It/KHf5DUs=</SECRET_TAG>"
```

Using AWS profile configured in ~/.aws/config
```
$cat ~/.aws/config

[profile <ProfileName>]
region = us-west-2
role_arn = arn:aws:iam:::role/*
source_profile = <source profile name configured in ~/.aws/credentials>

$ macgyver decrypt                         \
          --cryptoProvider="aws"           \
          --AWSprofileName="<ProfileName>" \
          --AWSlocationID="<LocatioID>"    \
          --AWScryptoKeyID="<KeyID>"       \
          --flags="-db_URL=<SECRET_TAG>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=</SECRET_TAG> -db_user=<SECRET_TAG>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg==</SECRET_TAG> -db_password=<SECRET_TAG>CiQAfxfF5dBTxNZuLubqzLbilN0pzavOV7gyq7ZZHiH2oAEKm3MSMQCPVWdQhmTYSQwjIk4Xk5sgROOm4ExM0NacutDa7C2Ldp5qovv3uCJD4It/KHf5DUs=</SECRET_TAG>"
```

Output
```
-db_URL=10.10.10.10 -db_user=root -db_password=password
```


---

### Using Base64 with environment variables

#### Decrypt
```
# time ./macgyver decrypt \
export db_URL="<SECRET_TAG>MTAuMTAuMTAuMTA=</SECRET_TAG>"
export db_user="<SECRET_TAG>cm9vdA==</SECRET_TAG>"
export db_password="password"

eval $(macgyver decrypt                 \
                --cryptoProvider=base64 \
                --keysType=env)
echo $db_URL
echo $db_user
echo $db_password
```

Output
```
10.10.10.10
root
password
```

---

### Using GCP KMS and service account JSON key by Google with environment variables

#### Decrypt
```
export db_URL="<SECRET_TAG>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=</SECRET_TAG>"
export db_user="<SECRET_TAG>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg==</SECRET_TAG>"
export db_password="password"

eval $(macgyver decrypt                        \
          --cryptoProvider=gcp                 \
          --keysType=env                       \
          --oAuthLocation=<oAuthLocation>.json \
          --GCPprojectID="<ProjectID>"         \
          --GCPlocationID="<LocationID>"       \
          --GCPkeyRingID="<KeyRingID>"         \
          --GCPcryptoKeyID="<cryptoKeyID>")

echo $db_URL
echo $db_user
echo $db_password

```

Output
```
10.10.10.10
root
password
```

---

### Using AWS KMS with environment variables

#### Decrypt
```
export db_URL="<SECRET_TAG>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=</SECRET_TAG>"
export db_user="<SECRET_TAG>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg==</SECRET_TAG>"
export db_password="password"

eval $(macgyver decrypt                                           \
          --cryptoProvider="aws"                                  \
          --AWSlocationID="us-west-2"                             \
          --AWScryptoKeyID="56c43f67-517b-4cbd-91ae-6d9c064b5671" \

echo $db_URL
echo $db_user
echo $db_password

```

Output
```
10.10.10.10
root
password
```

---

---

### Using Base64 with file type (yaml)

Currently, file type only support yaml file. And because of some limitation and easy usages, macgyver won't keep anchor and reference in yaml after parsing.

It only supports encrypt string type item, so if you want to encrypt other types (e.g. int, float, bool) you need to quote them.


#### Encrypt and Decrypt
```

echo "test: abc
test2: 
    abc: cde
    def: 1
    ghi: 
        - abc
        - true" > test2.yaml

cat test2.yaml

>> 
test: abc
test2: 
    abc: cde
    def: 1
    ghi: 
        - abc
        - true



macgyver encrypt \
--cryptoProvider=base64 \
--keysType=file \
--file=test2.yaml

cat test2.yaml
>>
test: <SECRET_TAG>YWJj</SECRET_TAG>
test2:
    abc: <SECRET_TAG>Y2Rl</SECRET_TAG>
    def: 1
    ghi:
        - <SECRET_TAG>YWJj</SECRET_TAG>
        - true


macgyver decrypt \
--cryptoProvider=base64 \
--keysType=file \
--file=test2.yaml

>> 
test: abc
test2:
    abc: cde
    def: 1
    ghi:
        - abc
        - true


```

---

### Behavior when mavgyver seeing the `secret tag`
- Lower case of `secret tag` to represent a plaintext and upper case of `SECCRET TAG` to represent a ciphertext
- When encryption:
    1. Pure secret will be regarded the entire text as a secret to be encrypted. Example:
       `-test1=secret-flag` => `-test1=<TAG>ciphertext-1</TAG>`
    2. Tagged secret without other characters, the secret between the secret tag will be encrypted. Example:
       `-test2=<tag>secret-flag</tag>` => `-test1=<TAG>ciphertext-2</TAG>`
    3. Tagged secret with other characters, only the secret between the secret tag will be encrypted. Example:
       `-test3=NotSecretPrefix/<tag>secret-flag</tag>/NotSecretSuffix` => `-test3=NotSecretPrefix/<TAG>ciphertext-3</TAG>/NotSecretSuffix`
- When decryption:
    1. Text without secret tag will be not decrypted. Example:
       `-test4=not-secret-flag` => `-test4=not-secret-flag`
    2. All tagged secrets between the secret tag will be decrypted. Example:
        `-test1=<TAG>Ciphertext-1</TAG>` => `-test1=secret-flag`
        `-test2=<TAG>Ciphertext-2</TAG>` => `-test2=secret-flag`
        and `-test3=NotSecretPrefix/<TAG>ciphertext-3</TAG>/NotSecretSuffix` => `-test3=NotSecretPrefix/secret-flag/NotSecretSuffix`

---

### Todo
- [x] Support environment variable
- [ ] Support text file
- [x] Refine cryptoProvide
  - [x] Base64
  - [x] GCP KMS
  - [x] AWS KMS
- [x] Customize secret tag (ex. `<kms>XXX</kms>`)
- [x] go vendor version control
- [ ] Add unit test
