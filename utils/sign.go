package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fileserver/conf"
	"strconv"
	"strings"
	"time"
)

func CheckSign(timestamp, sign string) bool {
	nowTimestamp := time.Now().Unix()
	fromTimestamp, _ := strconv.ParseInt(timestamp, 10, 64)
	if (nowTimestamp - fromTimestamp) > int64(60*5) {
		return false
	}
	key1 := conf.QsConfig.Key1
	key2 := conf.QsConfig.Key2
	if key1 == "" {
		key1 = "WISESOFT"
	}
	if key2 == "" {
		key2 = "MOBILE_IM_2019"
	}
	str := "key1=" + key1 + "&timestamp=" + timestamp + "&key2=" + key2
	if MD5Upper(str) == sign {
		return true
	} else {
		return false
	}
}

// 生成32位大写MD5
func MD5Upper(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return strings.ToUpper(hex.EncodeToString(ctx.Sum(nil)))
}
