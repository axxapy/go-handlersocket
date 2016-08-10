package handlersocket

import (
	"fmt"
	"strings"
	"io"
	"strconv"
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
	index_num   string
	assert_type string
	keys        []string
	values      []string
}

type cmd_delete struct {
	index_num   string
	assert_type string
	keys        []string
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
	limit := len(f.values)
	offset := 0
	where := []string{f.assert_type, strconv.Itoa(len(f.keys)), f.keys}
	if _, err := fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\t%s\n", f.index_num, strings.Join(where, "\t"), limit, offset, "U", strings.Join(f.values, "\t")); err != nil {
		return err
	}
	return nil
}
