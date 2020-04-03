package ratelimiter

import "time"

// New ratelimiter module
func New(rules []RateLimiterRule) Ratelimiter {
	module := &rlModule{
		rules: rules,
	}

	return module
}

// IsRateLimit ...
func (r *rlModule) IsRateLimit(requestTime time.Time) (bool, error) {
	// Check all rules
	for rule := range r.rules {

		// Check the rule
		isRateLimit, err := r.rules[rule].IsRateLimit(requestTime)

		// Failed to check
		if err != nil {
			return false, err
		}

		// Positive should ratelimit
		if isRateLimit {
			// Send log TODO
			// Either use datadog or iris

			return true, nil
		}
	}
	return false, nil
}
