package passwords

import (
	"go.uber.org/zap"
	"time"

	scrypt "github.com/elithrar/simple-scrypt"
)

const (
	// DefaultMaxTimeout default max timeout in ms
	DefaultMaxTimeout = 500 * time.Millisecond

	// DefaultMaxMemory default max memory in MB
	DefaultMaxMemory = 64
)

var logger, _ = zap.NewProduction()

// Options ...
type Options struct {
	maxTimeout time.Duration
	maxMemory  int
}

// NewOptions ...
func NewOptions(maxTimeout time.Duration, maxMemory int) *Options {
	return &Options{maxTimeout, maxMemory}
}

// ScryptPasswords ...
type ScryptPasswords struct {
	options *Options
	params  scrypt.Params
}

// NewScryptPasswords ...
func NewScryptPasswords(options *Options) Passwords {
	if options == nil {
		options = &Options{}
	}

	if options.maxTimeout == 0 {
		options.maxTimeout = DefaultMaxTimeout
	}
	if options.maxMemory == 0 {
		options.maxMemory = DefaultMaxMemory
	}

	params, err := scrypt.Calibrate(
		options.maxTimeout,
		options.maxMemory,
		scrypt.DefaultParams,
	)
	if err != nil {
		logger.Error("error calibrating scrypt params")
	}

	return &ScryptPasswords{options, params}
}

// CreatePassword ...
func (sp *ScryptPasswords) CreatePassword(password string) (string, error) {
	hash, err := scrypt.GenerateFromPassword([]byte(password), sp.params)
	return string(hash), err
}

// CheckPassword ...
func (sp *ScryptPasswords) CheckPassword(hash, password string) error {
	return scrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
