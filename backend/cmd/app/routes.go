package main

import (
	"net/http"
	"os"

	"github.com/jshelley8117/CodeCart/internal/client"
	"github.com/jshelley8117/CodeCart/internal/handler"
	"github.com/jshelley8117/CodeCart/internal/middleware"
	"github.com/jshelley8117/CodeCart/internal/persistence"
	"github.com/jshelley8117/CodeCart/internal/service"
)

func SetupRoutes(mux *http.ServeMux, resourceConfig ResourceConfig) {
	authMW := middleware.AuthMiddleware(resourceConfig.FirebaseAuth) // middleware wrapper to be used for JWT protected routes
	adminMW := middleware.AdminOnly                                  // middleware wrapper to be used for admin-only routes

	// ---------- USERS DOMAIN ----------
	userPersistence := persistence.NewUserPersistence(resourceConfig.GCloudDB)
	userService := service.NewUserService(userPersistence, resourceConfig.FirebaseAuth)
	userHandler := handler.NewUserHandler(userService)

	mux.HandleFunc("POST /api/v1/users", userHandler.HandleCreateUser)

	// ---------- CUSTOMERS DOMAIN ----------
	customerPersistence := persistence.NewCustomerPersistence(resourceConfig.GCloudDB)
	customerService := service.NewCustomerService(customerPersistence)
	customerHandler := handler.NewCustomerHandler(customerService)

	mux.Handle("POST /api/v1/customers", authMW(adminMW(http.HandlerFunc(customerHandler.HandleCreateCustomer))))
	mux.Handle("GET /api/v1/customers", authMW(adminMW(http.HandlerFunc(customerHandler.HandleGetAllCustomers))))
	mux.Handle("DELETE /api/v1/customers/{id}", authMW(adminMW(http.HandlerFunc(customerHandler.HandleDeleteCustomerById))))
	mux.Handle("PATCH /api/v1/customers/{id}", authMW(adminMW(http.HandlerFunc(customerHandler.HandleUpdateCustomerById))))

	// ---------- ADDRESS DOMAIN ----------
	addressPersistence := persistence.NewAddressPersistence(resourceConfig.GCloudDB)
	addressService := service.NewAddressService(addressPersistence)
	addressHandler := handler.NewAddressHandler(addressService)

	mux.Handle("POST /api/v1/addresses", authMW(http.HandlerFunc(addressHandler.HandleCreateAddress)))
	mux.Handle("GET /api/v1/addresses", authMW(http.HandlerFunc(addressHandler.HandleGetAllAddressesById)))
	mux.Handle("GET /api/v1/addresses/{id}", authMW(http.HandlerFunc(addressHandler.HandleGetAddressById)))
	mux.Handle("PATCH /api/v1/addresses/{id}", authMW(http.HandlerFunc(addressHandler.HandleUpdateAddressById)))
	mux.Handle("DELETE /api/v1/addresses/{id}", authMW(http.HandlerFunc(addressHandler.HandleDeleteAddressById)))

	// ---------- CLOUD FUNCTION POC DOMAIN ----------
	cloudFunctionClient := client.NewCloudFunctionClient(os.Getenv("GCP_IMP_SA"))
	cloudFunctionService := service.NewCloudFunctionService(
		cloudFunctionClient,
		service.CloudFunctionConfig{
			HelloWorldURL:  os.Getenv("CLOUD_FUNCTION_HELLO_WORLD_URL"),
			HelloWorld2URL: os.Getenv("CLOUD_FUNCTION_HELLO_WORLD_URL_2"),
		},
	)
	cloudFunctionHandler := handler.NewCloudFunctionHandler(cloudFunctionService)

	mux.HandleFunc("GET /api/v1/hw", cloudFunctionHandler.HandleGetHelloWorld)
	mux.HandleFunc("GET /api/v1/hw2", cloudFunctionHandler.HandleGetHelloWorld2)

	// ---------- ORDERS DOMAIN ----------
	orderPersistence := persistence.NewOrderPersistence(resourceConfig.GCloudDB)
	orderService := service.NewOrderService(orderPersistence)
	orderHandler := handler.NewOrderHandler(orderService)

	mux.Handle("POST /api/v1/orders", authMW(http.HandlerFunc(orderHandler.HandleCreateOrder)))
	mux.Handle("GET /api/v1/orders", authMW(http.HandlerFunc(orderHandler.HandleGetAllOrders)))
	mux.Handle("GET /api/v1/orders/{id}", authMW(http.HandlerFunc(orderHandler.HandleGetAllOrders)))
	mux.Handle("PATCH /api/v1/orders/{id}", authMW(http.HandlerFunc(orderHandler.HandleUpdateOrderById)))

	// ---------- PRODUCTS DOMAIN ----------
	productPersistence := persistence.NewProductPersistence(resourceConfig.GCloudDB)
	productService := service.NewProductService(productPersistence)
	productHandler := handler.NewProductHandler(productService)

	mux.Handle("POST /api/v1/products", authMW(adminMW(http.HandlerFunc(productHandler.HandleCreateProduct))))
	mux.HandleFunc("GET /api/v1/products", productHandler.HandleGetAllProducts)
	mux.HandleFunc("GET /api/v1/products/{id}", productHandler.HandleGetProductById)
	mux.HandleFunc("GET /api/v1/products/{id}/variants", productHandler.HandleGetAllProductVariantsByProductId)
	mux.Handle("PATCH /api/v1/products/{id}", authMW(adminMW(http.HandlerFunc(productHandler.HandleUpdateProductById))))
	mux.Handle("PATCH /api/v1/products/variants/{id}", authMW(adminMW(http.HandlerFunc(productHandler.HandleUpdateProductVariantById))))
	mux.Handle("DELETE /api/v1/products/{id}", authMW(adminMW(http.HandlerFunc(productHandler.HandleDeleteProductById))))
	mux.Handle("DELETE /api/v1/products/variants/{id}", authMW(adminMW(http.HandlerFunc(productHandler.HandleDeleteProductVariantById))))

	// ---------- INVENTORY DOMAIN ----------
	inventoryPersistence := persistence.NewInventoryPersistence((resourceConfig.GCloudDB))
	inventoryService := service.NewInventoryService(inventoryPersistence)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	mux.Handle("POST /api/v1/inventory", authMW(adminMW(http.HandlerFunc(inventoryHandler.HandleCreateInventory))))
	mux.Handle("GET /api/v1/inventory", authMW(adminMW(http.HandlerFunc(inventoryHandler.HandleGetAllInventory))))
	mux.Handle("GET /api/v1/inventory/{id}", authMW(adminMW(http.HandlerFunc(inventoryHandler.HandleGetInventoryById))))
	mux.Handle("PATCH /api/v1/inventory/{id}", authMW(adminMW(http.HandlerFunc(inventoryHandler.HandleUpdateInventoryById))))
	mux.Handle("DELETE /api/v1/inventory/{id}", authMW(adminMW(http.HandlerFunc(inventoryHandler.HandleDeleteInventoryById))))

	// ---------- OrderItem DOMAIN ----------
	orderItemPersistence := persistence.NewOrderItemPersistance(resourceConfig.GCloudDB)
	orderItemService := service.NewOrderItemService(orderItemPersistence)
	orderItemHandler := handler.NewOrderItemHandler(orderItemService)

	mux.HandleFunc("POST /api/v1/orders/{orderId}/item", orderItemHandler.HandleCreateOrderItem)
	mux.HandleFunc("GET /api/v1/orders/{orderId}/items", orderItemHandler.HandleGetAllOrderItems)
	mux.HandleFunc("PATCH /api/v1/orders/{orderId}/item/{id}", orderItemHandler.HandleUpdateOrderItemById)
	mux.HandleFunc("DELETE /api/v1/orders/{orderId}/item/{id}", orderItemHandler.HandleDeleteOrderItemById)

}
