package main

import (
	"net/http"
	"os"

	"github.com/jshelley8117/CodeCart/internal/client"
	"github.com/jshelley8117/CodeCart/internal/handler"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/service"
)

func SetupRoutes(mux *http.ServeMux, resourceConfig ResourceConfig) {
	// ---------- USERS DOMAIN ----------
	userPersistence := persistence.NewUserPersistence(resourceConfig.GCloudDB, resourceConfig.Logger)
	userService := service.NewUserService(userPersistence, resourceConfig.Logger)
	userHandler := handler.NewUserHandler(userService, resourceConfig.Logger)

	mux.HandleFunc("POST /api/v1/users", userHandler.HandleCreateUser)

	// ---------- CUSTOMERS DOMAIN ----------
	customerPersistence := persistence.NewCustomerPersistence(resourceConfig.GCloudDB, resourceConfig.Logger)
	customerService := service.NewCustomerService(customerPersistence, resourceConfig.Logger)
	customerHandler := handler.NewCustomerHandler(customerService, resourceConfig.Logger)

	mux.HandleFunc("POST /api/v1/customers", customerHandler.HandleCreateCustomer)
	mux.HandleFunc("GET /api/v1/customers", customerHandler.HandleGetAllCustomers)
	mux.HandleFunc("DELETE /api/v1/customers/{id}", customerHandler.HandleDeleteCustomerById)
	mux.HandleFunc("PATCH /api/v1/customers/{id}", customerHandler.HandleUpdateCustomerById)

	// ---------- ADDRESS DOMAIN ----------
	addressPersistence := persistence.NewAddressPersistence(resourceConfig.GCloudDB)
	addressService := service.NewAddressService(addressPersistence)
	addressHandler := handler.NewAddressHandler(addressService)

	mux.HandleFunc("POST /api/v1/addresses", addressHandler.HandleCreateAddress)
	mux.HandleFunc("GET /api/v1/addresses", addressHandler.HandleGetAllAddresses)

	// ---------- CLOUD FUNCTION POC DOMAIN ----------
	cloudFunctionClient := client.NewCloudFunctionClient(resourceConfig.TokenSource, os.Getenv("GCP_IMP_SA"), resourceConfig.Logger)
	cloudFunctionService := service.NewCloudFunctionService(
		cloudFunctionClient,
		os.Getenv("CLOUD_FUNCTION_HELLO_WORLD_URL"),
		resourceConfig.Logger,
	)
	cloudFunctionHandler := handler.NewCloudFunctionHandler(cloudFunctionService, resourceConfig.Logger)

	mux.HandleFunc("GET /api/v1/hw", cloudFunctionHandler.HandleGetHelloWorld)
}
