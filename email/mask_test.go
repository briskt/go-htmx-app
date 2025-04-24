package email

import (
	"testing"
)

func Test_MaskEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  string
	}{
		{
			name:  "very short name",
			email: "x@example.org",
			want:  "*@e******.o**",
		},
		{
			name:  "single domain",
			email: "john_doe@example.org",
			want:  "j***_d**@e******.o**",
		},
		{
			name:  "subdomain",
			email: "john.doe@email.example.org",
			want:  "j***.d**@e****.e******.o**",
		},
		{
			name:  "adjacent separators",
			email: "john..doe@email._example.org",
			want:  "j***.d**@e****.e******.o**",
		},
		{
			name:  "separator and start or end",
			email: "_john_doe@email.example.org.",
			want:  "j***_d**@e****.e******.o**",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskEmail(tt.email); got != tt.want {
				t.Errorf("MaskEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
