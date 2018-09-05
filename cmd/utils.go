package cmd

import "fmt"

type env struct {
	key   string
	value string
}

func (e *env) String() string {
	return fmt.Sprintf("%v", *e)
}
