package task

import "testing"

func TestError_StatusCode(t *testing.T) {
	tests := []struct {
		name    string
		subject *Error
		want    string
	}{
		{
			name: "no status code",
			subject: &Error{
				Status: "doesn't start with a number",
			},
			want: "",
		},
		{
			name: "starts with a status code",
			subject: &Error{
				Status: "418 I'm a teapot",
			},
			want: "418",
		},
		{
			name: "includes a number but doesn't start with it",
			subject: &Error{
				Status: "The number 42 should not be found",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.subject.StatusCode(); got != tt.want {
				t.Errorf("StatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
