package filterbase

import "github.com/google/cel-go/cel"

// initCELProgram - Initialize a CEL program with given CEL command and a set of environments
func initCELProgram(celCmd string, options ...cel.EnvOption) (cel.Program, *cel.Ast, error) {
	env, err := cel.NewEnv(options...)
	if err != nil {
		return nil, nil, err
	}

	ast, issue := env.Compile(celCmd)
	if issue.Err() != nil {
		return nil, nil, issue.Err()
	}

	program, err := env.Program(ast)
	if err != nil {
		return nil, nil, err
	}

	return program, ast, nil
}
