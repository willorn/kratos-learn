package server

import (
	"context"
	v1 "kratos-realworld/api/conduit/v1"
	"kratos-realworld/internal/conf"
	"kratos-realworld/internal/pkg/middleware/auth"
	"kratos-realworld/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

func NewSkipRoutersMatcher() selector.MatchFunc {

	skipRouters := map[string]struct{}{
		"/realworld.v1.Conduit/Login":        {},
		"/realworld.v1.Conduit/Register":     {},
		"/realworld.v1.Conduit/GetArticle":   {},
		"/realworld.v1.Conduit/ListArticles": {},
		"/realworld.v1.Conduit/GetComments":  {},
		"/realworld.v1.Conduit/GetTags":      {},
		"/realworld.v1.Conduit/GetProfile":   {},
	}

	return func(ctx context.Context, operation string) bool {
		if _, ok := skipRouters[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, jwtc *conf.JWT, s *service.ConduitService, logger log.Logger) *http.Server {
	// 创建选项切片
	var opts = []http.ServerOption{
		// 注册错误处理函数
		http.ErrorEncoder(errorEncoder),

		// 设置中间件
		http.Middleware(
			// panic 恢复中间件
			recovery.Recovery(),
			// JWT 认证中间件
			selector.Server(auth.JWTAuth(jwtc.Secret)).Match(NewSkipRoutersMatcher()).Build(),
			// 日志记录中间件
			logging.Server(logger),
		),
		// 添加跨域处理逻辑
		http.Filter(
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}),
				handlers.AllowedOrigins([]string{"*"}),
			),
		),
	}

	// 根据配置添加选项
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	// 创建 HTTP 服务器
	srv := http.NewServer(opts...)
	// 注册 gRPC 服务
	v1.RegisterConduitHTTPServer(srv, s)
	return srv
}
