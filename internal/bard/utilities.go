package bard

import (
	"crypto/md5"
	"encoding/hex"
)

func HashConversation(conversation []string) string {
	hash := md5.New()
	for _, message := range conversation {
		hash.Write([]byte(message))
	}
	return hex.EncodeToString(hash.Sum(nil))
}
