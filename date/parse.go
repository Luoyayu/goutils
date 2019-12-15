package date

import (
	"fmt"
	"time"
)

// ParseDate :layout default is RFC3339
func ParseDate(l string, v string) string {
	if l == "" {
		l = time.RFC3339
	}
	date_, _ := time.Parse(l, v)
	return fmt.Sprintf("%4d-%02d-%02d %02d:%02d", date_.Year(), int(date_.Month()), date_.Day(), date_.Hour(), date_.Minute())

}
