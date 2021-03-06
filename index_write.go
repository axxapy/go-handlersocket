package handlersocket

import (
	"strings"
	"strconv"
	"errors"
)

type WriteIndex struct {
	base_index
}

func (this *WriteIndex) Delete(assert_type string, keys map[string]interface{}) (modified_rows int, err error) {
	conn, err := this.conn_pool.getConnection()
	if err != nil {
		return 0, err
	}
	defer this.conn_pool.releaseConnection(conn)

	if err := this.open(conn); err != nil {
		return 0, err
	}

	conn.chan_write <- &cmd_delete{
		index_num:   conn.getIndexNum(this.spec),
		assert_type: assert_type,
		keys:        this.sortFieldsList(keys),
	}
	message := <-conn.chan_read

	if message.ReturnCode == "1" {
		return 0, errors.New("Error")
	}

	modified, err := strconv.Atoi(strings.TrimSpace(message.Data[1]))
	if err != nil {
		return 0, err
	}

	return modified, nil
}

func (this *WriteIndex) Update(assert_type string, keys map[string]interface{}, values map[string]interface{}) (modified_rows int, err error) {
	conn, err := this.conn_pool.getConnection()
	if err != nil {
		return 0, err
	}
	defer this.conn_pool.releaseConnection(conn)

	if err := this.open(conn); err != nil {
		return 0, err
	}

	conn.chan_write <- &cmd_update{
		index_num:   conn.getIndexNum(this.spec),
		assert_type: assert_type,
		keys:        this.sortFieldsList(keys),
		values:      this.sortFieldsList(values),
	}
	message := <-conn.chan_read

	if message.ReturnCode == "1" {
		return 0, errors.New("Error")
	}

	modified, err := strconv.Atoi(strings.TrimSpace(message.Data[1]))
	if err != nil {
		return 0, err
	}

	return modified, nil
}

func (this *WriteIndex) Insert(vals ...string) (err error) {
	conn, err := this.conn_pool.getConnection()
	if err != nil {
		return err
	}
	defer this.conn_pool.releaseConnection(conn)

	if err := this.open(conn); err != nil {
		return err
	}

	index_num := conn.getIndexNum(this.spec)

	values := strings.Join(vals, "\t")
	values_count := strconv.Itoa(len(vals))
	data := []string{"+", values_count, values}

	conn.chan_write <- &cmd_insert{command: index_num, params: data}
	message := <-conn.chan_read

	if message.ReturnCode == "1" {
		return errors.New("INSERT: Data Exists")
	}

	if message.ReturnCode != "0" {
		return errors.New("Error Inserting Data")
	}

	return nil
}
