# myerrors

Shared error / status-code / global-exception-handler library.

## Files

- `statuscode.go` — `Code` type, sentinel codes, and the single `Code -> HTTP status` mapping.
- `errors.go` — `AppError` type (wraps a cause, carries a `Code`), constructors, `CodeOf(err)` helper.
- `handler.go` — `WriteError`, `RecoverMiddleware` (panic recovery), `WrapHandler` (adapts an
  error-returning handler into a normal `http.HandlerFunc`).
- `example_test.go` — runnable example of the whole flow.

## Usage in a product

```go
mux := http.NewServeMux()
cfg := myerrors.HandlerConfig{Logger: slog.Default()}

mux.Handle("/users", myerrors.WrapHandler(getUser, cfg))

var handler http.Handler = mux
handler = myerrors.RecoverMiddleware(cfg)(handler) // wrap once, at the top

http.ListenAndServe(":8080", handler)
```

```go
func getUser(w http.ResponseWriter, r *http.Request) error {
    if !valid(r) {
        return myerrors.Validation("id is required")
    }
    user, err := repo.Get(id)
    if err != nil {
        return myerrors.Wrap(err, myerrors.CodeInternal, "could not load user").
            WithOp("UserService.Get")
    }
    return json.NewEncoder(w).Encode(user)
}
```

Branch on failures with `myerrors.CodeOf(err)` or `errors.Is(err, myerrors.NotFound(""))` —
never on `err.Error()` strings.

## Versioning workflow

1. This module has its own `go.mod`, versioned independently via git tags (`v1.0.0`, `v1.1.0`, ...).
2. **Codes are append-only.** Adding a code = MINOR. Changing what a code means or its HTTP
   status = MAJOR.
3. `HandlerConfig` and `ErrorResponse` grow by adding fields only — never rename/remove a field
   without a MAJOR bump.
4. At v2+, per Go's module rules, bump the module path too:
   `module github.com/yourcompany/myerrors/v2`, and consumers import
   `.../myerrors/v2`. This lets two products run different majors side by side during migration.
5. Consumers pin a version in their own `go.mod`:
   `require github.com/yourcompany/myerrors v1.2.0`, upgraded deliberately via
   `go get github.com/yourcompany/myerrors@v1.2.0`.
6. Private repo: consumers need `GOPRIVATE=github.com/yourcompany/*` and git auth configured,
   since `go get` clones directly rather than going through the public module proxy.

Tag releases from CI once `main` passes tests, and keep a `CHANGELOG.md` entry per tag so teams
can decide when to upgrade.