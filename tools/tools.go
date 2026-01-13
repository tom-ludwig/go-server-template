// tools/tools.go
//go:build tools

package tools

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "github.com/stripe/pg-schema-diff/cmd/pg-schema-diff"
)
