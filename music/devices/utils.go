package devices

import (
	"crypto/md5"
	"fmt"
)

func RemoveAfterNil(str string) string {
	breakingIndex := 0
	for i, r := range str {
		breakingIndex = i
		if r == 0 {
			break
		}
	}

	return str[:breakingIndex]
}

func BuildID(deviceId string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(deviceId)))
}