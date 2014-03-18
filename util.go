package gogue

import (
  "fmt"
  "strings"
)

func indent(s string) string {
  return strings.Join(strings.Split(fmt.Sprint(s), "\n"), "\n\t")
}

