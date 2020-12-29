// This file should only be run from the go generate utility, never directly.
// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func escape(text string) string {
	result := ""
	for _, r := range text {
		if r == '"' {
			result += "\\\""
		} else if r == '\\' {
			result += "\\\\"
		} else {
			result += string(r)
		}
	}
	return result
}

func main() {
	fh, err := os.Open("LICENSE.md")
	must(err)
	data, err := ioutil.ReadAll(fh)
	must(err)
	must(fh.Close())

	out, err := os.Create("license-data.go~")
	must(err)
	_, err = out.Write([]byte("// ⚠⚠⚠ Code generated by go generate; DO NOT EDIT ⚠⚠⚠\n"))
	must(err)
	_, err = out.Write([]byte("package main\n\nconst (\n"))
	must(err)

	lines := strings.Split(
		strings.ReplaceAll(string(data), "\r", ""),
		"\n")

	_, err = out.Write([]byte("\tlicenseText = "))
	must(err)

	i := 0
	for _, line := range lines {
		if i > 0 {
			_, err = out.Write([]byte(" + \n\t\t"))
		}
		i++
		_, err = out.Write([]byte(fmt.Sprintf("\"%s\\n\"", escape(line))))
		must(err)
	}
	_, err = out.Write([]byte("\n"))

	_, err = out.Write([]byte("\n)\n"))
	must(err)
	must(out.Close())
	must(os.Rename("license-data.go~", "license-data.go"))
}
