package server_test

import (
	"errors"
	"testing"

	"github.com/jsnfwlr/filamate/internal/server"
)

func TestConfigLoad(t *testing.T) {
	testCases := []struct {
		name          string
		env           map[string]string
		expectedErr   error
		expStaticType string
	}{
		{
			name: "embedded",
			env: map[string]string{
				"STATIC_TYPE": "embedded",
			},
			expStaticType: "embedded",
		},
		{
			name: "directory",
			env: map[string]string{
				"STATIC_TYPE": "directory",
			},
			expStaticType: "directory",
		},
		{
			name: "other",
			env: map[string]string{
				"STATIC_TYPE": "other",
			},
			expectedErr: errors.New("config validation failed: Key: 'EnvConfig.StaticType' Error:Field validation for 'StaticType' failed on the 'oneof' tag"),
		},
		{
			name:          "blank",
			env:           map[string]string{},
			expStaticType: "embedded",
		},
		{
			name: "empty_string",
			env: map[string]string{
				"STATIC_TYPE": "",
			},
			expStaticType: "embedded",
		},
		{
			name: "case_sensitive_invalid",
			env: map[string]string{
				"STATIC_TYPE": "EMBEDDED",
			},
			expectedErr: errors.New("config validation failed: Key: 'EnvConfig.StaticType' Error:Field validation for 'StaticType' failed on the 'oneof' tag"),
		},
		{
			name: "whitespace",
			env: map[string]string{
				"STATIC_TYPE": " embedded ",
			},
			expectedErr: errors.New("config validation failed: Key: 'EnvConfig.StaticType' Error:Field validation for 'StaticType' failed on the 'oneof' tag"),
		},
		{
			name: "numeric_value",
			env: map[string]string{
				"STATIC_TYPE": "123",
			},
			expectedErr: errors.New("config validation failed: Key: 'EnvConfig.StaticType' Error:Field validation for 'StaticType' failed on the 'oneof' tag"),
		},
		{
			name: "special_characters",
			env: map[string]string{
				"STATIC_TYPE": "embedded!@#",
			},
			expectedErr: errors.New("config validation failed: Key: 'EnvConfig.StaticType' Error:Field validation for 'StaticType' failed on the 'oneof' tag"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			cfg, err := server.LoadConfig()
			switch {
			case tc.expectedErr == nil && err != nil:
				t.Errorf("unexpected error: %v", err)
				return
			case tc.expectedErr != nil && err == nil:
				t.Errorf("missing error %v", tc.expectedErr)
				return
			case tc.expectedErr != nil && err != nil && err.Error() != tc.expectedErr.Error():
				t.Errorf("incorrect error:\nexpected:\n\t%v\nreceived:\n\t%v", tc.expectedErr, err)
				return
			}

			if cfg.StaticType() != tc.expStaticType {
				t.Errorf("incorrect config:\nexpected:\n\t%v\nreceived:\n\t%v", tc.expStaticType, cfg.StaticType())
			}
		})
	}
}
