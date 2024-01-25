// Package readliner implement an io.Reader that reads from a tty with history and completions.
package readliner

import (
	"os"
	"strings"

	"github.com/peterh/liner"
)

// ReadLiner is an io.Reader that can read from a tty using `readline`
type ReadLiner struct {
	liner       *liner.State
	completions []string
	history     string
	prompt      string
	eol         string
	buf         []byte
	err         error
}

const DefaultEOL = "\r\n"

// New creates a new ReadLiner and sets the tty in raw mode.
//
// `prompt` is printed before reading from each line.
// `history` should be the path to the history file.
func New(prompt, history string) *ReadLiner {
	rl := &ReadLiner{liner: liner.NewLiner(), history: history, prompt: prompt, eol: DefaultEOL}
	rl.liner.SetCtrlCAborts(true)

	if history != "" {
		if f, err := os.Open(history); err == nil {
			defer f.Close()
			rl.liner.ReadHistory(f)
		}
	}

	return rl
}

// SetPrompt changes the `prompt` for the ReadLiner
func (r *ReadLiner) SetPrompt(prompt string) {
	r.prompt = prompt
}

// SetCompletions sets a `completer` with a list of completions words.
//
// If `begin` is true, words are completed only at the beginning of the line (i.e. command names).
// If false, the last word of the line is completed.
func (r *ReadLiner) SetCompletions(completions []string, begin bool) {
	r.completions = completions

	if completions != nil {
		r.liner.SetCompleter(func(line string) (c []string) {
			prefix := ""

			if !begin {
				if i := strings.LastIndexAny(line, " \t'!@#$%^&*()-_=+[]{:\";'}|\\,./<>"); i >= 0 {
					prefix = line[:i+1]
					line = line[i+1:]
				}
			}
			for _, n := range r.completions {
				if strings.HasPrefix(n, strings.ToLower(line)) {
					c = append(c, prefix+n)
				}
			}
			return
		})
	} else {
		r.liner.SetCompleter(nil)
	}
}

func (r *ReadLiner) SetEol(eol string) {
	r.eol = eol
}

// Close closes the ReadLiner and reset the TTY. If there is an history file, the current history
// is written to the file.
func (r *ReadLiner) Close() error {
	defer r.liner.Close()

	if r.history != "" {
		if f, err := os.Create(r.history); err != nil {
			//fmt.Println(err)
		} else {
			defer f.Close()
			r.liner.WriteHistory(f)
		}
	}

	return nil
}

// Read implements io.Reader.Read
func (r *ReadLiner) Read(b []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}

	l := len(r.buf)

	if l == 0 {
		line, err := r.liner.Prompt(r.prompt)
		if err != nil {
			r.err = err
			return 0, err
		}

		r.liner.AppendHistory(line)
		r.buf = []byte(line + r.eol)
	}

	n := 0

	if l > 0 {
		n = copy(b, r.buf)
		r.buf = r.buf[n:]
	}

	return n, nil
}
