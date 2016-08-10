package handlersocket

import (
	"strconv"
	"strings"
	"errors"
)

type base_index struct {
	conn_pool *connection_pool
	spec      *hs_index_spec
	opened    bool
}

type hs_index_spec struct {
	db_name    string
	table_name string
	index_name string
	columns    []string
}

func (this *hs_index_spec) hash() string {
	return this.db_name + ":" + this.table_name + ":" + this.index_name + ":" + strings.Join(this.columns, ",")
}

func (this *base_index) sortFieldsList(data map[string]interface{}) []string {
	result := make([]string, len(data))
	num := 0
	for _, name := range this.spec.columns {
		if val, ok := data[name]; ok {
			result[num] = val
			num++
		}
	}
	return result
}

func (this *base_index)parseResult(resp Response) []map[string]interface{} {
	fieldCount, _ := strconv.Atoi(resp.Data[0])
	remainingFields := len(resp.Data) - 1
	if fieldCount <= 0 {
		return nil;
	}

	rs := remainingFields / fieldCount
	rows := make([]map[string]interface{}, rs)

	offset := 1

	for r := 0; r < rs; r++ {
		d := make(map[string]interface{}, fieldCount)
		for f := 0; f < fieldCount; f++ {
			if len(resp.Data[offset + f]) == 1 && int(resp.Data[offset + f][0]) == 0 {
				d[this.spec.columns[f]] = nil
			} else {
				d[this.spec.columns[f]] = resp.Data[offset + f]
			}
		}
		rows[r] = d
		offset += fieldCount
	}

	return rows
}

func (this *base_index) open(conn *hs_Connection) error {
	if _, ok := conn.indexes[this.spec.hash()]; ok {
		return nil
	}

	if err := conn.connect(); err != nil {
		return err
	}

	index_num := conn.getIndexNum(this.spec)
	cols := strings.Join(this.spec.columns, ",")
	a := []string{index_num, this.spec.db_name, this.spec.table_name, this.spec.index_name, cols}

	//conn.mutex.Lock()
	conn.chan_write <- &cmd_openindex{command: "P", params: a}
	message := <-conn.chan_read
	//conn.mutex.Unlock()

	if message.ReturnCode != "0" {
		return errors.New("Error Opening Index: " + message.toString())
	}

	this.opened = true

	return nil
}
