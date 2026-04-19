// Package portrank assigns risk levels to open ports based on port number,
// protocol, and configurable per-port weights.
//
// Ports are scored and grouped into four levels: low, medium, high, and
// critical. The ranking can be used to prioritise alerts or filter output.
package portrank
