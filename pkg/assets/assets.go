package assets

import (
	_ "embed"
)

//go:embed openapi.json
var OpenAPI []byte

//go:embed scalar.html
var ScalarUI []byte