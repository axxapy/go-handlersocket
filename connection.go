package handlersocket

import (
	"strings"
	"bufio"
	"net"
	"strconv"
	"sync"
	"io"
)

type hs_Connection struct {
	index_counter int

	indexes map[string]string

	conn_spec string

	mutex *sync.Mutex

	conn net.Conn
	chan_read chan Response
	chan_write chan hs_chan_writer
}

func NewConnection(addr string) *hs_Connection {
	return &hs_Connection{
		conn_spec: addr,
		mutex: &sync.Mutex{},
		indexes: make(map[string]string),
	}
}

func (this *hs_Connection) close() {
	if this.conn == nil {
		return
	}
	this.indexes = make(map[string]string)
	this.index_counter = 0
	this.conn.Close()
	this.conn = nil
	close(this.chan_read)
	close(this.chan_write)
}

func (this *hs_Connection) reader() {
	defer this.close()

	br := bufio.NewReader(this.conn)
	var retString string
	var bytes []byte
	for {
		b, err := br.ReadByte()
		if err != nil {
			// TODO(adg) handle error
			if err == io.EOF {
				break
			}
			break
		}

		if string(b) != "\n" {
			bytes = append(bytes, b)
		} else {
			retString = string(bytes)
			strs := strings.Split(retString, "\t") //, -1)
			hsr := Response{ReturnCode: strs[0], Data: strs[1:]}
			this.chan_read <- hsr
			retString = ""
			bytes = []byte{}
		}
	}
}

func (this *hs_Connection) writer() {
	defer this.close()

	bw := bufio.NewWriter(this.conn)

	for f := range this.chan_write {
		if err := f.write(bw); err != nil {
			panic(err)
		}

		if err := bw.Flush(); err != nil {
			panic(err)
		}
	}
}

func (this *hs_Connection) connect() error {
	if this.conn != nil {
		return nil
	}

	remote_net := "tcp"
	remote_addr := this.conn_spec
	arr := strings.Split(this.conn_spec, "://")
	if len(arr) == 2 {
		remote_net = arr[0]
		remote_addr = arr[1]
	}

	conn, err := net.Dial(remote_net, remote_addr)
	if err != nil {
		return err
	}

	this.conn = conn
	this.chan_read = make(chan Response)
	this.chan_write = make(chan hs_chan_writer)

	go this.writer()
	go this.reader()

	return nil
}

func (this *hs_Connection) getIndexNum(index *hs_index_spec) string {
	hash := index.hash()
	if val, ok := this.indexes[hash]; ok {
		return val
	}
	//this.mutex.Lock()
	this.index_counter++
	this.indexes[hash] = strconv.Itoa(this.index_counter)
	//this.mutex.Unlock()
	return this.indexes[hash]
}
