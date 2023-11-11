package mapexamples

var bob = &User{name: "Bob", age: 26}

// see the User structure above
var userRefs = map[string][]*User{
	"admin":   {bob},
	"manager": {bob, {name: "Alice", age: 35}},
	"":        {{name: "Tom", age: 18}},
	"cto":     {{name: "Chris", age: 41}},
}
