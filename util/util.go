/**
 * ImindLab
 *
 * Create by SongLi on 2023/12/28
 * Copyright © 2023 imind.tech All rights reserved.
 */

package util

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"sort"
	"strings"
	"time"
	"unsafe"
)

func GzipResponse(w http.ResponseWriter, buffer []byte) {

	w.Header().Add("Accept-Charset", "utf-8")
	w.Header().Add("Content-Type", "application/x-protobuf")
	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	defer gz.Close()
	gz.Write(buffer)

	gz.Flush()
}

func Md5Hash(s string) string {
	signByte := []byte(s)
	hash := md5.New()
	hash.Write(signByte)
	return hex.EncodeToString(hash.Sum(nil))
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func IsContains[K comparable](item K, items ...K) bool {
	for _, val := range items {
		if item == val {
			return true
		}
	}
	return false
}

func GenSignature(params map[string]string, secretKey string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	buf := bytes.Buffer{}
	for _, key := range keys {
		buf.WriteString(key)
		buf.WriteString(params[key])
	}
	buf.WriteString(secretKey)
	return Md5Hash(buf.String())
}

func EncodeParams(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	buf := bytes.Buffer{}
	for _, key := range keys {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(params[key])
	}
	return buf.String()
}

func HttpDo(ctx context.Context, params map[string]string, url string) (string, error) {
	query := EncodeParams(params)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(query))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func PostJSON(ctx context.Context, jsonBody []byte, url string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	ret := map[string]interface{}{}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func AddCookie(rw http.ResponseWriter, req *http.Request, name, value string, expire time.Duration) {
	cookie, err := req.Cookie(name)
	if err != nil {
		cookie = &http.Cookie{}
		cookie.Name = name
	}
	SetCookie(rw, cookie, value, expire)
}

func SetCookie(rw http.ResponseWriter, cookie *http.Cookie, value string, expire time.Duration) {
	cookie.Value = value
	cookie.Expires = time.Now().Add(expire)
	http.SetCookie(rw, cookie)
}

func MapMerge[K comparable, V any](to map[K]V, from ...map[K]V) map[K]V {
	for _, m := range from {
		for k, v := range m {
			to[k] = v
		}
	}
	return to
}

func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = strings.TrimSpace(strings.Split(ip, ",")[0])
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

// GetId gets the next available id
func GetId(ids []int) int {
	if len(ids) == 0 {
		return 1
	}
	max := int(^uint16(0))
	sort.Ints(ids)
	id := ids[len(ids)-1] + 1
	for {
		if id > max {
			id = 1
		}
		available := true
		for _, v := range ids {
			if v == id {
				available = false
				break
			}
		}
		if available {
			return id
		}
		id++
	}
}

func RandDuration(max int64) time.Duration {
	var rand = rand.New(rand.NewSource(max))
	rand.Seed(time.Now().UnixNano())
	rd := rand.Int63n(max)
	return time.Duration(rd) * time.Second
}

func AppendString(keys ...string) string {
	var ret bytes.Buffer
	for i, v := range keys {
		if i > 0 {
			ret.WriteByte('_')
		}
		ret.WriteString(v)
	}
	return ret.String()
}

func GetPascalCase(name string) string {
	return strings.ReplaceAll(cases.Title(language.English).String(name), "-", "")
}

func GetCamelCase(name string) string {
	components := strings.SplitN(name, "-", 2)
	if len(components) == 2 {
		var buffer bytes.Buffer
		buffer.WriteString(components[0])
		buffer.WriteString(strings.ReplaceAll(cases.Title(language.English).String(components[1]), "-", ""))
		return buffer.String()
	}
	return components[0]
}

func GetPtrFuncName() (string, string) {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])

	info := strings.Split(f.Name(), ".")
	cnt := len(info)

	layer := info[cnt-2]
	layer = strings.Replace(layer, "(*", "", 1)
	layer = strings.Replace(layer, ")", "", 1)

	return layer, info[cnt-1]
}

func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])

	info := strings.Split(f.Name(), ".")
	cnt := len(info)

	return strings.Join(info[cnt-2:], ".")
}
