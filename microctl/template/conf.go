/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package template

import (
	"os"
	"text/template"
)

// 生成conf
func CreateConf(data *Data) error {
	// 生成conf.yaml
	var tpl = `service:
  name: {{.Service}}
  port: #监听端口
    http: 80
    grpc: 50051
  profile:
    rate: 1

db:
  logLevel: 4
  max:
    open: 100
    idle: 5
    life: 30
  imind:
    tablePrefix: tbl
    master:
      host: 127.0.0.1
      port: 3306
      user: root
      pass: 123456
      name: mind
    replica:
      host: 127.0.0.1
      port: 3306
      user: root
      pass: 123456
      name: mind

redis:
  addr: '127.0.0.1:6379'
  db: 0

kafka:
  business:
    producer:
      - '127.0.0.1:9092'
    consumer:
      - '127.0.0.1:9092'
    topic:
      {{.Service}}Create: {{.Service}}_create
      {{.Service}}Update: {{.Service}}_update

tracing:
  agent: '172.16.50.50:6831'
  type: const
  param: 1
  name:
    client: imind-{{.Service}}-cli
    server: imind-{{.Service}}-srv

log:
  path: './logs/ms.log'
  level: -1
  age: 7
  size: 128
  backup: 30
  compress: true
  format: json
`
	t, err := template.New("conf").Parse(tpl)
	if err != nil {
		return err
	}

	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/conf/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "conf.yaml"

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成tls.crt
	tpl = `-----BEGIN CERTIFICATE-----
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
	t, err = template.New("tlscrt").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/conf/ssl/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "tls.crt"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成tls.key
	tpl = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQChZGDvL7RcSSjz
f5EZnVhvRGGaP6dajljQuuAJ0rt2PDsaFliFBT0BzL87yFKzQqNGh/i+Dd3wFLo3
F82MlN1BGsR3d8UO+lY6raV5FL/5h10GrDDsZp6bpUNwnk7G1zNuIvKMAKRjeCkw
vXJQedH9A9VWX98uNVPpVLEiqmDgMLunqizF5wExrWO6esthKYR601Myi3x86FI7
IMjxA0gk1OQK+84KS7t8cfgM/n+bDCY8Z2maa5O7Ll09v5L/0MSlarAgFNbG4ZkB
E6I4U6hnuLv7h8A74LaGb/lGKA28McPHTPAkpleBFJatFM6WjLXDJl0KEbPKl29r
QvgotNtDAgMBAAECggEAJP5LMcgvgU/LsTS2X7avRKHZ9W9NvvYN6ZpMLDQ/f/SC
X5Jrq+Htf/Ici2l5e1l073/PRlELZDJ8FJDCLs0YggnOqsurQamkBkMzQLO+5UVf
z128wRMsm+Sftrqyt+UwBri/+9NR2vL6Dg/+me+ycFpmlivXjlHu7/bXu2huWDS2
Nx0RB/jwO/4PcFr2RwGVgvNt/9oMUegZP9Iqzp9kBxwafNL/7ASSYDNd7VlhnrSG
fB8A8tvOLFCjnmuXEMOz+znc+0Y8JFmbxd122Qd07L406+tgAz4CAgGpP4JiRMmj
RhORh3dCzpXabsTyA5ybZdJMG/X26EB/istVyHhwUQKBgQDQ1xjkjpoDnWFz45mm
euYT9plCdDaGG2ktR3oPw7lvBTCy4D38TZNfi0W98ZBStIh8boNtqqBVk/6I6ini
7RtwKHJ2WNRGBxC+rPIS7WL0MEsxX+R9MU0rY80zwFScvHLxsiWKA76T9sX4eV4O
lGmHEA7wpiv0mU2+9E04XIOJxQKBgQDF1lb/Qox7+oHvQLxKDQTAXVsTcsSvSKSM
U/JX9S8U5bKMSj5ERiLCuWqzFz/515VVD+J3WtY774PP1/GejFazs4FxE4KPi+CH
cmj8fMPRJl3v3g+l6BVp8BiKz7DZwFFgwqhpPx+woGp7FebLUIS+7thEpQeK8Gc+
OrhAxbaJZwKBgQCjhto6FaNpkzF84kos/uzr0tudGoybJBmOV/qvH24zDZhdaJOA
3Wm5lb+NhPqimDSLYqnNFJ8pg5H6uYkE5O7oOvOt0c6d6uhktd1zjqg+VxZ52gF7
OkCX6jUDAeX/ONy0fu9AC8CN8dyAvOA2gGXFWYCpVST0CZrEHF3e9SoWlQKBgQCQ
45pRVe7HOb8Bdxqu7Pvm2jhCdRJBAWWpdC2PZ3y0xEjQX+tcWzVIAT14rfVnyBCQ
/JIyMW+m85JInPFS2ZsB/tw08UH0WU/2Qr9K8yECQyQW8T3qlp9gN7vxpYvy3dt7
jvSCJ/3QgJubS338txqRLyFqnKZ6hfhG5gBdR6+YzwKBgEXOInVrRZW5sEyVSgFs
1h63Y5BEbyGuoSpV/tsu6T49S2OQYfJysq8yxL+du9TBzgeSTpxJjgt1mVJjVheU
yZ7kDxFkoAqp1NsmK+RCKQ70TIxuueby7ZiMNuR6ZenKutRHOL1O8GSAvz6bToo/
VRHj9cnoQrOf7n2ABB5uM89O
-----END PRIVATE KEY-----
`
	t, err = template.New("tlskey").Parse(tpl)
	if err != nil {
		return err
	}

	fileName = dir + "tls.key"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
