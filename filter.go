package filterbase

import (
	"context"
	"github.com/google/cel-go/cel"
	"github.com/pkg/errors"
)

var ErrInvalidCELProgram = errors.New("invalid CEL program")
var ErrInvalidCELResultType = errors.New("invalid CEL expression result, filter expression must evaluate to a bool")

// Filters the slice of [T] by creating and invoking a CEL filter expression for each row
func Filter[T any](ctx context.Context, ss []T, envOptions []cel.EnvOption, fn func(T) map[string]interface{}, celexpression string) ([]T, error) {
	// Empty filter expression, return all rows
	if celexpression == "" {
		return ss, nil
	}

	// Setup CEL program and evaluate expression
	var celEnvOptions = envOptions
	program, ast, err := initCELProgram(celexpression, celEnvOptions...)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidCELProgram, err.Error())
	}

	// Check if output type is bool
	if ast.OutputType().String() != "bool" {
		return nil, errors.Wrapf(ErrInvalidCELResultType, "got %s", ast.OutputType().String())
	}

	rows := []T{}
	for _, row := range ss {
		val, _, err := program.ContextEval(ctx, fn(row))

		if err != nil {
			return nil, err
		}

		if val.Value() == true {
			rows = append(rows, row)
		}
	}

	return rows, nil
}
