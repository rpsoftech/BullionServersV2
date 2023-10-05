package interfaces

const (
	ERROR_INTERNAL_SERVER                       = 99
	ERROR_ENTITY_NOT_FOUND                      = 100
	ERROR_OTP_EXPIRED                           = 201
	ERROR_OTP_INVALID                           = 202
	ERROR_OTP_REQ_NOT_EXIST                     = 203
	ERROR_RECENT_OTP_REQ_EXIST                  = 204
	ERROR_INVALID_TOKEN_SIGNATURE               = 205
	ERROR_TOKEN_EXPIRED                         = 206
	ERROR_TOKEN_NOT_BEFORE                      = 207
	ERROR_ROLE_NOT_EXISTS                       = 208
	ERROR_ROLE_NOT_AUTHORIZED                   = 209
	ERROR_USER_NOT_FOUND                        = 210
	ERROR_INVALID_PASSWORD                      = 211
	ERROR_DO_NOT_HAVE_ACCESS_TO_ANY_LOCATION_ID = 212
	ERROR_PERMISSION_NOT_ALLOWED                = 213
	ERROR_INVALID_INPUT                         = 214
	ERROR_GENERAL_USER_REQ_NOT_FOUND            = 215
	ERROR_GENERAL_USER_REQ_REJECTED             = 216
	ERROR_GENERAL_USER_REQ_PENDING              = 217
	ERROR_GENERAL_USER_REQ_EXISTS               = 218
	ERROR_TOKEN_ROLE_NOT_FOUND                  = 219
	ERROR_GENERAL_USER_INVALID_PASSWORD         = 220
	ERROR_INVALID_ENTITY                        = 220
)
