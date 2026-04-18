package ast

const (
	Empty = ""
	Null  = "null"
	True  = "true"
	False = "false"

	SimpleCodeStr    = "code"
	SimpleMessageStr = "message"
	SimpleDataStr    = "data"
	SimpleSuccessStr = "success"
	SimpleErrorStr   = "error"
	SimpleUnknownStr = "unknown"

	SimplePageStr        = "page"
	SimpleLimitStr       = "limit"
	SimpleOffsetStr      = "offset"
	SimpleTotalStr       = "total"
	SimpleIdStr          = "id"
	SimpleNameStr        = "name"
	SimpleDescriptionStr = "description"
	SimpleTypeStr        = "type"
	SimpleVersionStr     = "version"
	SimpleStatusStr      = "status"
	SimpleUrlStr         = "url"
	SimpleCreateTimeStr  = "create_time"
	SimpleUpdateTimeStr  = "update_time"
	SimpleDeleteTimeStr  = "delete_time"
	SimpleAccountStr     = "account"
	SimpleTokenStr       = "token"
	SimpleExpireTimeStr  = "expire_time"
	SimpleKeyStr         = "key"
	SimpleValStr         = "value"
)

const BitSize = 32 << (^uint(0) >> 63)
