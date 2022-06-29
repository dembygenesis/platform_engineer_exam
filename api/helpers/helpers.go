package helpers

func WrapStrInErrMap(str string) map[string][]string {
	w := make(map[string][]string)
	w["errors"] = []string{str}
	return w
}

func WrapErrInErrMap(err error) map[string][]string {
	w := make(map[string][]string)
	w["errors"] = []string{err.Error()}
	return w
}
