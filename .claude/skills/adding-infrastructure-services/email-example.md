# Multi-implementation example: email

## `internal/infrastructure/email.go`

```go
package infrastructure

import "context"

type EmailMessage struct {
	To      string
	Subject string
	Body    string
}

type EmailSender interface {
	Send(ctx context.Context, msg EmailMessage) error
}
```

## `internal/infrastructure/email/local.go`

```go
package email

import (
	"context"
	"log/slog"

	"github.com/Sahil2k07/go-base/internal/infrastructure"
)

type Local struct{}

var _ infrastructure.EmailSender = (*Local)(nil)

func NewLocal() *Local { return &Local{} }

func (l *Local) Send(ctx context.Context, msg infrastructure.EmailMessage) error {
	slog.InfoContext(ctx, "email not sent (local provider)", "to", msg.To, "subject", msg.Subject, "body", msg.Body)
	return nil
}
```

## `internal/infrastructure/email/ses.go`

```go
package email

type SES struct {
	// client and sender address from config
}

var _ infrastructure.EmailSender = (*SES)(nil)

func NewSES(cfg configs.SESConfig) *SES { ... }

func (s *SES) Send(ctx context.Context, msg infrastructure.EmailMessage) error { ... }
```

## `internal/apps/email.go`

```go
package apps

import (
	"fmt"

	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
	"github.com/Sahil2k07/go-base/internal/infrastructure/email"
)

func NewEmailSender(cfg configs.EmailConfig, environment string) (infrastructure.EmailSender, error) {
	switch cfg.Provider {
	case "local":
		if environment == "production" {
			return nil, fmt.Errorf("email provider %q is not allowed in production", cfg.Provider)
		}
		return email.NewLocal(), nil
	case "ses":
		return email.NewSES(cfg.SES), nil
	default:
		return nil, fmt.Errorf("unknown email provider %q", cfg.Provider)
	}
}
```

## `configs/server.toml`

```toml
[email]
provider = "local"   # local | ses
```

The `local` provider logs the message body so reset links can be copied during development. That is why it is refused in production.
