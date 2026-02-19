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
	userPersistence := persistence.NewUserPersistence(resourceConfig.GCloudDB)
	userService := service.NewUserService(userPersistence)
	userHandler := handler.NewUserHandler(userService)

	mux.HandleFunc("POST /api/v1/users", userHandler.HandleCreateUser)

	// ---------- CUSTOMERS DOMAIN ----------
	customerPersistence := persistence.NewCustomerPersistence(resourceConfig.GCloudDB)
	customerService := service.NewCustomerService(customerPersistence)
	customerHandler := handler.NewCustomerHandler(customerService)

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
	mux.HandleFunc("GET /api/v1/addresses/{id}", addressHandler.HandleGetAddressById)
	mux.HandleFunc("PATCH /api/v1/addresses/{id}", addressHandler.HandleUpdateAddressById)
	mux.HandleFunc("DELETE /api/v1/addresses/{id}", addressHandler.HandleDeleteAddressById)

	// ---------- CLOUD FUNCTION POC DOMAIN ----------
	cloudFunctionClient := client.NewCloudFunctionClient(resourceConfig.TokenSource, os.Getenv("GCP_IMP_SA"))
	cloudFunctionService := service.NewCloudFunctionService(
		cloudFunctionClient,
		os.Getenv("CLOUD_FUNCTION_HELLO_WORLD_URL"),
	)
	cloudFunctionHandler := handler.NewCloudFunctionHandler(cloudFunctionService)

	mux.HandleFunc("GET /api/v1/hw", cloudFunctionHandler.HandleGetHelloWorld)

	// ---------- ORDERS DOMAIN ----------
	orderPersistence := persistence.NewOrderPersistence(resourceConfig.GCloudDB)
	orderService := service.NewOrderService(orderPersistence)
	orderHandler := handler.NewOrderHandler(orderService)

	mux.HandleFunc("POST /api/v1/orders", orderHandler.HandleCreateOrder)
	mux.HandleFunc("GET /api/v1/orders", orderHandler.HandleGetAllOrders)
	mux.HandleFunc("GET /api/v1/orders/{id}", orderHandler.HandleGetAllOrders)
	mux.HandleFunc("PATCH /api/v1/orders/{id}", orderHandler.HandleUpdateOrderById)

	// ---------- PRODUCTS DOMAIN ----------
	productPersistence := persistence.NewProductPersistence(resourceConfig.GCloudDB)
	productService := service.NewProductService(productPersistence)
	productHandler := handler.NewProductHandler(productService)

	mux.HandleFunc("POST /api/v1/products", productHandler.HandleCreateProduct)
	mux.HandleFunc("GET /api/v1/products", productHandler.HandleGetAllProducts)
	mux.HandleFunc("GET /api/v1/products/{id}", productHandler.HandleGetProductById)
	mux.HandleFunc("GET /api/v1/products/{id}/variants", productHandler.HandleGetAllProductVariantsByProductId)
	mux.HandleFunc("PATCH /api/v1/products/{id}", productHandler.HandleUpdateProductById)
	mux.HandleFunc("PATCH /api/v1/products/variants/{id}", productHandler.HandleUpdateProductVariantById)
	mux.HandleFunc("DELETE /api/v1/products/{id}", productHandler.HandleDeleteProductById)
	mux.HandleFunc("DELETE /api/v1/products/variants/{id}", productHandler.HandleDeleteProductVariantById)

}
