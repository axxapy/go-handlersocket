package handlersocket

import (
	"strings"
	"strconv"
	"errors"
)

type WriteIndex struct {
	base_index
}

func (this *WriteIndex) Update() {

}

func (this *WriteIndex) Insert(vals ...string) (err error) {
	conn, err := this.conn_pool.getConnection()
	if err != nil {
		return nil, err
	}
	defer this.conn_pool.releaseConnection(conn)

	if err := this.open(conn); err != nil {
		return nil, err
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
