package utils

import (
	"fmt"

	"github.com/kevinburke/nacl"
	"github.com/spf13/viper"
)

// GetCombinedKeys returns the key, if keys has exactly 1 item.
// If there are more, the first item is expected to be a format string,
// the other items are parameters to that format.
// Example parameters: ["%s %s", "uuid", "name"]
func GetCombinedKeys(keys []string) string {
	switch len(keys) {
	case 0:
		return ""
	case 1:
		return keys[0]
	default:
		values := make([]any, len(keys)-1)
		for i, v := range keys[1:] {
			values[i] = v
		}
		return fmt.Sprintf(keys[0], values...)
	}
}

// Implements the KeepassxcClientProfile interface for viper config
type ViperKeepassxcProfile struct{}

func (p ViperKeepassxcProfile) GetAssocName() string {
	return viper.GetString(ConfigKeypathAssocName)
}

func (p ViperKeepassxcProfile) GetAssocKey() nacl.Key {
	b64String := viper.GetString(ConfigKeypathAssocKey)
	if b64String == "" {
		return nil
	}
	return B64ToNaclKey(b64String)
}

func (p ViperKeepassxcProfile) SetAssoc(name string, key nacl.Key) error {
	viper.Set(ConfigKeypathAssocName, name)
	viper.Set(ConfigKeypathAssocKey, NaclKeyToB64(key))
	return viper.WriteConfig()
}
