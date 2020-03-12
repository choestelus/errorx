// Package errorx provides frequently used error constructs in API development
package errorx

import "fmt"

// E defines common error information for inspecting
// and displaying to various format
type E struct {
	e      error
	Code   string
	Fields map[string]interface{}
}

func (e E) Error() string {
	if e.e != nil {
		return e.Error()
	}
	return fmt.Sprintf("%+v", e)
}
