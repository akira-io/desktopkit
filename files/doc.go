// Package files performs filesystem actions that are visible to the user:
// opening a path with the default application, opening a URL in the default
// browser, and revealing a path inside the platform's file manager.
//
// It does not read, write, or copy files. Use the standard library os and io
// packages for those concerns.
package files
