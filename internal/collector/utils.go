package collector

import "strings"

// ProcessFQDN removes the specified domain suffix from an FQDN if present
func ProcessFQDN(fqdn, domainSuffix string) string {
	if domainSuffix == "" {
		return fqdn
	}

	if strings.HasSuffix(fqdn, domainSuffix) {
		return fqdn[:len(fqdn)-len(domainSuffix)]
	}

	return fqdn
}

// CleanIPAddress removes CIDR notation if present
func CleanIPAddress(ip string) string {
	if idx := strings.Index(ip, "/"); idx != -1 {
		return ip[:idx]
	}
	return ip
}
