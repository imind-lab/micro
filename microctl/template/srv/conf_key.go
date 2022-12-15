/**
 *  MindLab
 *
 *  Create by songli on 2022/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package srv

import (
	"github.com/imind-lab/micro/microctl/template"
)

// 生成conf/conf.yaml
func CreateConfKey(data *template.Data) error {
	var tpl = `-----BEGIN PRIVATE KEY-----
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

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/conf/ssl/"

	name := "tls.key"

	return template.CreateFile(data, tpl, path, name)
}
