package handlersocket

import (
	"testing"
	"bytes"
)

func test_write(t *testing.T, cmd hs_chan_writer, want string) {
	var b bytes.Buffer;

	err := cmd.write(&b);
	if err != nil {
		t.Fatal(err);
	}

	got := b.String();
	if got != want {
		t.Errorf("Expect '%s' but got '%s'", want, got)
	}
}

func TestOpenIndexFormatting(t *testing.T) {
	spec := &hs_index_spec{db_name: "db_name", table_name: "table_name", index_name: "index_name", columns: []string{"col1", "col2"}}
	cmd := &cmd_openindex{index_num: "10", index_spec: spec}

	test_write(t, cmd, "P	10	db_name	table_name	index_name	col1,col2\n");
}

func TestFindFormatting(t *testing.T) {
	cmd := &cmd_find{command: "10", params: []string{"=", "2", "one	two"}, limit: 100, offset: 2}
	test_write(t, cmd, "10	=	2	one	two	100	2\n")
}

func TestInsertFormatting(t *testing.T) {
	cmd := &cmd_insert{command: "10", params: []string{"+", "2", "one	two"}}
	test_write(t, cmd, "10	+	2	one	two\n")
}

func TestUpdateFormatting(t *testing.T) {
	cmd := &cmd_update{
		index_num:   "10",
		assert_type: "=",
		keys:        []string{"key1", "key2"},
		values:      []string{"value1", "value2"},
	}
	test_write(t, cmd, "10	=	2	key1	key2	2	0	U	value1	value2\n")
}

func TestDeleteFormattingg(t *testing.T) {
	cmd := &cmd_delete{
		index_num:   "10",
		assert_type: "=",
		keys:        []string{"key1", "key2"},
	}
	test_write(t, cmd, "10	=	2	key1	key2	0	0	D\n")
}