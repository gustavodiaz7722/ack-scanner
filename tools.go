//go:build tools

package tools

// This file pins tool and library dependencies that are not yet imported
// by production code but are required by the project.

import (
	_ "github.com/google/go-github/v60/github"
	_ "github.com/olekukonko/tablewriter"
	_ "golang.org/x/oauth2"
	_ "gopkg.in/yaml.v3"
	_ "pgregory.net/rapid"
)
