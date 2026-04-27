// Package portjournal provides a persistent, append-only log of port
// lifecycle events observed by portwatch.
//
// Each entry captures the time, port number, protocol, optional service
// name, event kind (opened / closed / changed), and a free-text note.
//
// Entries are stored as a JSON array on disk and can be queried or
// printed via the format helpers in this package.
package portjournal
