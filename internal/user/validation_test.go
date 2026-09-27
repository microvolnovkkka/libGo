package user

import (
	"errors"
	"strings"
	"testing"
)

func TestNormalizeUsername(t *testing.T) {
	// Каждый сценарий содержит входную строку и ожидаемый результат.
	tests := []struct {
		name     string
		username string
		want     string
	}{
		{name: "trim_and_lowercase", username: " Fedor_12 ", want: "fedor_12"},
		{name: "already_normalized", username: "fedor_12", want: "fedor_12"},
		{name: "tabs_and_newlines", username: "\tFedor\n", want: "fedor"},
		{name: "preserve_inner_space", username: " FE DOR ", want: "fe dor"},
		{name: "whitespace_only", username: " \t\n ", want: ""},
		{name: "empty", username: "", want: ""},
	}

	for _, tt := range tests {
		// t.Run показывает каждый сценарий отдельным подтестом.
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeUsername(tt.username)
			if got != tt.want {
				t.Errorf("NormalizeUsername(%q): получили %q, ожидали %q", tt.username, got, tt.want)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  error
	}{
		{name: "allowed_characters", username: "az09_", wantErr: nil},
		{name: "minimum_length", username: "abc", wantErr: nil},
		{name: "maximum_length", username: strings.Repeat("a", 32), wantErr: nil},
		{name: "empty", username: "", wantErr: ErrInvalidUsername},
		{name: "too_short", username: "ab", wantErr: ErrInvalidUsername},
		{name: "too_long", username: strings.Repeat("a", 33), wantErr: ErrInvalidUsername},
		{name: "uppercase", username: "Fedor", wantErr: ErrInvalidUsername},
		{name: "inner_space", username: "fe dor", wantErr: ErrInvalidUsername},
		{name: "outer_spaces", username: " fedor ", wantErr: ErrInvalidUsername},
		{name: "cyrillic", username: "федор", wantErr: ErrInvalidUsername},
		{name: "punctuation", username: "fedor!", wantErr: ErrInvalidUsername},
		{name: "invalid_utf8", username: "\xffab", wantErr: ErrInvalidUsername},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			// Проверяем нужную ошибку, а не её текст. errors.Is(nil, nil) тоже даёт true.
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateUsername(%q): получили ошибку %v, ожидали %v", tt.username, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{name: "empty", password: "", wantErr: ErrInvalidPassword},
		{name: "too_short", password: strings.Repeat("a", 7), wantErr: ErrInvalidPassword},
		{name: "minimum_length", password: strings.Repeat("a", 8), wantErr: nil},
		{name: "maximum_length", password: strings.Repeat("a", 128), wantErr: nil},
		{name: "too_long", password: strings.Repeat("a", 129), wantErr: ErrInvalidPassword},
		// Кириллица помогает заметить случайный подсчёт байт вместо рун.
		{name: "unicode_too_short", password: strings.Repeat("я", 7), wantErr: ErrInvalidPassword},
		{name: "unicode_minimum", password: strings.Repeat("я", 8), wantErr: nil},
		{name: "unicode_maximum", password: strings.Repeat("я", 128), wantErr: nil},
		{name: "unicode_too_long", password: strings.Repeat("я", 129), wantErr: ErrInvalidPassword},
		// Пробел входит в длину пароля: семь букв и пробел должны пройти проверку.
		{name: "leading_space_counts", password: " aaaaaaa", wantErr: nil},
		{name: "trailing_space_counts", password: "aaaaaaa ", wantErr: nil},
		{name: "invalid_utf8", password: "\xffabcdefgh", wantErr: ErrInvalidPassword},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidatePassword: получили ошибку %v, ожидали %v", err, tt.wantErr)
			}
		})
	}
}
