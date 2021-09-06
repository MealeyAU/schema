package file

import "strings"

type Path string

// Extension returns the file extension of a given path
func (fp Path) Extension() string {
	parts := strings.Split(string(fp), ".")
	return parts[len(parts) - 1]
}

// Parent returns the parent of the given path
func (fp Path) Parent() Path {
	withoutExtension := strings.Split(string(fp), ".")[0]
	pathComponents := strings.Split(withoutExtension, "/")

	componentLen := len(pathComponents)
	if componentLen <= 1 {
		return "/"
	} else {
		return Path(strings.Join(pathComponents[0:componentLen-1], "/"))
	}
}


