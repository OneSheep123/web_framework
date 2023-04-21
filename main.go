package main

import (
	"fmt"
	"strings"
)

func main() {
	path := "/a/b/c"
	split := strings.Split(strings.Trim(path, "/"), "/")
	fmt.Println(split)
}
