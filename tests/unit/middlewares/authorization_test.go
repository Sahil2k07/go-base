package middlewares

import (
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sahil2k07/go-base/internal/constants"
	"github.com/Sahil2k07/go-base/internal/middlewares"
	"github.com/Sahil2k07/go-base/tests/unit"
)

const (
	rolesQuery       = `FROM "roles"`
	permissionsQuery = `FROM "permissions"`
)

func TestRequireRoles(t *testing.T) {
	unit.LoadConfig(t)
	header := "Bearer " + unit.AccessToken(t, 7, 1, "a@x.com")

	tests := []struct {
		name  string
		roles []string
		setup func(m sqlmock.Sqlmock)
		want  int
	}{
		{
			name:  "user has the role",
			roles: []string{"admin"},
			setup: func(m sqlmock.Sqlmock) {
				unit.ExpectCount(m, rolesQuery, 1, 7, "admin")
			},
			want: http.StatusOK,
		},
		{
			name:  "user has the second of two roles",
			roles: []string{"admin", "moderator"},
			setup: func(m sqlmock.Sqlmock) {
				unit.ExpectCount(m, rolesQuery, 0, 7, "admin")
				unit.ExpectCount(m, rolesQuery, 1, 7, "moderator")
			},
			want: http.StatusOK,
		},
		{
			name:  "user has none of the roles",
			roles: []string{"admin", "moderator"},
			setup: func(m sqlmock.Sqlmock) {
				unit.ExpectCount(m, rolesQuery, 0, 7, "admin")
				unit.ExpectCount(m, rolesQuery, 0, 7, "moderator")
			},
			want: http.StatusForbidden,
		},
		{
			name:  "database error",
			roles: []string{"admin"},
			setup: func(m sqlmock.Sqlmock) {
				unit.ExpectQueryError(m, rolesQuery, 7, "admin")
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(unit.MockDB(t))

			rec := unit.Serve(t, header, middlewares.AuthMiddleware(), middlewares.RequireRoles(tt.roles...))

			unit.AssertStatus(t, rec, tt.want)
		})
	}
}

func TestRequirePermissions(t *testing.T) {
	unit.LoadConfig(t)
	header := "Bearer " + unit.AccessToken(t, 7, 1, "a@x.com")

	tests := []struct {
		name        string
		permissions []constants.Permission
		setup       func(m sqlmock.Sqlmock)
		want        int
	}{
		{
			name:        "user has all permissions",
			permissions: []constants.Permission{constants.READ_PERMISSION, constants.READ_ROLE},
			setup: func(m sqlmock.Sqlmock) {
				unit.ExpectCount(m, permissionsQuery, 1, 7, "permission.read")
				unit.ExpectCount(m, permissionsQuery, 1, 7, "role.read")
			},
			want: http.StatusOK,
		},
		{
			name:        "second permission is missing",
			permissions: []constants.Permission{constants.READ_PERMISSION, constants.READ_ROLE},
			setup: func(m sqlmock.Sqlmock) {
				unit.ExpectCount(m, permissionsQuery, 1, 7, "permission.read")
				unit.ExpectCount(m, permissionsQuery, 0, 7, "role.read")
			},
			want: http.StatusForbidden,
		},
		{
			name:        "first permission is missing and the check stops",
			permissions: []constants.Permission{constants.READ_PERMISSION, constants.READ_ROLE},
			setup: func(m sqlmock.Sqlmock) {
				unit.ExpectCount(m, permissionsQuery, 0, 7, "permission.read")
			},
			want: http.StatusForbidden,
		},
		{
			name:        "database error",
			permissions: []constants.Permission{constants.READ_PERMISSION},
			setup: func(m sqlmock.Sqlmock) {
				unit.ExpectQueryError(m, permissionsQuery, 7, "permission.read")
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(unit.MockDB(t))

			rec := unit.Serve(t, header, middlewares.AuthMiddleware(), middlewares.RequirePermissions(tt.permissions...))

			unit.AssertStatus(t, rec, tt.want)
		})
	}
}
