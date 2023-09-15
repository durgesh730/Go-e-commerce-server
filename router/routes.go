package router

import (
	"github.com/durgesh730/authenticationInGo/controllers"
	"github.com/durgesh730/authenticationInGo/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// users routes
	router.HandleFunc("/user/signup", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/user/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/getuserData", controllers.GetUsersEndpoint).Methods("GET")
	router.HandleFunc("/user/AddAddress", middleware.AuthMiddleware(controllers.AddAddress)).Methods("PUT")

	//product
	router.HandleFunc("/product/Maleproduct", controllers.MaleCreateProducts).Methods("POST")
	router.HandleFunc("/product/Femaleproduct", controllers.FemaleCreateProducts).Methods("POST")
	router.HandleFunc("/getproduct", controllers.GetProducts).Methods("GET")
	router.HandleFunc("/getQueryproduct", controllers.GetQueryProducts).Methods("GET")
	router.HandleFunc("/getproductbyId", controllers.GetProductsbyId).Methods("GET")

	//cart routes
    router.HandleFunc("/cart/createproducts", middleware.AuthMiddleware(controllers.CreateProductCart)).Methods("POST")
	router.HandleFunc("/cart/getproducts", middleware.AuthMiddleware(controllers.GetProductFromCart)).Methods("GET")
	router.HandleFunc("/cart/deleteproducts", controllers.DeleteProductFromCart).Methods("DELETE")
	router.HandleFunc("/cart/updateproducts", controllers.UpdateProductFromCart).Methods("PUT")

	return router
}
