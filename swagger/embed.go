package swagger

import _ "embed"

//go:embed redoc.standalone.js
var RedocJS []byte

//go:embed swagger.html
var SwaggerHTML []byte
