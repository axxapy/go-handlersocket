package handlersocket

import (
	"fmt"
	"strings"
	"io"
)

type hs_cmd_openindex struct {
	command string
	params  []string
}

type hs_cmd_find struct {
	command string
	params  []string
	limit   int
	offset  int
}

func (f *hs_cmd_openindex) write(w io.Writer) error {

	if _, err := fmt.Fprintf(w, "%s\t%s\n", f.command, strings.Join(f.params, "\t")); err != nil {

		return err
	}

	return nil
}

func (f *hs_cmd_find) write(w io.Writer) error {

	if _, err := fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", f.command, strings.Join(f.params, "\t"), f.limit, f.offset); err != nil {
		return err
	}

	return nil
}
