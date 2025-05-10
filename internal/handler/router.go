package handler

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	customMiddleware "github.com/invinciblewest/gophermart/internal/middleware"
	"github.com/invinciblewest/gophermart/internal/usecase"
)

func NewRouter(h *Handler, authUseCase usecase.AuthUseCase) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Recoverer)
	r.Use(customMiddleware.LoggerMiddleware)
	r.Use(chiMiddleware.Compress(5))

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", h.RegisterUser)
			r.Post("/login", h.LoginUser)

			withAuth := r.With(customMiddleware.AuthMiddleware(authUseCase))
			withAuth.Post("/orders", h.AddOrder)
			withAuth.Get("/orders", h.GetUserOrders)
			withAuth.Route("/balance", func(withAuth chi.Router) {
				withAuth.Get("/", h.GetUserBalance)
				withAuth.Post("/withdraw", h.WithdrawBalance)
			})
			withAuth.Get("/withdrawals", h.GetWithdrawals)
		})
	})

	return r
}
