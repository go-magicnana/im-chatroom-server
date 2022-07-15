package util

import "encoding/json"

func ToJsonString(v any) string {
	b, e := json.Marshal(v)

	Panic(e)

	return string(b)

}
