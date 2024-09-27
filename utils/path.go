package utils

import "strings"

func CleanPath(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	if strings.HasSuffix(p, "/") {
		p = p[:len(p)-1]
	}

	return p
}

func CleanDirPath(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}

	return p
}
