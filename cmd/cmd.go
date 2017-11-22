package cmd

import (
	"fmt"
)

type MissingParameters struct {
	Param string
}

func (mp MissingParameters) Error() string {
	return fmt.Sprintf("parameter %s is missing", mp.Param)
}

type BodyEmpty struct{}

func (be BodyEmpty) Error() string {
	return fmt.Sprintf("body cannot be empty")
}
