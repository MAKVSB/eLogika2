package utils

import (
	"net"
	"regexp"
	"strings"
)

// ipInRange checks if ip is within [start, end] range
func ipInRange(ip, start, end net.IP) bool {
	// Compare byte slices
	if len(ip) != len(start) {
		return false
	}
	for i := 0; i < len(ip); i++ {
		if ip[i] < start[i] || ip[i] > end[i] {
			return false
		}
	}
	return true
}

// IsIPAllowed checks if currentIP matches any of the rules in allowedList
func IsIPAllowed(allowedList string, currentIP string) bool {
	ip := net.ParseIP(currentIP)
	if ip == nil {
		return false
	}

	rules := strings.Split(allowedList, ";")
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		// Case 1: exact match (single IP)
		if !strings.Contains(rule, "-") {
			if ip.Equal(net.ParseIP(rule)) {
				return true
			}
			continue
		}

		// Case 2: range (start-end)
		parts := strings.Split(rule, "-")
		if len(parts) != 2 {
			continue
		}
		start := net.ParseIP(strings.TrimSpace(parts[0]))
		end := net.ParseIP(strings.TrimSpace(parts[1]))
		if start == nil || end == nil {
			continue
		}

		if ipInRange(ip, start, end) {
			return true
		}
	}

	return false
}

// IsValidIPCondition checks if the condition string has valid format:
// single IPs (x.x.x.x) or ranges (x.x.x.x-x.x.x.x) separated by ;
func IsValidIPCondition(cond string) bool {
	pattern := `^(?:\d{1,3}(?:\.\d{1,3}){3}(?:-\d{1,3}(?:\.\d{1,3}){3})?)(?:;(?:\d{1,3}(?:\.\d{1,3}){3}(?:-\d{1,3}(?:\.\d{1,3}){3})?))*$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(cond)
}
