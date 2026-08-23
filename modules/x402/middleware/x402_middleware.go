package middleware

import (
	"net/http"
)

func PaymentRequired(
	amount string,
	asset string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				paymentHeader :=
					r.Header.Get("X-Payment")

				if paymentHeader == "" {

					w.Header().Set(
						"X-Payment-Required",
						"true",
					)

					w.Header().Set(
						"X-Payment-Amount",
						amount,
					)

					w.Header().Set(
						"X-Payment-Asset",
						asset,
					)

					http.Error(
						w,
						"payment required",
						http.StatusPaymentRequired,
					)

					return
				}

				// Phase 5:
				// verify x402 payment proof
				// before executing endpoint.

				next.ServeHTTP(w, r)
			},
		)
	}
}
