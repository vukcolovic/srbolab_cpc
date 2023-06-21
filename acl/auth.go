package acl

import (
	"net/http"
)

const (
	ADMIN_ROLE    = "ADMIN_ROLE"
	EMPLOYEE_ROLE = "EMPLOYEE_ROLE"
	CUSTOMER_ROLE = "CUSTOMER_ROLE"
)

func Auth(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}
