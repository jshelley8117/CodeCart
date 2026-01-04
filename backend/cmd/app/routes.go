package main

import (
	"net/http"

	"github.com/jshelley8117/CodeCart/internal/handler"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/service"
)

func SetupRoutes(mux *http.ServeMux, resourceConfig ResourceConfig) {
	// ---------- USERS DOMAIN ----------
	userPersistence := persistence.NewUserPersistence(resourceConfig.GCloudDB)
	userService := service.NewUserService(userPersistence)
	userHandler := handler.NewUserHandler(userService)

	mux.HandleFunc("POST /api/v1/users", userHandler.HandleCreateUser)

	// ---------- CUSTOMERS DOMAIN ----------
	customerPersistence := persistence.NewCustomerPersistence(resourceConfig.GCloudDB)
	customerService := service.NewCustomerService(customerPersistence)
	customerHandler := handler.NewCustomerHandler(customerService)

	mux.HandleFunc("POST /api/v1/customers", customerHandler.HandleCreateCustomer)

	// ---------- ADDRESS DOMAIN ----------
	addressPersistence := persistence.NewAddressPersistence(resourceConfig.GCloudDB)
	addressService := service.NewAddressService(addressPersistence)
	addressHandler := handler.NewAddressHandler(addressService)

	mux.HandleFunc("POST /api/v1/addresses", addressHandler.HandleCreateAddress)

}
