package jsonrpc2

// some helper variables
var (
	EmptyArrayBytes, _ = json.Marshal(map[string]interface{}{}) // {}
	EmptyListBytes, _  = json.Marshal([]interface{}{})          // []
)
