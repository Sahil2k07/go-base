---
paths:
  - "internal/features/auth/**"
  - "internal/middlewares/**"
  - "internal/infrastructure/crypto.go"
---

# Auth and security

- Access token: HS256-signed, then A256GCM-encrypted JWT with user id, profile id and email. `encryption_secret` must be exactly 32 bytes.
- Refresh token: HS256-signed JWT (user id, session id, random jti) with its own secret. Only its SHA-256 hash is stored (`user_sessions.refresh_token_hash`). An access token must never verify as a refresh token or the reverse.
- Rotation: one session per login (`session_identifier`, a uuid generated in Go). Refresh swaps the stored hash with one conditional update (`RotateUserSession`). If no row matched, the token was reused or expired: delete the session and return 401.
- Signin returns the same 401 message for an unknown email and a wrong password.
- Passwords: bcrypt with `bcrypt_cost` from config; 8 to 72 characters with upper, lower, digit and symbol; checked on signup only. Hash in the service, never in the store.
- Roles and permissions are not in tokens. Middleware checks the database on every request so revocation is immediate.
- Status codes: missing or invalid token 401; authenticated but lacking role/permission 403; a failed role/permission lookup 500, never 403.
- `RequireRoles` passes on any one role, `RequirePermissions` needs all. Both run after `AuthMiddleware()`. Read claims with `infrastructure.GetAuthClaims(c)`.
- Never log or return passwords, tokens, hashes or secrets, and keep `json:"-"` on them.
- Changes to tokens, sessions or these middlewares need unit and integration tests (rotation and reuse live in `tests/integration/features/auth_test.go`).
