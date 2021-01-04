package gcpcomputetimer

import (
	"fmt"
)

// SecondsToHuman translate seconds into a human readable string with a
// specific length.
//
//   fmt.Printf("time: %s\n", SecondsToHuman(500))
//
func SecondsToHuman(seconds int64) string {
	var retv string

	if seconds == 0 {
		return "                 "
	}

	if seconds > 86400 {
		retv = retv + fmt.Sprintf("%3d days ", seconds/86400)
		seconds = seconds % 86400
	} else {
		retv = retv + "         "
	}

	if seconds > 3600 {
		retv = retv + fmt.Sprintf("%02d:", seconds/3600)
		seconds = seconds % 3600
	} else {
		retv = retv + "00:"
	}

	if seconds > 60 {
		retv = retv + fmt.Sprintf("%02d:", seconds/60)
		seconds = seconds % 60
	} else {
		retv = retv + "00:"
	}

	retv = retv + fmt.Sprintf("%02d", seconds)

	return retv
}
