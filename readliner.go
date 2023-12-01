package readliner

import (
	"os"
	"strings"

	"github.com/peterh/liner"
)

type ReadLiner struct {
	liner       *liner.State
	completions []string
	history     string
	prompt      string
	eol         string
	buf         []byte
	err         error
}

func New(prompt, history string) *ReadLiner {
	rl := &ReadLiner{liner: liner.NewLiner(), history: history, prompt: prompt, eol: "\n"}
	rl.liner.SetCtrlCAborts(true)

	if history != "" {
		if f, err := os.Open(history); err == nil {
			defer f.Close()
			rl.liner.ReadHistory(f)
		}
	}

	return rl
}

func (r *ReadLiner) SetPrompt(prompt string) {
	r.prompt = prompt
}

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
		r.buf = []byte(line + "\n\n")
	}

	n := 0

	if l > 0 {
		n = copy(b, r.buf)
		r.buf = r.buf[n:]
	}

	return n, nil
}
