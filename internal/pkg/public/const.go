package public

// 需要设置的上下文相关的一些全局 Key
const (
	CtxResponseKey = "response" // response
	CtxUserIDKey   = "userID"   // userID
	CtxUserKey     = "user"

	// 验证相关 Key
	CtxValidatorKey  = "ValidatorKey"
	CtxTranslatorKey = "TranslatorKey"
)

// Context 之外的其他全局 Key
const (
	// Session Key 用于 session 认证
	KeySessionUser      = "user"
	KeyAdminSessionInfo = "AdminSessionInfoKey"
)

// requestlog 中使用
const (
	HeaderTraceID    = "com-header-rid"
	HeaderSpanID     = "com-header-spanid"
	ContextStartTime = "startExecTime"
	ContextTrace     = "trace"
)

const (
	LoadTypeHTTP = 0
	LoadTypeTCP  = 1
	LoadTypeGRPC = 2

	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)

const (
	// 流量统计的 Key
	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	// 流量前缀
	FlowTotal          = "flow_total"
	FlowServicePrefix  = "flow_service_"
	FlowCountAppPrefix = "flow_app_"
)
