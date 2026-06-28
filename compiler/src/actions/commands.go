package actions

// Purposely use an empty body as a noop so that you have to explicitly send a
// request body.
var NOOP = command{Action: 0o0}
var RECOMPILE = command{Action: 0o1}
