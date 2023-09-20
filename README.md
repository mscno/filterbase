# Filterbase

![ci](https://github.com/mscno/filterbase/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/mscno/filterbase/graph/badge.svg?token=BZ0WIFKBRJ)](https://codecov.io/gh/mscno/filterbase)

Filterbase is a simple, lightweight, and fast library for filtering data in go using CEL expressions.

## Installation

```bash
go get github.com/mscno/filterbase
```

## Usage

In order to use the filtering package you need to implement a simple filter function for your struct.
This is required for the filtering package to be able to filter your data and get the right params and field values from the structs.

Example:
```go
type User struct {
	Name  string
	Age   int
	Admin bool
}

// Filters the slice acs by invoke CEL program for each AC
func FilterUsers(ctx context.Context, acs []*User, celCmd string) ([]*User, error) {
	var celEnvOptions = []cel.EnvOption{
		cel.Variable("name", cel.StringType),
		cel.Variable("age", cel.IntType),
		cel.Variable("admin", cel.BoolType),
	}

	accessorFn := func(row *User) map[string]interface{} {

		return map[string]any{
			"name":  row.Name,
			"age":   row.Age,
			"admin": row.Admin,
		}
	}

	return Filter(ctx, acs, celEnvOptions, accessorFn, celCmd)
}
```
Then you can use the filter function like this:
```go
filteredUsers := FilterUsers(ctx, users, "name == 'John' && age > 18")
```
Or you can create your own specific filter functions for your structs:
```go
func FilterUsersByAge(users []*User, age int) []users {
	ctx := context.Background()
    return FilterUsers(ctx, users, fmt.Sprintf("age > %d", age))
}
```
and use it like this:
```go
filteredUsers := FilterUsersByAge(users, 18)
```
