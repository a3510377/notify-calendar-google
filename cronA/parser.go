package cronA

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ＊  ＊  ＊  ＊　＊  ＊
// ┬  ┬  ┬  ┬  ┬  ┬
// │  │  │  │  │  └ day of week (0 - 7, 1L - 7L) (0 or 7 is Sun)
// │  │  │  │  └── month (1 - 12)
// │  │  │  └──── day of month (1 - 31, L)
// │  │  └────── hour (0 - 23)
// │  └──────── minute (0 - 59)
// └────────── second (0 - 59, optional)

var (
	CronExpressionMap = []string{"second", "minute", "hour", "dayOfMonth", "month", "dayOfWeek"}
	defaults          = []string{"0", "0", "*", "*", "*", "*"}
	// https://man.freebsd.org/cgi/man.cgi?crontab%285%29
	predefined = map[string]string{
		"@yearly": "0 0 1 1 *", "@annually": "0 0 1 1 *",
		"@monthly": "0 0 1 * *",
		"@weekly":  "0 0 * * 0",
		"@daily":   "0 0 * * *", "@midnight": "0 0 * * *",
		"@hourly":       "0 * * * *",
		"@every_minute": "*/1 * * * *",
		"@every_second": "* * * * *",
		// "@reboot": "-", // cron is no reboot
	}
)

type (
	CronExpression struct{}
)

func (c *CronExpression) Parse(spec string) (s *Schedule, err error) {
	if len(spec) == 0 {
		return nil, fmt.Errorf("empty spec string")
	}

	loc := time.Local
	if strings.HasPrefix(spec, "TZ=") || strings.HasPrefix(spec, "CRON_TZ=") {
		i := strings.Index(spec, " ")
		if loc, err = time.LoadLocation(spec[strings.Index(spec, "=")+1 : i]); err != nil {
			return nil, err
		}
		spec = strings.TrimSpace(spec[i:])
	}

	if fined, ok := predefined[spec]; ok { // set from predefined
		spec = fined
	}

	fields := strings.Fields(spec)
	if fl := len(fields); fl < 6 {
		fields = append(defaults[:len(defaults)-len(fields)], fields...)
	} else if fl != 6 { // fl > 6
		return nil, fmt.Errorf("unexpected number of fields: %d", fl)
	}

	field := func(field string, r bounds) uint64 {
		if err != nil {
			return 0
		}
		var bits uint64
		bits, err = parseField(field, r)
		return bits
	}

	var (
		second     = field(fields[0], seconds)
		minute     = field(fields[1], minutes)
		hour       = field(fields[2], hours)
		dayOfMonth = field(fields[3], dom)
		month      = field(fields[4], months)
		dayOfWeek  = field(fields[5], dow)
	)

	if err != nil {
		return nil, err
	}

	fmt.Println(strings.Join(fields, " "))
	fmt.Println(loc.String())

	return &Schedule{second, minute, hour, dayOfMonth, month, dayOfWeek, loc}, nil
}

func parseField(field string, r bounds) (uint64, error) {
	var bits uint64

	for _, expr := range strings.Split(field, ",") {
		bit, err := parseRange(expr, r)
		if err != nil {
			return 0, err
		}
		bits |= bit
	}

	return bits, nil
}

func parseInt(expr string) (uint64, error) {
	num, err := strconv.Atoi(expr)
	if err != nil {
		return 0, err
	}
	if num < 0 {
		return 0, fmt.Errorf("negative number (%d) not allowed: %s", num, expr)
	}

	return uint64(num), nil
}

func parseFromIntOrName(expr string, r bounds) (uint64, error) {
	if u, ok := r.names[expr]; ok {
		return u, nil
	}

	num, err := parseInt(expr)
	if err != nil {
		return 0, err
	}

	if num < r.min {
		return 0, fmt.Errorf("not less than %d (%d)", r.min, num)
	}
	if num > r.max {
		return 0, fmt.Errorf("not greater than %d (%d)", r.min, num)
	}

	return 0, nil
}

func parseRange(expr string, r bounds) (u uint64, err error) {
	if expr == "*" || expr == "?" { // all values || no specific value
		return r.all(), nil
	}

	start, end := r.min, r.max
	// TODO: this
	// start, err = parseFromIntOrName(lowAndHigh[0], r)
	switch lowAndHigh := strings.Split(expr, "-"); len(lowAndHigh) {
	case 1:
		end = start
	case 2:
		end, err = parseFromIntOrName(lowAndHigh[0], r)
	default:
		err = fmt.Errorf("too many hyphens: %s", expr)
	}
	if err != nil {
		return 0, err
	}

	var step uint64
	switch rangeAndStep := strings.Split(expr, "/"); len(rangeAndStep) {
	case 1:
		step = 1
	case 2:
		if step, err = parseInt(rangeAndStep[1]); step == 0 {
			err = fmt.Errorf("step of range should be a positive number: %s", expr)
		}
	default:
		err = fmt.Errorf("too many slashes: %s", expr)
	}
	if err != nil {
		return 0, err
	}

	// if start < r.min {
	// 	return 0, fmt.Errorf("beginning of range (%d) below minimum (%d): %s", start, r.min, expr)
	// }
	// if end > r.max {
	// 	return 0, fmt.Errorf("end of range (%d) above maximum (%d): %s", end, r.max, expr)
	// }
	// if start > end {
	// 	return 0, fmt.Errorf("beginning of range (%d) beyond end of range (%d): %s", start, end, expr)
	// }

	fmt.Println(start, end)
	return getBits(start, end, step), nil
}

func NewCronExpression() *CronExpression { return &CronExpression{} }
