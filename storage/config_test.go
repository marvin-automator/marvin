package storage

import (
	"testing"
	"github.com/bigblind/marvin/domain"
)

func TestGetConfigNoSaved(t *testing.T) {
	WithTestDB(t, func(s Store) {
		c, err := s.GetConfig()
		if err != nil {
			t.Error(err)
		}

		if c != domain.DefaultConfig {
			t.Errorf("Should have gotten the default config, but didn't.\n%v\n%v", c, domain.DefaultConfig)
		}
	})
}


func TestSaveAndGetConfig(t *testing.T) {
	WithTestDB(t, func(s Store) {
		c1 := domain.DefaultConfig
		c1.AccountsEnabled = true

		err := s.SaveConfig(c1)
		if err != nil {
			t.Error(err)
		}

		c2, err := s.GetConfig()
		if err != nil {
			t.Error(err)
		}

		if c2 != c1 {
			t.Errorf("Saved and loaded config aren't equal\nSaved:  %v\nloaded: %v", c1, c2)
		}
	})
}
