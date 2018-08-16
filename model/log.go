package twisk

import "context"

// Logger represents logging interface
type Logger interface {
	// source, msg, error, params
	Log(context.Context, string, string, error, map[string]interface{})
}
