# macgyver
A tool of decrypt and encrypt in Google Cloud Platform, which using key management. That tool friendly using golang's flags.

### Installation

```
go get -u github.com/17media/macgyver
```

### Usage
```
macgyver help
```

### Command Use

#### Using base64 with text

*Encrypt*
```
macgyver encrypt \
          --cryptoProvider=base64 \
          --flags="-db_URL=10.10.10.10 -db_user=root -db_password=password"
```

Output

```
-db_URL=<SECRET_TAG>MTAuMTAuMTAuMTA=</SECRET_TAG> -db_user=<SECRET_TAG>cm9vdA==</SECRET_TAG> -db_password=<SECRET_TAG>cGFzc3dvcmQ=</SECRET_TAG>
```

Decrypt
```
macgyver decrypt \
          --cryptoProvider=base64 \
          --flags="-db_URL=<SECRET_TAG>MTAuMTAuMTAuMTA=</SECRET_TAG> -db_user=<SECRET_TAG>cm9vdA==</SECRET_TAG> -db_password=<SECRET_TAG>cGFzc3dvcmQ=</SECRET_TAG>"
```

Output

```
-db_URL=10.10.10.10 -db_user=root -db_password=password
```

#### Using GCP KMS and service account JSON key by Google with text

Encrypt

```
macgyver encrypt \
          --cryptoProvider=gcp \
          --oAuthLocation=<oAuthLocation>.json \
          --GCPprojectID="<ProjectID>" \
          --GCPlocationID="<LocationID>" \
          --GCPkeyRingID="<KeyRingID>" \
          --GCPcryptoKeyID="<cryptoKeyID>" \
          --flags="-db_URL=10.10.10.10 -db_user=root -db_password=password"
```
Output
```
-db_URL=<SECRET_TAG>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=</SECRET_TAG> -db_user=<SECRET_TAG>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg==</SECRET_TAG> -db_password=<SECRET_TAG>CiQAfxfF5dBTxNZuLubqzLbilN0pzavOV7gyq7ZZHiH2oAEKm3MSMQCPVWdQhmTYSQwjIk4Xk5sgROOm4ExM0NacutDa7C2Ldp5qovv3uCJD4It/KHf5DUs=</SECRET_TAG>
```

Decrypt

```
macgyver decrypt \
          --cryptoProvider=gcp \
          --oAuthLocation=<oAuthLocation>.json \
          --GCPprojectID="<ProjectID>" \
          --GCPlocationID="<LocationID>" \
          --GCPkeyRingID="<KeyRingID>" \
          --GCPcryptoKeyID="<cryptoKeyID>" \
          --flags="-db_URL=<SECRET_TAG>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=</SECRET_TAG> -db_user=<SECRET_TAG>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg==</SECRET_TAG> -db_password=<SECRET_TAG>CiQAfxfF5dBTxNZuLubqzLbilN0pzavOV7gyq7ZZHiH2oAEKm3MSMQCPVWdQhmTYSQwjIk4Xk5sgROOm4ExM0NacutDa7C2Ldp5qovv3uCJD4It/KHf5DUs=</SECRET_TAG>"
```
Output
```
-db_URL=10.10.10.10 -db_user=root -db_password=password
```

#### Using base64 with environment variables

Decrypt
```
# time ./macgyver decrypt \
export db_URL="<SECRET_TAG>MTAuMTAuMTAuMTA=</SECRET_TAG>"
export db_user="<SECRET_TAG>cm9vdA==</SECRET_TAG>"
export db_password="password"

eval $(macgyver decrypt \
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

#### Using GCP KMS and service account JSON key by Google with environment variables


Decrypt

```
export db_URL="<SECRET_TAG>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=</SECRET_TAG>"
export db_user="<SECRET_TAG>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg==</SECRET_TAG>"
export db_password="password"

eval $(macgyver decrypt \
          --cryptoProvider=gcp \
          --keysType=env \
          --oAuthLocation=<oAuthLocation>.json \
          --GCPprojectID="<ProjectID>" \
          --GCPlocationID="<LocationID>" \
          --GCPkeyRingID="<KeyRingID>" \
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

### Behavior when mavgyver seeing the `secret tag`
- Lower case of `secret tag` to represent a plaintext and upper case of `secret tag` to represent a ciphertext
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


### Todo
- [x] Support environment variable
- [ ] refine cryptoProvide
  - [x] Base64
  - [ ] AWS
- [x] customize secret tag (ex. `<kms>XXX</kms>`)
- [x] go vendor version control
- [ ] add unit test
