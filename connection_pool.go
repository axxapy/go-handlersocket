package handlersocket

import (
	"sync"
	"errors"
	"strconv"
)

type connection_pool struct {
	addr              string
	conn_max          int

	mu                sync.Mutex // protects following fields

	conn_cnt          int
	chan_conn_free    chan *hs_Connection
}

func NewConnectionPool(addr string, max_connections int) *connection_pool {
	pool := &connection_pool{
		addr: addr,
		conn_max: max_connections,
		chan_conn_free: make(chan *hs_Connection, max_connections),
	}

	return pool
}

func (this *connection_pool) getConnection() (*hs_Connection, error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	var conn *hs_Connection

	select {
		case conn = <- this.chan_conn_free:
			return conn, nil

		default:
			if (this.conn_max <= this.conn_cnt) {
				return nil, errors.New("Connection limit exeeded: " + strconv.Itoa(this.conn_cnt))
			}

			conn = NewConnection(this.addr)
			this.conn_cnt++
	}

	return conn, nil
}

func (this *connection_pool) releaseConnection(conn *hs_Connection) {
	this.chan_conn_free <- conn
}
