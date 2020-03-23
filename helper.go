package jsonrpc2

var (
	EmptyArrayBytes, _ = json.Marshal(map[string]interface{}{}) // {}
	EmptyListBytes, _  = json.Marshal([]interface{}{})          // []
)
