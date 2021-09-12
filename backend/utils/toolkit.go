package utils

func CreateFileKey(username, repository, tag string) string {
	return username + "/" + repository + "/" + tag + ".tar"
}
