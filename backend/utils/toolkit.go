package utils

import "regexp"

func CreateFileKey(username, repository, tag string) string {
	return username + "/" + repository + "/" + tag + ".tar.gz"
}

func IsValidName(name string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9\-_]+$`).MatchString(name)
}
