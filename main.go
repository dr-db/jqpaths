package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
)

var simpleKeyRegex *regexp.Regexp

func init() {
	simpleKeyRegex = regexp.MustCompile(`^[a-zA-z][a-zA-z0-9]*$`)
}

func stringPathPart(k string) string {
	if simpleKeyRegex.MatchString(k) {
		return "." + k
	}
	return fmt.Sprintf("[%q]", k)
}

func showPaths(prefix string, v interface{}) {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, v := range t {
			showPaths(prefix+stringPathPart(k), v)
		}
	case []interface{}:
		for i, v := range t {
			showPaths(fmt.Sprintf("%s[%d]", prefix, i), v)
		}
	case bool, float64:
		fmt.Printf("%s = %v\n", prefix, v)
	case nil:
		fmt.Printf("%s = null\n", prefix)
	case string:
		fmt.Printf("%s = %q\n", prefix, v)
	default:
		fmt.Printf("%s = UNKNOWN TYPE (%T)=%v", prefix, t, t)
	}
}

func main() {
	var v interface{}
	if err := json.NewDecoder(os.Stdin).Decode(&v); err != nil {
		log.Fatal(err)
	}
	showPaths("", v)
}
