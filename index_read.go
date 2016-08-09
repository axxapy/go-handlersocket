package handlersocket

import (
	"strings"
	"strconv"
	"errors"
	"github.com/axxapy/go-handlersocket/lg"
)

type hs_IndexRead struct {
	hs_index
}

const TAG = "index:read"

func (this *hs_IndexRead) Find(oper string, limit int, offset int, vals ...string) ([]map[string]interface{}, error) {
	conn, err := this.conn_pool.getConnection()
	if err != nil {
		return nil, err
	}

	defer this.conn_pool.releaseConnection(conn)

	lg.V(TAG, "OPENING INDEX...")
	if err := this.open(conn); err != nil {
		lg.V(TAG, "FAILED TO OPEN INDEX: " + err.Error())
		return nil, err
	}
	lg.V(TAG, "INDEX OPENED")

	cols := strings.Join(vals, "\t")
	colCount := strconv.Itoa(len(vals))
	a := []string{oper, colCount, cols}

	index_num := conn.getIndexNum(this.spec)

	//conn.mutex.Lock()
	conn.chan_write <- &hs_cmd_find{command: index_num, params: a, limit: limit, offset: offset}
	message := <-conn.chan_read
	//conn.mutex.Unlock()

	if message.ReturnCode != "0" {
		return nil, errors.New("Failed to select data: " + message.toString())
	}

	return this.parseResult(message), nil
}
