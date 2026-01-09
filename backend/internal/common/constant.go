package common

const (
	ERR_REQ_BODY_READ_FAIL         = "Failed to read request body"
	ERR_REQ_UNMARSH_FAIL           = "Failed to Unmarshal JSON to Go Type"
	ERR_VALIDATION_FAIL            = "Validation failed"
	ERR_CLIENT_REQUEST_FAIL        = "Server failed to process request"
	ERR_CLIENT_DB_PERSISTENCE_FAIL = "Failed to save data"
	ERR_CLIENT_DB_RETRIEVAL_FAIL   = "Failed to retrieve requested data"
	ERR_CLIENT_DB_DELETE_FAIL      = "Failed to remove requested data"
)
