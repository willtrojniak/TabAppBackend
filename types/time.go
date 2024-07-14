package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
	"strings"
	"time"
)

type Time struct {
	time.Duration
}

func (t *Time) ScanTime(v pgtype.Time) error {
	t.Duration = time.Microsecond * time.Duration(v.Microseconds)
	return nil
}

func (t Time) TimeValue() (pgtype.Time, error) {
	return pgtype.Time{Microseconds: t.Duration.Microseconds(), Valid: true}, nil
}

func (d *Time) String() string {
	return fmt.Sprintf("%02d:%02d", int(d.Hours()), int(d.Minutes())%60)
}

func (d *Time) UnmarshalJSON(b []byte) (err error) {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		segments := strings.Split(value, ":")
		if len(segments) != 2 {
			return errors.New("Invalid time (number of time segments)")
		}
		hours, err := strconv.Atoi(segments[0])
		if err != nil || hours < 0 || hours > 23 {
			return errors.New("Invalid time (hour)")
		}
		minutes, err := strconv.Atoi(segments[1])
		if err != nil || minutes < 0 || minutes > 59 {
			return errors.New("Invalid time (minute)")
		}

		d.Duration, err = time.ParseDuration(fmt.Sprintf("%vh%vm", hours, minutes))
		if err != nil {
			return errors.New("Invalid time")
		}
		return nil
	default:
		return errors.New("Invalid time")
	}
}

func (t Time) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
}
