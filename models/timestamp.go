package models

import (
	"fmt"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Timestamp data as unix time in seconds
type Timestamp time.Time

// MarshalJSON defines how encoding/json marshals the object to JSON,
// the result is a string of the UNIX timestamp
func (t Timestamp) MarshalJSON() ([]byte, error) {
	ts := t.Time().Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

// UnmarshalJSON defines how encoding/json unmarshals the object from JSON,
// a UNIX timestamp string is converted to int which is used for the Timestamp
// object value
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(int64(ts), 0))

	return nil
}

// GetBSON defines how labix.org/v2/mgo marshals the object to BSON,
// the result is a time.Time object which is then handled by mgo
func (t Timestamp) GetBSON() (interface{}, error) {
	if t.Time().IsZero() {
		return nil, nil
	}

	return t.Time(), nil
}

// SetBSON defines how labix.org/v2/mgo unmarshals the object from BSON,
// the raw BSON data is unmarshaled to a time.Time object which is used for the
// Timestamp object value
func (t *Timestamp) SetBSON(raw bson.Raw) error {
	var tm time.Time

	if err := raw.Unmarshal(&tm); err != nil {
		return err
	}

	*t = Timestamp(tm)

	return nil
}

// String returns the string representation of the Timestamp object,
// it is equal to the time.Time string representation of the Timestamp object
// value
func (t Timestamp) String() string {
	return t.Time().String()
}

// Time returns a time.Time object with the same time value as the Timestamp
// object
func (t Timestamp) Time() time.Time {
	if time.Time(t).IsZero() {
		return time.Unix(0, 0)
	}

	return time.Time(t)
}

// Now returns a pointer to a Timestamp object with the current time,
// it is equal to creating a Timestamp object from time.Now()
func Now() *Timestamp {
	t := Timestamp(time.Now())
	return &t
}

// Unix calls the Unix() method of a time.Time object with the same time values
// as the timestamp object
func (t Timestamp) Unix() int64 {
	return t.Time().Unix()
}

// Time returns a pointer to a Timestamp object which is created
// from a time.Time object
func Time(t time.Time) *Timestamp {
	ts := Timestamp(t)
	return &ts
}

// Unix returns a pointer to a Timestamp object which is created from
// a UNIX timestamp
func Unix(sec, nsec int64) *Timestamp {
	t := time.Unix(sec, nsec)
	return Time(t)
}
