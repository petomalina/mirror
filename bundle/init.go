package bundle

import "github.com/petomalina/mirror"

var (
	// L is a global logger that can be reconfigured by third parties
	// to customize logging
	// TODO: copying from the core library is kinda hacky, maybe create a separate log?
	L = mirror.L
)
