// Code generated by "stringer -type Alert -output types_string.go"; DO NOT EDIT.

package function

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnknownAlert-0]
	_ = x[LONG-1]
	_ = x[REDUCE-2]
	_ = x[CLOSE-3]
	_ = x[STOP_LOSS-4]
}

const _Alert_name = "UnknownAlertLONGREDUCECLOSESTOP_LOSS"

var _Alert_index = [...]uint8{0, 12, 16, 22, 27, 36}

func (i Alert) String() string {
	if i < 0 || i >= Alert(len(_Alert_index)-1) {
		return "Alert(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Alert_name[_Alert_index[i]:_Alert_index[i+1]]
}
