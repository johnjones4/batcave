package util

import "time"

func FormatTime(t time.Time) string {
	return t.Format("January 2, 2006 at 3:04pm")
}
