package server

import (
	"RMS/handler"
	"RMS/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	chi.Router
}

func SetupRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/rms", func(rms chi.Router) {
		rms.Post("/register", handler.Register)
		rms.Post("/log-in", handler.Login)
		rms.Route("/auth", func(auth chi.Router) {
			auth.Use(middleware.AuthMiddleware)
			auth.Post("/image", handler.UploadImage)
			// auth.Post("/bulk-image", handler.BulkInsert)
			auth.Route("/admin", func(admin chi.Router) {
				admin.Use(middleware.AdminMiddleware)
				// admin.Post("/", handler.CreateAdmin)
				// admin.Get("/", handler.GetAdminDetails)
				admin.Get("/restaurants", handler.GetAllRestaurants)
				admin.Get("/sub-admin", handler.GetSubAdmin)
				admin.Get("/users", handler.GetUsers)
				admin.Post("/sub-admin", handler.CreateSubAdmin)
				admin.Post("/user", handler.CreateUser)
				admin.Post("/restaurants", handler.CreateRestaurants)
				admin.Route("/{restaurantID}", func(restaurant chi.Router) {
					restaurant.Get("/", handler.GetDishes)
					restaurant.Post("/", handler.CreateDishes)
					restaurant.Route("/{dishID}", func(dish chi.Router) {
						dish.Post("/bulk-dishes", handler.BulkInsertDishes)
						dish.Put("/", handler.UpdateDish)
						dish.Delete("/", handler.DeleteDish)
					})
				})
				// admin.Put("/", handler.UpdateAdmin)
				// admin.Delete("/", handler.DeleteAdmin)
				// admin.Route("/{ID}", func(subAdminID chi.Router) {
				//	subAdminID.Delete("/sub_admin", handler.DeleteSubAdmin)
				//	subAdminID.Route("/{ID}", func(restaurantID chi.Router) {
				//	})
				// })
				// admin.Route("/{ID}", func(userID chi.Router) {
				//	userID.Route("/{role}", func(role chi.Router) {
				//		role.Post("/set_role", handler.SetRole)
				//	})
				// })
			})

			auth.Route("/sub-admin", func(subAdmin chi.Router) {
				subAdmin.Use(middleware.SubAdminMiddleware)
				subAdmin.Post("/restaurants", handler.CreateRestaurants)
				subAdmin.Get("/all-restaurants", handler.GetAllRestaurants)
				subAdmin.Get("/restaurants", handler.GetRestaurants)
				// subAdmin.Post("/user", handler.CreateUser)
				subAdmin.Route("/{restaurantID}", func(restaurant chi.Router) {
					restaurant.Get("/dishes", handler.GetDishes)
					restaurant.Post("/dishes", handler.CreateDishes)
					restaurant.Route("/{dishID}", func(dish chi.Router) {
						dish.Post("/bulk-dishes", handler.BulkInsertDishes)
						dish.Put("/dish", handler.UpdateDish)
						dish.Delete("/dish", handler.DeleteDish)
					})
				})

				// subAdmin.Post("/add-dishes", handler.AddDishes)
				// subAdmin.Get("/dishes", handler.SubAdminGetDishes)
				// subAdmin.Put("/sub_admin", handler.UpdateSubAdmin)
				// subAdmin.Put("/restaurants", handler.UpdateRestaurants)
				// subAdmin.Delete("/sub_admin", handler.DeleteSubAdmin)
				// subAdmin.Delete("/restaurants", handler.DeleteRestaurants)
			})

			auth.Route("/user", func(user chi.Router) {
				user.Use(middleware.AuthMiddleware)
				user.Get("/restaurants", handler.GetAllRestaurants)
				user.Post("/add-address", handler.AddAddress)
				user.Put("/user", handler.UpdateUser)
				user.Delete("/", handler.DeleteUser)
				user.Route("/{restaurantID}", func(restaurant chi.Router) {
					restaurant.Get("/dishes", handler.GetDishes)
					restaurant.Route("/{addressID}", func(address chi.Router) {
						restaurant.Get("/distance", handler.GetDistance)
					})
				})
				// user.Get("/user", handler.GetUser)
			})
		})
	})
	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}

// sqrl for bulk insert
