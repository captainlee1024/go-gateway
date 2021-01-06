package public

import "encoding/json"

func ObjectToJson(s interface{}) string {
	bts, _ := json.Marshal(s)
	return string(bts)
}
