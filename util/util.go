/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package util

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"

	"golang.org/x/text/cases"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/language"
	"google.golang.org/grpc/metadata"
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

func MapFillStruct(data map[string]interface{}, result interface{}) {
	t := reflect.ValueOf(result).Elem()
	for k, v := range data {
		key := cases.Title(language.Und, cases.NoLower).String(k)
		val := t.FieldByName(key)
		if val.IsValid() {
			val.Set(reflect.ValueOf(v))
		}
	}
}

func CheckFileType(ext string, exts []string) bool {
	for _, v := range exts {
		if v == ext {
			return false
		}
	}
	return true
}

func DecodeInt32(code byte) int32 {
	num := []byte{code}
	return int32(binary.BigEndian.Uint32(num))
}

func Md5Hash(s string) string {
	signByte := []byte(s)
	hash := md5.New()
	hash.Write(signByte)
	return hex.EncodeToString(hash.Sum(nil))
}

// Base64Encode base64 加密
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

//Base64Decode  base64 解密
func Base64Decode(s string) string {
	var b []byte
	var err error
	x := len(s) * 3 % 4
	switch {
	case x == 2:
		s += "=="
	case x == 1:
		s += "="
	}
	if b, err = base64.StdEncoding.DecodeString(s); err != nil {
		return string(b)
	}
	return string(b)
}

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func MysqlEscape(src string) string {
	if len(src) == 0 {
		return ""
	}
	i := 0
	dest := make([]rune, len(src)*2)
	for _, c := range src {
		switch c {
		case '\r', '\n', '\\', '\'', '"':
			dest[i] = '\\'
			dest[i+1] = c
			i = i + 2
		case '\032':
			dest[i] = '\\'
			dest[i+1] = 'Z'
			i = i + 2
		default:
			dest[i] = c
			i = i + 1

		}
	}
	return string(dest[:i])
}

func StripTags(content string) string {
	re := regexp.MustCompile(`<(.|\n)*?>`)
	return re.ReplaceAllString(content, "")
}

//将utf-8编码的字符串转换为GBK编码
func Convert2GBK(str string) string {
	if ret, err := simplifiedchinese.GBK.NewEncoder().String(str); err == nil {
		return ret
	}
	return ""
}

//将GBK编码的字符串转换为utf-8编码
func Convert2UTF8(str string) string {
	if ret, err := simplifiedchinese.GBK.NewDecoder().String(str); err == nil {
		return ret
	}
	return ""
}

func InSliceInt(item int, items ...int) bool {
	for _, i := range items {
		if item == i {
			return true
		}
	}
	return false
}

func InSliceInt32(item int32, items ...int32) bool {
	for _, i := range items {
		if item == i {
			return true
		}
	}
	return false
}

func InSliceString(item string, items ...string) bool {
	for _, i := range items {
		if item == i {
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

	body, err := ioutil.ReadAll(rsp.Body)
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

	body, err := ioutil.ReadAll(rsp.Body)
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

func MapMerge(to map[string]interface{}, from ...map[string]interface{}) map[string]interface{} {
	for _, m := range from {
		for k, v := range m {
			to[k] = v
		}
	}
	return to
}

func JoinIntSlice(ids []int) string {
	var buf bytes.Buffer

	for _, id := range ids {
		if buf.Len() > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.Itoa(id))

	}
	return buf.String()
}

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
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

func GetString(req *http.Request, key string, def string) string {
	val := req.PostFormValue(key)
	if len(val) == 0 {
		val = req.PostFormValue(strings.ToUpper(key))
		if len(val) == 0 {
			val = req.FormValue(key)
			if len(val) == 0 {
				val = req.FormValue(strings.ToUpper(key))
				if len(val) == 0 {
					return def
				}
			}
		}
	}
	return StripTags(val)
}

func GetInt(req *http.Request, key string, def int) int {
	val := req.PostFormValue(key)
	if len(val) == 0 {
		val = req.PostFormValue(strings.ToUpper(key))
		if len(val) == 0 {
			val = req.FormValue(key)
			if len(val) == 0 {
				val = req.FormValue(strings.ToUpper(key))
				if len(val) == 0 {
					return def
				}
			}
		}
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}

func GetInt32(req *http.Request, key string, def int) int32 {
	i := GetInt(req, key, def)
	return int32(i)
}

func GetValues(req *http.Request, key string) []string {
	if values, ok := req.PostForm[key]; ok {
		return values
	}
	if values, ok := req.Form[key]; ok {
		return values
	}
	return []string{}
}

func GetInt32Values(req *http.Request, key string) []int32 {
	values := GetValues(req, key)
	if len(values) > 0 {
		return MapStringToInt32(values, StringToInt32)
	}
	return []int32{}
}

func MapStringToInt32(value []string, fn func(string) int32) []int32 {
	ret := make([]int32, 0, len(value))
	for _, val := range value {
		ret = append(ret, fn(val))
	}
	return ret
}

func StringToInt32(val string) int32 {
	ret, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return int32(ret)
}

// 过滤 emoji 表情
func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
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

func DurationFormat(duration int) string {
	var buf bytes.Buffer
	hour := 0
	if duration > 3600 {
		hour = duration / 3600
		buf.WriteString(fmt.Sprintf("%02d:", hour))
	}
	buf.WriteString(fmt.Sprintf("%02d:", (duration-3600*hour)/60))
	buf.WriteString(fmt.Sprintf("%02d", duration%60))
	return buf.String()
}

func DurationNewFormat(duration int) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%02d:", duration/60))
	buf.WriteString(fmt.Sprintf("%02d", duration%60))
	return buf.String()
}

func ReplaceAll(special string, str ...string) string {
	for _, v := range str {
		special = strings.ReplaceAll(special, v, "")
	}
	return special
}

func GetMetaInt(md metadata.MD, key string, def int) int {
	vals := md.Get(key)
	if len(vals) > 0 {
		val, err := strconv.Atoi(vals[0])
		if err == nil {
			return val
		}
	}
	return def
}

func GetMetaString(md metadata.MD, key string, def string) string {
	vals := md.Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return def
}

func GetMetaFloat(md metadata.MD, key string, def float32) float32 {
	vals := md.Get(key)
	if len(vals) > 0 {
		val, err := strconv.ParseFloat(vals[0], 32)
		if err == nil {
			return float32(val)
		}
	}
	return def
}

func Contains(list []interface{}, value interface{}) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func ContainsString(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

//获取下一个可用消息id
func GetId(ids []int) int {
	if len(ids) == 0 {
		return 1
	} else {
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
}

func RandDuration(max int) time.Duration {
	rand.Seed(time.Now().UnixNano())
	rd := rand.Intn(max)
	return time.Duration(rd) * time.Second
}

func MapConvertValue(data map[string]string) map[string]interface{} {
	ret := make(map[string]interface{})
	for key, val := range data {
		ret[key] = val
	}
	return ret
}

func AppendString(keys ...string) string {
	var result []byte
	for _, key := range keys {
		result = append(result, key...)
	}
	return string(result)
}

func GetPtrFuncName() (string, string) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])

	info := strings.Split(f.Name(), ".")
	cnt := len(info)

	layer := info[cnt-2]
	layer = strings.Replace(layer, "(*", "", 1)
	layer = strings.Replace(layer, ")", "", 1)

	return layer, info[cnt-1]
}

func GetFuncName() (string, string) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])

	info := strings.Split(f.Name(), ".")
	cnt := len(info)

	return info[cnt-2], info[cnt-1]
}
