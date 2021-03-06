// Code generated by "stringer -type=HTTPMethod"; DO NOT EDIT.

package webutil

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[GET-0]
	_ = x[HEAD-1]
	_ = x[POST-2]
	_ = x[PUT-3]
	_ = x[CONNECT-4]
	_ = x[DELETE-5]
	_ = x[OPTIONS-6]
	_ = x[PATCH-7]
	_ = x[TRACE-8]
}

const _HTTPMethod_name = "GETHEADPOSTPUTCONNECTDELETEOPTIONSPATCHTRACE"

var _HTTPMethod_index = [...]uint8{0, 3, 7, 11, 14, 21, 27, 34, 39, 44}

func (i HTTPMethod) String() string {
	if i < 0 || i >= HTTPMethod(len(_HTTPMethod_index)-1) {
		return "HTTPMethod(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _HTTPMethod_name[_HTTPMethod_index[i]:_HTTPMethod_index[i+1]]
}
