package interceptor

type contextKey string

const (
	// CtxKeyUser is the context key for user claims
	CtxKeyUser contextKey = "user"
)
