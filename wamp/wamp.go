package wamp

const (
	Hello        = 1
	Welcome      = 2
	Abort        = 3
	Challenge    = 4
	Authenticate = 5
	Goodbye      = 6
	Error        = 8
	Publish      = 16
	Published    = 17
	Subscribe    = 32
	Subscribed   = 33
	Unsubscribe  = 34
	Unsubscribed = 35
	Event        = 36
	Call         = 48
	Cancel       = 49
	Result       = 50
	Register     = 64
	Registred    = 65
	Unregister   = 66
	Unregistred  = 67
	Invocation   = 68
	Interrupt    = 69
	Yield        = 70
)

const RAND_MAX = 999999999

type (
	args     []interface{}
	kwargs   map[string]interface{}
	details  map[string]map[string]map[string]interface{}
	callback func(args []interface{}, kwargs map[string]interface{})
	proc     func([]interface{}, map[string]interface{}) ([]interface{}, map[string]interface{})
)
