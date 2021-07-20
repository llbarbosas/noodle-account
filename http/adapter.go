package http

import (
	"context"
	"log"
	"net/url"

	core "github.com/llbarbosas/noodle-account/core/infra"
	"github.com/llbarbosas/noodle-account/core/util"

	"github.com/gofiber/fiber/v2"
)

type FiberRequest struct {
	fiberCtx *fiber.Ctx
	ctx      context.Context
}

func (r FiberRequest) Body(out interface{}) error {
	return r.fiberCtx.BodyParser(out)
}

func (r FiberRequest) Query(out interface{}) error {
	return r.fiberCtx.QueryParser(out)
}

func (r FiberRequest) Header(key string) string {
	return r.fiberCtx.Get(key, "")
}

func (r FiberRequest) Params(key string) string {
	return r.fiberCtx.Params(key, "")
}

func (r FiberRequest) RawQuery() string {
	return r.fiberCtx.Context().QueryArgs().String()
}

func (r FiberRequest) QueryOne(key string) string {
	return r.fiberCtx.Query(key, "")
}

func (r FiberRequest) Cookie(cookie *core.Cookie) {
	fCookie := &fiber.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Expires:  cookie.Expires,
		Secure:   cookie.Secure,
		HTTPOnly: cookie.HTTPOnly,
		SameSite: cookie.SameSite,
	}

	r.fiberCtx.Cookie(fCookie)
}

func (r FiberRequest) Cookies(key string) string {
	return r.fiberCtx.Cookies(key, "")
}

func (r FiberRequest) ClearCookie(key ...string) {
	r.fiberCtx.ClearCookie(key...)
}

func (r FiberRequest) Context() context.Context {
	return r.ctx
}

func (r FiberRequest) WithContext(ctxFns ...core.PartialContextFn) core.Request {
	for _, f := range ctxFns {
		r.ctx = f(r.ctx)
	}

	return r
}

func (r FiberRequest) IP() string {
	return r.fiberCtx.IP()
}

func (r FiberRequest) Host() string {
	return r.fiberCtx.Hostname()
}

func NewFiberRequest(fiberCtx *fiber.Ctx) core.Request {
	return &FiberRequest{
		fiberCtx: fiberCtx,
		ctx:      context.Background(),
	}
}

func FiberHandlerAdapter(h core.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := NewFiberRequest(c)
		response := h(request)

		log.Printf("%s: %s %d", util.GetFunctionName(h), response.Type, response.StatusCode)

		if len(response.Headers) > 0 {
			for header, value := range response.Headers {
				c.Set(header, value)
			}
		}

		switch response.Type {
		case "redirect":
			if response.RedirectSafe {
				return c.Redirect(SafeRedirect(response.RedirectURL, core.DashboardRoute, c.Hostname()))
			}

			return c.Redirect(response.RedirectURL)
		case "file":
			return c.Status(response.StatusCode).Render(response.File, response.FileBindData)
		// case "json":
		default:
			return c.Status(response.StatusCode).JSON(response.Body)
		}
	}
}

func FiberRenderFile(file string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render(file, "")
	}
}

func SafeRedirect(path, fallback, hostname string) string {
	pathUrl, err := url.Parse(path)

	if err != nil || pathUrl.Host != "" && pathUrl.Host != hostname {
		return fallback
	}

	return path
}
