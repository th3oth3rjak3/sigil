package backends

import "sigil/internal/ast"

type CompilerBackend interface {
	Execute(program *ast.Program, debug bool) error
}
