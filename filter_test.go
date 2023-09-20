package filterbase

import (
	"context"
	"github.com/google/cel-go/cel"
	"reflect"
	"testing"
)

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

func TestFilterUsers(t *testing.T) {
	type args struct {
		ctx    context.Context
		rows   []*User
		celCmd string
	}
	tests := []struct {
		name    string
		args    args
		want    []*User
		wantErr bool
	}{
		{
			"empty filter string",
			args{context.Background(),
				[]*User{
					{Name: "user123", Age: 1, Admin: true}}, ""},
			[]*User{
				{Name: "user123", Age: 1, Admin: true},
			},
			false,
		},
		{
			"empty filter string 2",
			args{context.Background(),
				[]*User{
					{Name: "user123", Age: 1, Admin: true}}, ``},
			[]*User{
				{Name: "user123", Age: 1, Admin: true},
			},
			false,
		},
		{
			"1 arg",
			args{context.Background(),
				[]*User{
					{Name: "user123", Age: 1, Admin: true}}, `name == "user321" || name == "user123"`},
			[]*User{
				{Name: "user123", Age: 1, Admin: true},
			},
			false,
		},
		{
			"1 arg not ok",
			args{context.Background(),
				[]*User{
					{Name: "user123", Age: 2, Admin: true},
				},
				`name == "user123" && admin == false`},
			[]*User{},
			false,
		},
		{
			"3 args",
			args{context.Background(),
				[]*User{
					{Name: "user123", Age: 3, Admin: true},
				},
				`name == "user123" && admin == true && age > 2`},
			[]*User{
				{Name: "user123", Age: 3, Admin: true},
			},
			false,
		},
		{
			"3 args not ok",
			args{context.Background(),
				[]*User{
					{Name: "user123", Age: 4, Admin: true},
				},
				`name == "user123" && admin == true && age == 2`},
			[]*User{},
			false,
		},
		{
			"bad variable name",
			args{context.Background(),
				[]*User{
					{Name: "user123", Age: 4, Admin: true},
				},
				`name == "user123" && public == true`},
			nil,
			true,
		},
		{
			"not bool result",
			args{context.Background(),
				[]*User{
					{Name: "user123", Age: 4, Admin: true},
				},
				`"hello"`},
			nil,
			true,
		},
		{
			"invalid cel expression",
			args{context.Background(),
				[]*User{
					{Name: "user123", Age: 4, Admin: true},
				},
				`"hello+`},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FilterUsers(tt.args.ctx, tt.args.rows, tt.args.celCmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}
