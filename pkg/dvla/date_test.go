package dvla

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestItUnmarshalsDate(t *testing.T) {
	raw := `"2024-11-23"`

	date := &Date{}
	require.Nil(t, date.UnmarshalJSON([]byte(raw)))
	require.Equal(t, 2024, time.Time(*date).Year())
	require.Equal(t, time.November, time.Time(*date).Month())
	require.Equal(t, 23, time.Time(*date).Day())
}
