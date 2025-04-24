package email

import (
	"bytes"
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed logo_test.svg
var files embed.FS

func TestRawEmail(t *testing.T) {
	raw, err := rawEmail(
		"to@example.com",
		"from@example.com",
		"test subject",
		`<h4>body</h4><img src="cid:logo"><p>End of body</p>`,
		map[string]string{"logo": "logo_test.svg"},
		&files)
	require.NoError(t, err)

	require.Greater(t, len(raw), 1000)
}

func Test_encodeFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
		want     string
	}{
		{
			name:     "file not found",
			filename: "gone.txt",
			wantErr:  true,
			want:     "",
		},
		{
			name:     "good",
			filename: "logo_test.svg",
			wantErr:  false,
			want:     "PHN2ZyB3aWR0aD0iMTAiIGhlaWdodD0iMTAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHJlY3Qgd2lkdGg9IjEwIiBoZWlnaHQ9IjEwIiBzdHlsZT0iZmlsbDpyZWQiLz48L3N2Zz4K",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)
			buf := bytes.Buffer{}
			err := encodeFile(&files, tt.filename, &buf)
			if tt.wantErr {
				assert.Error(err)
				return
			}
			assert.NoError(err)

			got := buf.String()
			assert.Equal(tt.want, got)
		})
	}
}
