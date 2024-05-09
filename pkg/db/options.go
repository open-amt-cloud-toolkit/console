package db

import "time"

// Option -.
type Option func(*SQL)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *SQL) {
		c.maxPoolSize = size
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *SQL) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *SQL) {
		c.connTimeout = timeout
	}
}
