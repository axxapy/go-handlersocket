package handlersocket

import (
	"fmt"
	"strings"
	"io"
)

type cmd_openindex struct {
	command string
	params  []string
}

type cmd_find struct {
	command string
	params  []string
	limit   int
	offset  int
}

type cmd_insert struct {
	command string
	params  []string
}

type cmd_update struct {
	command  string
	criteria []string
	limit    int
	offset   int
	mop      string
	newvals  []string
}

func (f *cmd_openindex) write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s\t%s\n", f.command, strings.Join(f.params, "\t")); err != nil {
		return err
	}
	return nil
}

func (f *cmd_find) write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", f.command, strings.Join(f.params, "\t"), f.limit, f.offset); err != nil {
		return err
	}
	return nil
}

func (f *cmd_insert) write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s\t%s\n", f.command, strings.Join(f.params, "\t")); err != nil {
		return err
	}
	return nil
}

func (f *cmd_update) write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\t%s\n", f.command, strings.Join(f.criteria, "\t"), f.limit, f.offset, f.mop, strings.Join(f.newvals, "\t")); err != nil {
		return err
	}
	return nil
}
