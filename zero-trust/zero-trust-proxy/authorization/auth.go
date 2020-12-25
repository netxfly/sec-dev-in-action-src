package authorization

import (
	"net"
	"net/http"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"

	"github.com/dgrijalva/jwt-go"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/logger"
	"sec-dev-in-action-src/zero-trust/zero-trust-proxy/vars"
)

type SecProxy struct {
	cel         *cel.Env
	expressions []cel.Program
	expMap      map[cel.Program]string
}

func inNetwork(clientIP ref.Val, network ref.Val) ref.Val {
	snet, ok := network.Value().(string)
	if !ok {
		return types.False
	}
	sip, ok := clientIP.Value().(string)
	if !ok {
		return types.False
	}

	_, subnet, _ := net.ParseCIDR(snet)
	ip := net.ParseIP(sip)

	if subnet.Contains(ip) {
		return types.True
	}

	return types.False
}

// NewAuthorization creates a new authorization service with a given set of rules
func NewAuthorization(expressions []string) *SecProxy {
	env, err := cel.NewEnv(cel.Declarations(
		decls.NewVar("request.host", decls.String),
		decls.NewVar("request.path", decls.String),
		decls.NewVar("request.ip", decls.String),
		decls.NewVar("request.email", decls.String),
		decls.NewVar("request.time", decls.Timestamp),

		decls.NewFunction("network",
			decls.NewInstanceOverload("network_string_string", []*exprpb.Type{decls.String, decls.String}, decls.String)),
	))
	if err != nil {
		logger.Log.Fatal(err)
	}
	expMap := make(map[cel.Program]string)
	programs := make([]cel.Program, 0, len(expressions))
	for _, exp := range expressions {
		ast, issues := env.Compile(exp)
		if issues != nil && issues.Err() != nil {
			logger.Log.Fatalf("Invalid CEL expression: %s", issues.Err())
		}

		// declare function overloads
		funcs := cel.Functions(
			&functions.Overload{
				Operator: "network",
				Binary:   inNetwork,
			})

		p, err := env.Program(ast, funcs)
		if err != nil {
			logger.Log.Fatalf("Error while creating CEL program: %q", err)
		}

		programs = append(programs, p)
		expMap[p] = exp
	}

	return &SecProxy{
		cel:         env,
		expressions: programs,
		expMap:      expMap,
	}
}

// Middleware evaluates authorization rules against a request
func (p *SecProxy) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debugf("Authorizing request %q", r.URL)

		context := getContext(r)

		for _, exp := range p.expressions {
			out, evalDetail, err := exp.Eval(context)
			logger.Log.Debugf("exp: %v, context: %v, out: %v,  detail: %v, err: %v",
				p.expMap[exp], context, out, evalDetail, err)

			if err != nil {
				logger.Log.Errorf("Error evaluating expression: %v, err: %v", p.expMap[exp], err)
				w.WriteHeader(http.StatusForbidden)
				_, err := w.Write([]byte("鉴权失败"))
				_ = err
				return
			}

			logger.Log.Warningf("exp: %v, content: %v, out.Value: %v", p.expMap[exp], context, out.Value())
			if out.Value() == false {
				_, err := w.Write([]byte("鉴权失败"))
				_ = err
				w.WriteHeader(http.StatusForbidden)
				return
			}

			logger.Log.Debugf("authZ result: %v", out)
		}

		next.ServeHTTP(w, r)
	})
}

func getContext(r *http.Request) map[string]interface{} {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logger.Log.Error(err)
	}
	return map[string]interface{}{
		"request.host":  r.Host,
		"request.path":  r.RequestURI,
		"request.ip":    ip,
		"request.time":  time.Now().UTC().Format(time.RFC3339),
		"request.email": getEmail(r),
	}
}

func getEmail(r *http.Request) string {
	email := ""
	cookie, err := r.Cookie(vars.CookieName)
	tokenString := r.Header.Get(vars.HeaderName)
	if err == http.ErrNoCookie && tokenString == "" {
		return email
	}

	if tokenString == "" && cookie != nil {
		tokenString = cookie.Value
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})

	if token == nil {
		return email
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		email = claims["sub"].(string)
	}
	logger.Log.Warningf("claims: %v, ok: %v, email: %v", claims, ok, email)
	return email
}
