# Test skeletons

Replace `blog` and names with the real ones. Import app packages with an alias; test packages are named after the area.

## Contents

- Validation (unit)
- Service with a fake repository (unit)
- Middleware backed by the database (unit)
- Endpoint (integration)
- Roles in integration tests

## Validation (unit)

`tests/unit/features/blog/views_test.go`

```go
package blog

import (
	"testing"

	appblog "github.com/Sahil2k07/go-base/internal/features/blog"
)

func TestCreateBlogRequestValidate(t *testing.T) {
	valid := appblog.CreateBlogRequest{Title: "Hello"}

	tests := []struct {
		name    string
		mutate  func(*appblog.CreateBlogRequest)
		wantErr bool
	}{
		{"valid", func(*appblog.CreateBlogRequest) {}, false},
		{"empty title", func(r *appblog.CreateBlogRequest) { r.Title = "" }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := valid
			tt.mutate(&req)

			if err := req.ValidateRequest(); (err != nil) != tt.wantErr {
				t.Fatalf("wantErr %v, got %v", tt.wantErr, err)
			}
		})
	}
}
```

## Service with a fake repository (unit)

Embed the interface and override only what the test needs. Any other call panics, which flags an unexpected dependency.

```go
type fakeBlogRepo struct {
	blog.BlogRepository
	titleTaken bool
	err        error
}

func (f fakeBlogRepo) IsTitleTaken(string) (bool, error) { return f.titleTaken, f.err }

func TestCreateBlogConflict(t *testing.T) {
	svc := appblog.NewBlogService(fakeBlogRepo{titleTaken: true})

	_, err := svc.CreateBlog(1, appblog.CreateBlogRequest{Title: "Hello"})

	var appErr *app_errors.AppError
	if !errors.As(err, &appErr) || appErr.Code != http.StatusConflict {
		t.Fatalf("want a 409 AppError, got %v", err)
	}
}
```

Services that use the crypto service: `infrastructure.NewCryptoService(unit.LoadConfig(t).Server)`.

## Middleware backed by the database (unit)

```go
func TestRequireRoles(t *testing.T) {
	unit.LoadConfig(t)
	header := "Bearer " + unit.AccessToken(t, 7, 1, "a@x.com")

	m := unit.MockDB(t)
	unit.ExpectCount(m, `FROM "roles"`, 1, 7, "admin")

	rec := unit.Serve(t, header, appmw.AuthMiddleware(), appmw.RequireRoles("admin"))

	unit.AssertStatus(t, rec, http.StatusOK)
}
```

A database failure: `unit.ExpectQueryError(m, `FROM "roles"`, 7, "admin")`, expecting 500.

## Endpoint (integration)

`tests/integration/features/blog_test.go`. Use the real routes and request bodies; copy the sign-up helper and response types from `auth_test.go` (move them to `helpers_test.go` in the same folder when a second file needs them).

```go
func TestCreateBlog(t *testing.T) {
	env := integration.Setup(t)
	token := signUp(t, env, "a@x.com")

	t.Run("created", func(t *testing.T) {
		status, body := env.Call(t, http.MethodPost, "/api/blogs", token, map[string]string{"title": "Hello"})
		if status != http.StatusCreated {
			t.Fatalf("want 201, got %d: %s", status, body)
		}
	})

	t.Run("no token", func(t *testing.T) {
		status, _ := env.Call(t, http.MethodPost, "/api/blogs", "", map[string]string{"title": "Hello"})
		if status != http.StatusUnauthorized {
			t.Fatalf("want 401, got %d", status)
		}
	})
}
```

## Roles in integration tests

Roles are read from the database on every request, so seed them after sign-up and the existing token works.

```go
var user models.User
env.DB.Where("email = ?", "a@x.com").First(&user)

role := models.Role{Name: "admin"}
env.DB.Create(&role)
env.DB.Create(&models.UserRole{UserID: user.ID, RoleID: role.ID})
```
