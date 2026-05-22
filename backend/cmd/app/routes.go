package main

import (
	"net/http"
	"os"

	"github.com/jshelley8117/CodeCart/internal/client"
	"github.com/jshelley8117/CodeCart/internal/common"
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

	mux.HandleFunc("POST "+common.APIBasePath+"/users", userHandler.HandleCreateUser)

	// ---------- CUSTOMERS DOMAIN ----------
	customerPersistence := persistence.NewCustomerPersistence(resourceConfig.GCloudDB)
	customerService := service.NewCustomerService(customerPersistence)
	customerHandler := handler.NewCustomerHandler(customerService)

	mux.Handle("POST "+common.APIBasePath+"/customers", authMW(adminMW(http.HandlerFunc(customerHandler.HandleCreateCustomer))))
	mux.Handle("GET "+common.APIBasePath+"/customers", authMW(adminMW(http.HandlerFunc(customerHandler.HandleGetAllCustomers))))
	mux.Handle("DELETE "+common.APIBasePath+"/customers/{id}", authMW(adminMW(http.HandlerFunc(customerHandler.HandleDeleteCustomerById))))
	mux.Handle("PATCH "+common.APIBasePath+"/customers/{id}", authMW(adminMW(http.HandlerFunc(customerHandler.HandleUpdateCustomerById))))

	// ---------- ADDRESS DOMAIN ----------
	addressPersistence := persistence.NewAddressPersistence(resourceConfig.GCloudDB)
	addressService := service.NewAddressService(addressPersistence)
	addressHandler := handler.NewAddressHandler(addressService)

	mux.Handle("POST "+common.APIBasePath+"/addresses", authMW(http.HandlerFunc(addressHandler.HandleCreateAddress)))
	mux.Handle("GET "+common.APIBasePath+"/addresses", authMW(http.HandlerFunc(addressHandler.HandleGetAllAddressesById)))
	mux.Handle("GET "+common.APIBasePath+"/addresses/{id}", authMW(http.HandlerFunc(addressHandler.HandleGetAddressById)))
	mux.Handle("PATCH "+common.APIBasePath+"/addresses/{id}", authMW(http.HandlerFunc(addressHandler.HandleUpdateAddressById)))
	mux.Handle("DELETE "+common.APIBasePath+"/addresses/{id}", authMW(http.HandlerFunc(addressHandler.HandleDeleteAddressById)))

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

	mux.HandleFunc("GET "+common.APIBasePath+"/hw", cloudFunctionHandler.HandleGetHelloWorld)
	mux.HandleFunc("GET "+common.APIBasePath+"/hw2", cloudFunctionHandler.HandleGetHelloWorld2)

	// ---------- ORDERS DOMAIN ----------
	orderPersistence := persistence.NewOrderPersistence(resourceConfig.GCloudDB)
	orderService := service.NewOrderService(orderPersistence)
	orderHandler := handler.NewOrderHandler(orderService)

	mux.Handle("POST "+common.APIBasePath+"/orders", authMW(http.HandlerFunc(orderHandler.HandleCreateOrder)))
	mux.Handle("GET "+common.APIBasePath+"/orders", authMW(http.HandlerFunc(orderHandler.HandleGetAllOrders)))
	mux.Handle("GET "+common.APIBasePath+"/orders/{id}", authMW(http.HandlerFunc(orderHandler.HandleGetAllOrders)))
	mux.Handle("PATCH "+common.APIBasePath+"/orders/{id}", authMW(http.HandlerFunc(orderHandler.HandleUpdateOrderById)))

	// ---------- PRODUCTS DOMAIN ----------
	productPersistence := persistence.NewProductPersistence(resourceConfig.GCloudDB)
	productService := service.NewProductService(productPersistence)
	productHandler := handler.NewProductHandler(productService)

	mux.Handle("POST "+common.APIBasePath+"/products", authMW(adminMW(http.HandlerFunc(productHandler.HandleCreateProduct))))
	mux.HandleFunc("GET "+common.APIBasePath+"/products", productHandler.HandleGetAllProducts)
	mux.HandleFunc("GET "+common.APIBasePath+"/products/{id}", productHandler.HandleGetProductById)
	mux.HandleFunc("GET "+common.APIBasePath+"/products/{id}/variants", productHandler.HandleGetAllProductVariantsByProductId)
	mux.Handle("PATCH "+common.APIBasePath+"/products/{id}", authMW(adminMW(http.HandlerFunc(productHandler.HandleUpdateProductById))))
	mux.Handle("PATCH "+common.APIBasePath+"/products/variants/{id}", authMW(adminMW(http.HandlerFunc(productHandler.HandleUpdateProductVariantById))))
	mux.Handle("DELETE "+common.APIBasePath+"/products/{id}", authMW(adminMW(http.HandlerFunc(productHandler.HandleDeleteProductById))))
	mux.Handle("DELETE "+common.APIBasePath+"/products/variants/{id}", authMW(adminMW(http.HandlerFunc(productHandler.HandleDeleteProductVariantById))))

	// ---------- INVENTORY DOMAIN ----------
	inventoryPersistence := persistence.NewInventoryPersistence((resourceConfig.GCloudDB))
	inventoryService := service.NewInventoryService(inventoryPersistence)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	mux.Handle("POST "+common.APIBasePath+"/inventory", authMW(adminMW(http.HandlerFunc(inventoryHandler.HandleCreateInventory))))
	mux.Handle("GET "+common.APIBasePath+"/inventory", authMW(adminMW(http.HandlerFunc(inventoryHandler.HandleGetAllInventory))))
	mux.Handle("GET "+common.APIBasePath+"/inventory/{id}", authMW(adminMW(http.HandlerFunc(inventoryHandler.HandleGetInventoryById))))
	mux.Handle("PATCH "+common.APIBasePath+"/inventory/{id}", authMW(adminMW(http.HandlerFunc(inventoryHandler.HandleUpdateInventoryById))))
	mux.Handle("DELETE "+common.APIBasePath+"/inventory/{id}", authMW(adminMW(http.HandlerFunc(inventoryHandler.HandleDeleteInventoryById))))

	// ---------- PAYMENTS DOMAIN ----------
	paymentPersistence := persistence.NewPaymentPersistence(resourceConfig.GCloudDB)
	paymentService := service.NewPaymentService(paymentPersistence)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	mux.Handle("POST /api/v1/payments", authMW(http.HandlerFunc(paymentHandler.HandleCreatePayment)))
	mux.Handle("GET /api/v1/payments/{id}", authMW(http.HandlerFunc(paymentHandler.HandleGetPaymentById)))
	mux.Handle("PATCH /api/v1/payments/{id}", authMW(http.HandlerFunc(paymentHandler.HandleUpdatePaymentById)))
	mux.Handle("DELETE /api/v1/payments/{id}", authMW(http.HandlerFunc(paymentHandler.HandleDeletePaymentById)))

	// ---------- OrderItem DOMAIN ----------
	orderItemPersistence := persistence.NewOrderItemPersistance(resourceConfig.GCloudDB)
	orderItemService := service.NewOrderItemService(orderItemPersistence)
	orderItemHandler := handler.NewOrderItemHandler(orderItemService)

	mux.HandleFunc("POST "+common.APIBasePath+"/orders/{orderId}/item", orderItemHandler.HandleCreateOrderItem)
	mux.HandleFunc("GET "+common.APIBasePath+"/orders/{orderId}/items", orderItemHandler.HandleGetAllOrderItems)
	mux.HandleFunc("PATCH "+common.APIBasePath+"/orders/{orderId}/item/{id}", orderItemHandler.HandleUpdateOrderItemById)
	mux.HandleFunc("DELETE "+common.APIBasePath+"/orders/{orderId}/item/{id}", orderItemHandler.HandleDeleteOrderItemById)

	// ---------- ORDER DISCOUNTS DOMAIN ----------
	orderDiscountPersistence := persistence.NewOrderDiscountPersistence(resourceConfig.GCloudDB)
	orderDiscountService := service.NewOrderDiscountService(orderDiscountPersistence)
	orderDiscountHandler := handler.NewOrderDiscountHandler(orderDiscountService)

	mux.Handle("POST "+common.APIBasePath+"/orders/{orderId}/discount", authMW(http.HandlerFunc(orderDiscountHandler.HandleCreateOrderDiscount)))
	mux.Handle("GET "+common.APIBasePath+"/orders/{orderId}/discount/{id}", authMW(http.HandlerFunc(orderDiscountHandler.HandleGetOrderDiscountById)))

	// ---------- LOCATIONS DOMAIN ----------
	locationPersistence := persistence.NewLocationPersistence(resourceConfig.GCloudDB)
	locationService := service.NewLocationService(locationPersistence)
	locationHandler := handler.NewLocationHandler(locationService)

	mux.Handle("POST /api/v1/locations", authMW(adminMW(http.HandlerFunc(locationHandler.HandleCreateLocation))))
	mux.HandleFunc("GET /api/v1/locations", locationHandler.HandleGetAllLocations)
	mux.HandleFunc("GET /api/v1/locations/{id}", locationHandler.HandleGetLocationById)
	mux.Handle("PATCH /api/v1/locations/{id}", authMW(adminMW(http.HandlerFunc(locationHandler.HandleUpdateLocationById))))
	mux.Handle("DELETE /api/v1/locations/{id}", authMW(adminMW(http.HandlerFunc(locationHandler.HandleDeleteLocationById))))

}
