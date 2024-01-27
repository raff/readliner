# readliner
An io.Reader that uses `github.com/peterh/liner` to read from a tty (with command completion and editing).

This is useful for cases where you are trying to create an interactive cli (i.e. REPL) but the input is being
processed by a separate entity (like an expression parser) that operates on an io.Reader.
