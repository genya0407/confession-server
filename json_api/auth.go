package json_api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/genya0407/confession-server/usecase"
	"github.com/julienschmidt/httprouter"
)

func requireAuthorizationBearer(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer realm=""`)
	w.WriteHeader(http.StatusUnauthorized)
}

func invalidRequest(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer error="invalid_request"`)
	w.WriteHeader(http.StatusBadRequest)
}

func invalidToken(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
	w.WriteHeader(http.StatusUnauthorized)
}

const token68regexp = `[A-Za-z0-9\-._~+/]+=*`

func fetchBearerToken(s string) (string, error) {
	rep := regexp.MustCompile(fmt.Sprintf(`Bearer (%s)`, token68regexp))
	matched := rep.FindStringSubmatch(s)
	if len(matched) != 2 {
		log.Println(s)
		return "", errors.New("Invalid Request")
	}

	token := matched[1]
	return token, nil
}

func fetchTokenByHeader(r *http.Request) (string, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		log.Println(authorization)
		return "", errors.New("Authorization required")
	}

	token, err := fetchBearerToken(authorization)
	if err != nil {
		log.Println(token)
		return "", errors.New("Invalid Request")
	}

	return token, nil
}

func fetchTokenByQueryParams(r *http.Request) (string, error) {
	tokens, ok := r.URL.Query()["access_token"]
	if !ok {
		return "", errors.New("No token specified")
	}

	token := tokens[0]

	re := regexp.MustCompile(token68regexp)
	if !re.MatchString(token) {
		log.Println(token)
		return "", errors.New("Invalid Request")
	}

	return token, nil
}

func fetchToken(r *http.Request) (string, error) {
	token, headerErr := fetchTokenByHeader(r)
	if headerErr == nil {
		log.Println("Token sent by header")
		return token, nil
	}

	token, paramsErr := fetchTokenByQueryParams(r)
	if paramsErr == nil {
		log.Println("Token sent by parameter")
		return token, nil
	}

	if paramsErr.Error() == "No token specified" {
		log.Println("Authorization by header was tried but failed")
		return "", headerErr
	}

	log.Println("Authorization by query parameter was tried but failed")
	return "", paramsErr
}

func handleError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case "Authorization required":
		requireAuthorizationBearer(w)
	case "Invalid Request":
		invalidRequest(w)
	default:
		panic(err.Error())
	}
}

func AuthorizeAnonymous(route func(http.ResponseWriter, *http.Request, httprouter.Params, usecase.AnonymousLoginInfoDTO)) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := fetchToken(r)
		if err != nil {
			handleError(w, err)
			return
		}

		route(w, r, ps, usecase.AnonymousLoginInfoDTO{SessionToken: token})
	}
}

func AuthorizeAccount(route func(http.ResponseWriter, *http.Request, httprouter.Params, usecase.AccountLoginInfoDTO)) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := fetchToken(r)
		if err != nil {
			handleError(w, err)
			return
		}

		route(w, r, ps, usecase.AccountLoginInfoDTO{SessionToken: token})
	}
}
