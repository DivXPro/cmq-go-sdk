package cmq

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// MakeSignPlainText
func MakeSignPlainText(
	params interface{},
	method string,
	host string,
	path string,
) string {
	url := host + path
	paramStr := MakeParamStr(params)
	reg := regexp.MustCompile(`(?i:^https)|(?i:^http)://`)
	// TODO: 用正则方式替换https://,http://
	return method + reg.ReplaceAllString(url, "") + paramStr
}

// MakeParamStr 构建提交参数的字符串
func MakeParamStr(params interface{}) string {
	var paramStr string
	v := reflect.ValueOf(params)
	elemType := v.Elem().Type()
	count := v.NumField()
	var keyNames []string
	for i := 0; i < count; i++ {
		if elemType.Field(i).Name == "Signature" {
			continue
		}
		keyName := strings.Replace(elemType.Field(i).Name, "_", ".", -1)
		keyNames = append(keyNames, keyName)
	}
	// 对键值进行排序
	sort.Strings(keyNames)
	for n := 0; n < len(keyNames); n++ {
		if n == 0 {
			paramStr = paramStr + "?"
		} else {
			paramStr = paramStr + "&"
		}
		switch v.Field(n).Kind() {
		case reflect.String:
			paramStr = paramStr + keyNames[n] + "=" + v.Field(n).String()
		case reflect.Int64:
			paramStr = paramStr + keyNames[n] + "=" + strconv.FormatInt(v.Field(n).Int(), 10)
		case reflect.Float64:
			paramStr = paramStr + keyNames[n] + "=" + strconv.FormatFloat(v.Field(n).Float(), 'E', -1, 64)
		case reflect.Bool:
			paramStr = paramStr + keyNames[n] + "=" + strconv.FormatBool(v.Field(n).Bool())
		}
	}
	return paramStr
}

// Sign 对数据进行签名
func Sign(signatureMethod string, secretKey string, data []byte) (signature string) {
	var h hash.Hash
	switch signatureMethod {
	case "sha1":
		h = hmac.New(sha1.New, []byte(secretKey))
	case "sha256":
		h = hmac.New(sha256.New, []byte(secretKey))
	default:
		h = hmac.New(sha1.New, []byte(secretKey))
	}
	h.Write(data)
	return string(h.Sum(nil))
}
