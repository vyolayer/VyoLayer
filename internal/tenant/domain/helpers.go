package domain

// import "time"

// // NullableTime safely wraps time.Time to handle nil pointers and zero values
// // specifically for mapping to external layers like Protobuf or JSON.
// type NullableTime struct {
// 	time time.Time
// }

// // NewNullableTime creates a wrapper from a standard time value.
// func NewNullableTime(t time.Time) *NullableTime {
// 	return &NullableTime{time: t}
// }

// // NewNullableTimeFromPtr safely handles incoming *time.Time (like from GORM models)
// func NewNullableTimeFromPtr(t *time.Time) *NullableTime {
// 	if t == nil {
// 		return nil
// 	}
// 	return &NullableTime{time: *t}
// }

// // Time returns the underlying time.Time. If nil, returns a zero time.
// func (t *NullableTime) Time() time.Time {
// 	if t == nil {
// 		return time.Time{}
// 	}
// 	return t.time
// }

// // Valid returns true if the pointer is not nil AND the time is not zero.
// func (t *NullableTime) Valid() bool {
// 	if t == nil {
// 		return false
// 	}
// 	return !t.time.IsZero()
// }

// // Format safely formats the time. Returns an empty string if nil or zero.
// func (t *NullableTime) Format(layout string) string {
// 	if !t.Valid() {
// 		return ""
// 	}
// 	return t.time.Format(layout)
// }

// // GetTime is a convenience method for RFC3339. Returns "" if invalid.
// func (t *NullableTime) GetTime() string {
// 	return t.Format(time.RFC3339)
// }

// // Ptr returns a true pointer to the underlying time, or nil if invalid.
// func (t *NullableTime) Ptr() *time.Time {
// 	if !t.Valid() {
// 		return nil
// 	}
// 	// Return a copy's address to prevent external mutation
// 	timeCopy := t.time
// 	return &timeCopy
// }

// // PtrString returns a pointer to an RFC3339 string, or nil if invalid.
// func (t *NullableTime) PtrString() *string {
// 	if !t.Valid() {
// 		return nil
// 	}
// 	s := t.time.Format(time.RFC3339)
// 	return &s
// }
