package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"srbolab_cpc/handlers"
	"srbolab_cpc/logoped"
	"srbolab_cpc/util"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const (
	ADMIN   = "ADMINISTRATOR"
	TEACHER = "PREDAVAC"
)

var apiRoleMap = make(map[string]string)

func init() {
	//fill route name - role map
	apiRoleMap["users_register"] = ADMIN
	apiRoleMap["users_delete"] = ADMIN
	apiRoleMap["users_update"] = ADMIN
}

func AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "login") || strings.Contains(r.URL.Path, "corporate-ip") || strings.Contains(r.URL.Path, "create-not-verified") || strings.Contains(r.URL.Path, "seminar-days/jmbg") || strings.Contains(r.URL.Path, "/client-test/create") || strings.Contains(r.URL.Path, "/surveys/active") || strings.Contains(r.URL.Path, "/surveys/client-survey/create") || strings.Contains(r.URL.Path, "/seminar-days/teachers/id") || strings.Contains(r.URL.Path, "/partners/list") {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
			return
		}

		claims := jwt.MapClaims{}
		var tokenHeader = r.Header.Get("Authorization")

		json.NewEncoder(w).Encode(r)
		tokenHeader = strings.TrimSpace(tokenHeader)
		if tokenHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			handlers.SetAuthErrorResponse(w, errors.New("Greška autentifikacije"))
			return
		}

		tkn, err := jwt.ParseWithClaims(tokenHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			logoped.ErrorLog.Println("Error no token on auth token: ", err)
			handlers.SetAuthErrorResponse(w, errors.New("Nema tokena"))
			return
		}
		if !tkn.Valid {
			logoped.ErrorLog.Println("Error token is not valid on auth token: ", err)
			handlers.SetAuthErrorResponse(w, errors.New("Token nije validan"))
			return
		}

		if !handlers.IsCorporateIpMethod(r) {
			logoped.ErrorLog.Println("Error token is valid but ip is not corporate: ", err)
			handlers.SetAuthErrorResponse(w, errors.New("Adresa nije firmina!"))
			return
		}

		exp := claims["ExpiresAt"].(float64)
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			handlers.SetAuthErrorResponse(w, errors.New("Token istekao"))
			return
		}

		var userRoles []string
		act := claims["Roles"].([]interface{})
		for _, v := range act {
			userRoles = append(userRoles, v.(string))
		}

		routeName := mux.CurrentRoute(r).GetName()
		roleForRoute, exist := apiRoleMap[routeName]
		if exist && !util.Contains(userRoles, roleForRoute) {
			logoped.ErrorLog.Println("Error there is no permission for route ", routeName, " error: ", err)
			handlers.SetErrorResponse(w, errors.New("Akcija nije dozvoljena, nemate zahtevane permisije"))
			return
		}

		if time.Unix(int64(exp), 0).Before(time.Now().Add(6 * time.Hour)) {
			refreshedClaims := jwt.MapClaims{
				"Id":        claims["Id"],
				"ExpiresAt": time.Now().Add(time.Hour * 8).Unix(),
				"Roles":     claims["Roles"],
			}

			refreshedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshedClaims)
			refreshedSignedToken, err := refreshedToken.SignedString([]byte("secret"))
			if err != nil {
				logoped.ErrorLog.Println("Error login user, error signing token: ", err)
				return
			}

			w.Header().Add("Access-Control-Expose-Headers", "Authorization")
			w.Header().Set("Authorization", refreshedSignedToken)
		}

		next.ServeHTTP(w, r)
	})
}
