package scanner

// CommonPorts is a curated list of well-known ports to scan by default.
var CommonPorts = []int{
	21,   // FTP
	22,   // SSH
	23,   // Telnet
	25,   // SMTP
	53,   // DNS
	80,   // HTTP
	110,  // POP3
	143,  // IMAP
	443,  // HTTPS
	465,  // SMTPS
	587,  // SMTP submission
	993,  // IMAPS
	995,  // POP3S
	3306, // MySQL
	5432, // PostgreSQL
	6379, // Redis
	8080, // HTTP alt
	8443, // HTTPS alt
	27017, // MongoDB
}

// WellKnownService maps common port numbers to service names.
var WellKnownService = map[int]string{
	21:    "ftp",
	22:    "ssh",
	23:    "telnet",
	25:    "smtp",
	53:    "dns",
	80:    "http",
	110:   "pop3",
	143:   "imap",
	443:   "https",
	465:   "smtps",
	587:   "smtp-submission",
	993:   "imaps",
	995:   "pop3s",
	3306:  "mysql",
	5432:  "postgresql",
	6379:  "redis",
	8080:  "http-alt",
	8443:  "https-alt",
	27017: "mongodb",
}

// LookupService returns the service name for a given port, or "unknown".
func LookupService(port int) string {
	if name, ok := WellKnownService[port]; ok {
		return name
	}
	return "unknown"
}
