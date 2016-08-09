package handlersocket

import (
	"strings"
	"io"
)

type HandlerSocket struct {
	conn_pool_read  *connection_pool
	conn_pool_write *connection_pool
}

func New(addr_read string, addr_write string, limit_read int, limit_write int) *HandlerSocket {
	connection := &HandlerSocket{
		conn_pool_read: NewConnectionPool(addr_read, limit_read),
		conn_pool_write: NewConnectionPool(addr_write, limit_write),
	}

	return connection
}

type Response struct {
	ReturnCode string
	Data       []string
}

func (this *Response) toString() string {
	return this.ReturnCode + " " + strings.Join(this.Data, " ")
}

type hs_chan_writer interface {
	write(w io.Writer) (err error)
}

func (this *HandlerSocket) OpenReadIndex(db string, table string, index_name string, columns ...string) *ReadIndex {
	spec := &hs_index_spec{db_name: db, table_name: table, index_name: index_name, columns: columns}
	return &ReadIndex{base_index{
		spec: spec,
		conn_pool: this.conn_pool_read,
		opened:false,
	}}
}

type hs_IndexWrite struct {
	base_index
}
