// +build !debug

package tomato

import "time"

const (
	taskDuration     = 25 * time.Minute
	restDuration     = 5 * time.Minute
	longRestDuration = 15 * time.Minute
)
