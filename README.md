# macgyver
A tool of decrypt and encrypt in Google Cloud Platform, which using key management. That tool friendly using golang's flags.

### Installation

```
go get github.com/17media/macgyver
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
-db_URL=<secret_prefix>MTAuMTAuMTAuMTA= -db_user=<secret_prefix>cm9vdA== -db_password=<secret_prefix>cGFzc3dvcmQ=
```

Decrypt
```
macgyver decrypt \
          --cryptoProvider=base64 \
          --flags="-db_URL=<secret_prefix>MTAuMTAuMTAuMTA= -db_user=<secret_prefix>cm9vdA== -db_password=<secret_prefix>cGFzc3dvcmQ="
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
-db_URL=<secret_prefix>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU= -db_user=<secret_prefix>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg== -db_password=<secret_prefix>CiQAfxfF5dBTxNZuLubqzLbilN0pzavOV7gyq7ZZHiH2oAEKm3MSMQCPVWdQhmTYSQwjIk4Xk5sgROOm4ExM0NacutDa7C2Ldp5qovv3uCJD4It/KHf5DUs=
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
          --flags="-db_URL=<secret_prefix>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU= -db_user=<secret_prefix>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg== -db_password=<secret_prefix>CiQAfxfF5dBTxNZuLubqzLbilN0pzavOV7gyq7ZZHiH2oAEKm3MSMQCPVWdQhmTYSQwjIk4Xk5sgROOm4ExM0NacutDa7C2Ldp5qovv3uCJD4It/KHf5DUs="
```
Output
```
-db_URL=10.10.10.10 -db_user=root -db_password=password
```

#### Using base64 with environment variables

Decrypt
```
# time ./macgyver decrypt \
export db_URL="<secret_prefix>MTAuMTAuMTAuMTA="
export db_user="<secret_prefix>cm9vdA=="
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
export db_URL="<secret_prefix>CiQAfxfF5QJgZYEvFhWwtv/x4Fou2R/8EqLheUDV+cdod3pS0rASNACPVWdQ+uFI6GtGWICaqA1xgfTVnBE+Gp4F1BkAohhdIPjQvnx+kqUPxebOiK1GKKmkMoU=
export db_user="<secret_prefix>CiQAfxfF5WuD0AfFN882MOtICNNNZ4Pj/QYERYiL/brcLcTRV9ISLQCPVWdQ8S1KZwNaZc6dIAXdoe8MIi26TcG1y5oeAqsxNxUp1Uxtz8mf1+8jvg=="
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


### Todo
- [x] Support environment variable
- [ ] refine cryptoProvide
  - [x] Base64
  - [ ] AWS
- [x] customize prefix (ex. `<kms>`)
- [x] go vendor version control
