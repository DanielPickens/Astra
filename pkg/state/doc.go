// Package state gives access to the state of the astra process stored in a local file
// The state of an instance is stored in a file .astra/devstate.${PID}.json.
// For compatibility with previous versions of astra, the `devstate.json` file contains
// the state of the first instance of astra.
package state
