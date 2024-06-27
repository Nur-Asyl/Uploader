package utils

func SaltKey(key, salt string) string {
	return key + "_" + salt
}
