package kemono

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		d       Time
		args    args
		want    time.Time
		wantErr require.ErrorAssertionFunc
	}{
		{"string", Time{}, args{[]byte(`"2024-05-02T14:21:02.807702"`)}, time.Date(2024, time.May, 2, 14, 21, 2, 807702000, time.UTC), require.NoError},
		{"unix", Time{}, args{[]byte(`1714659663`)}, time.Date(2024, time.May, 2, 14, 21, 3, 0, time.UTC), require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, tt.d.UnmarshalJSON(tt.args.bytes))
			assert.Equal(t, tt.want, time.Time(tt.d))
		})
	}
}
