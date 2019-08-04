package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

var simpleKeyRegex *regexp.Regexp

func init() {
	simpleKeyRegex = regexp.MustCompile(`^[a-zA-z][a-zA-z0-9]*$`)
}

func extendPath(prefix, k string) string {
	if simpleKeyRegex.MatchString(k) {
		return prefix + "." + k
	}
	if prefix == "" {
		prefix = "."
	}
	return fmt.Sprintf("%s[%q]", prefix, k)
}

func showPaths(prefix string, v interface{}, results map[string]string) {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, v := range t {
			showPaths(extendPath(prefix, k), v, results)
		}
	case []interface{}:
		for i, v := range t {
			if prefix == "" {
				prefix = "."
			}
			showPaths(fmt.Sprintf("%s[%d]", prefix, i), v, results)
		}
	case bool, float64:
		results[prefix] = fmt.Sprintf("%v", v)
	case nil:
		results[prefix] = "null"
	case string:
		results[prefix] = fmt.Sprintf("%q", v)
	default:
		results[prefix] = fmt.Sprintf("UNKNOWN TYPE (%T)=%v", t, t)
	}
}

func main() {
	var outputStyle string
	flag.StringVar(&outputStyle, "o", "plain", "output format (\"plain\" or \"json\")")
	flag.Parse()

	var v interface{}
	if err := json.NewDecoder(os.Stdin).Decode(&v); err != nil {
		log.Fatal(err)
	}
	results := make(map[string]string)
	showPaths("", v, results)
	switch outputStyle {
	case "plain":
		for k, v := range results {
			fmt.Printf("%s => %s\n", k, v)
		}
	case "json":
		if err := json.NewEncoder(os.Stdout).Encode(results); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown output style: %q", outputStyle)
	}
}
