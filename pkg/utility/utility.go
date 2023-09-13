package utility

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func HashUrlToFileName(url string) string {
	hash := sha1.New()
	hash.Write([]byte(url))
	return hex.EncodeToString(hash.Sum(nil))
}

func ExtractProbableFileEnding(url string) string {
	split := strings.Split(url, "/")
	fileNameSplit := strings.Split(split[len(split)-1], ".")
	probableFileEnding := fileNameSplit[len(fileNameSplit)-1]
	if len(probableFileEnding) > 4 {
		return "unknown"
	}
	return probableFileEnding
}
