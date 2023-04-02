package store

import (
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	testData := map[string]struct {
		cfg *mysql.Config
		err error
	}{
		"valid": {
			TestConfig(),
			nil,
		},
		"no config": {
			nil,
			ErrNoConfig,
		},
	}

	for name, td := range testData {
		t.Run(name, func(t *testing.T) {
			_, err := NewClient(td.cfg)
			if err != nil {
				require.ErrorIs(t, err, td.err)
				return
			}
			require.Nil(t, td.err)
		})
	}
}
