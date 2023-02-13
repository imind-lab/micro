/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package api

import (
    "github.com/imind-lab/micro/v2/microctl/template"
)

// 生成conf/conf.yaml
func CreateConfCrt(data *template.Data) error {
    var tpl = `-----BEGIN CERTIFICATE-----
MIIDqjCCApKgAwIBAgIJAJEfwWYkH/bjMA0GCSqGSIb3DQEBCwUAMIGFMQswCQYD
VQQGEwJDTjELMAkGA1UECBMCQkoxDDAKBgNVBAcTA1NKUzETMBEGA1UEChMKaW1p
bmQudGVjaDENMAsGA1UECxMEdGVjaDEVMBMGA1UEAxQMKi5pbWluZC50ZWNoMSAw
HgYJKoZIhvcNAQkBFhFzb25nbGlAaW1pbmQudGVjaDAeFw0yMTA5MzAwOTE0NDNa
Fw0zMTA5MjgwOTE0NDNaMIGFMQswCQYDVQQGEwJDTjELMAkGA1UECBMCQkoxDDAK
BgNVBAcTA1NKUzETMBEGA1UEChMKaW1pbmQudGVjaDENMAsGA1UECxMEdGVjaDEV
MBMGA1UEAxQMKi5pbWluZC50ZWNoMSAwHgYJKoZIhvcNAQkBFhFzb25nbGlAaW1p
bmQudGVjaDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKFkYO8vtFxJ
KPN/kRmdWG9EYZo/p1qOWNC64AnSu3Y8OxoWWIUFPQHMvzvIUrNCo0aH+L4N3fAU
ujcXzYyU3UEaxHd3xQ76VjqtpXkUv/mHXQasMOxmnpulQ3CeTsbXM24i8owApGN4
KTC9clB50f0D1VZf3y41U+lUsSKqYOAwu6eqLMXnATGtY7p6y2EphHrTUzKLfHzo
UjsgyPEDSCTU5Ar7zgpLu3xx+Az+f5sMJjxnaZprk7suXT2/kv/QxKVqsCAU1sbh
mQETojhTqGe4u/uHwDvgtoZv+UYoDbwxw8dM8CSmV4EUlq0UzpaMtcMmXQoRs8qX
b2tC+Ci020MCAwEAAaMbMBkwFwYDVR0RBBAwDoIMKi5pbWluZC50ZWNoMA0GCSqG
SIb3DQEBCwUAA4IBAQAS0wVvtQV2yAcfigOOI1uSnTieOIIMVbkZlmjrT/gePgwB
Txkgd9exFFOppEjIKIQ4N40s23ZDp/iGVrxILpu4DWTVIZM8LrA767vPbDqx65ym
ZIZfoKuFmmSEs5+38mQuvOxugIm7h9ZJbq1MyCyXvY9Gs73esq1f0wCQP8c4YAgz
s3HQHAnvNL8OJCxncdM7LUtz+B9jMpaMyp/JYy4vUNSYZ3Eb10caaCJsr5rQNqxv
mhKwZrOEDMKVvheEfHXHQ2tsDWWwr4vr4f8UolM/M17VGtWdg2hqs/g9Ohze2d2X
SbjNQfx9MfJPBSWL8ZxVpOGP/swfwCUewheKaGXm
-----END CERTIFICATE-----
`

    path := "./" + data.Domain + "/" + data.Repo + "/" + data.Service + "-api/conf/ssl/"
    name := "tls.crt"

    return template.CreateFile(data, tpl, path, name)
}
