// Package applog monta o logger da API (slog) e leva o id da requisição do
// contexto para cada linha de log.
package applog

import (
	"context"
	"io"
	"log/slog"
)

// New cria o logger: JSON em produção, para ferramenta de log ler campo a
// campo; texto em desenvolvimento, para ler no terminal.
func New(production bool, w io.Writer) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	var h slog.Handler = slog.NewTextHandler(w, opts)
	if production {
		h = slog.NewJSONHandler(w, opts)
	}
	return slog.New(contextHandler{h})
}

type requestIDKey struct{}

// WithRequestID guarda o id da requisição no contexto.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID devolve o id da requisição guardado no contexto, ou "".
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// contextHandler acrescenta request_id em todo log feito com o contexto da
// requisição (slog.InfoContext, slog.WarnContext...).
type contextHandler struct{ slog.Handler }

func (h contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := RequestID(ctx); id != "" {
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return contextHandler{h.Handler.WithAttrs(attrs)}
}

func (h contextHandler) WithGroup(name string) slog.Handler {
	return contextHandler{h.Handler.WithGroup(name)}
}
