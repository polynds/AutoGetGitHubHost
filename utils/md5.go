package utils

import (
	"crypto/md5"
	"fmt"
)

func Str2md5(str string) string {
	data := []byte(str)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return fmt.Sprintf("%x", md5Ctx.Sum(nil))
}
