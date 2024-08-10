package gonuts

// Generic function to create a pointer to any type
func Ptr[T any](v T) *T {
	return &v
}

// Type-specific functions for common types

// StrPtr returns a pointer to the given string
func StrPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the given int
func IntPtr(i int) *int {
	return &i
}

// BoolPtr returns a pointer to the given bool
func BoolPtr(b bool) *bool {
	return &b
}

// Float64Ptr returns a pointer to the given float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// Float32Ptr returns a pointer to the given float32
func Float32Ptr(f float32) *float32 {
	return &f
}

// Int64Ptr returns a pointer to the given int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// Int32Ptr returns a pointer to the given int32
func Int32Ptr(i int32) *int32 {
	return &i
}

// Uint64Ptr returns a pointer to the given uint64
func Uint64Ptr(u uint64) *uint64 {
	return &u
}

// Uint32Ptr returns a pointer to the given uint32
func Uint32Ptr(u uint32) *uint32 {
	return &u
}
