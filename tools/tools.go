//go:build tools
// +build tools

// Package tools is a dummy package that will be ignored for builds, but included for dependencies
package tools

import (
	_ "github.com/matryer/moq" // moq is used for generating mocks from interfaces.
)
