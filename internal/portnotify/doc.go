// Package portnotify evaluates a set of open ports against user-defined watch
// rules and emits structured events when a port matches.
//
// Rules are expressed as "port[/proto][=label]" strings, for example:
//
//	22/tcp=ssh-watch
//	3306=mysql-alert
//
// Events can be printed with [Print] / [Fprint] or consumed programmatically
// via [Notifier.Check].
package portnotify
