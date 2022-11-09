package utils

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func GetMd5Sum(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
