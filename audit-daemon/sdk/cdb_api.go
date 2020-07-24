package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/chubaofs/chubaofs-tools/audit-daemon/util"
)

const DESC Sequence = "desc"
const ASC Sequence = "asc"

type Sequence string

type DBConfig struct {
	Addr        string // Address of chubaodb
	RaftTable   string // 'raft' table name of chubaodb
	DentryTable string // 'dentry' table name of chubaodb
}

func NewDBConfig(dbAddr, raftTable, dentryTable string) *DBConfig {
	return &DBConfig{
		Addr:        dbAddr,
		RaftTable:   raftTable,
		DentryTable: dentryTable,
	}
}

func (dbc *DBConfig) Insert(table, index string, body []byte) (err error) {
	chubaodbAddr := fmt.Sprintf("%v/put/%v", dbc.Addr, table)
	url := fmt.Sprintf("http://%v/%v", chubaodbAddr, index)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		LOG.Errorf("send insert item request[%s]: new request err: [%s]", url, err.Error())
		return
	}

	do, err := http.DefaultClient.Do(req)
	if err != nil {
		LOG.Errorf("send insert item request[%s]: do client err: [%s]", url, err.Error())
		return
	}

	if do.StatusCode != 200 {
		resBody := do.Body
		buf := new(bytes.Buffer)
		buf.ReadFrom(resBody)
		LOG.Warningf("send insert request[%s]: status code: [%v]", url, do.StatusCode)
		return fmt.Errorf("post has status:[%d] body:[%s]", do.StatusCode, buf.String())
	}

	all, err := ioutil.ReadAll(do.Body)
	_ = do.Body.Close()
	if err != nil {
		LOG.Errorf("send insert item request[%s]: read body err: [%s]", url, err.Error())
		return err
	}

	var resp Response
	if err = json.Unmarshal(all, &resp); err != nil {
		LOG.Errorf("send insert item request[%s]: unmarshal err: [%s]", url, err.Error())
		return err
	}

	if resp.Code > 200 {
		LOG.Warningf("send insert item request[%s]: response code: [%v]", url, resp.Code)
		return fmt.Errorf(resp.Msg)
	}

	return nil
}

func (dbc *DBConfig) Query(key, value string, size int) ([]*HitItem, error) {
	url := fmt.Sprintf("%v/search/%v?query=%v:%v&size=%v", dbc.Addr, dbc.RaftTable, key, value, size)

	respData, err := SendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var cdbResp ResponseCDB
	if err := json.Unmarshal(respData, &cdbResp); err != nil {
		return nil, err
	}

	if cdbResp.Info.Success != 1 {
		err = errors.New(cdbResp.Info.Message)
		return nil, err
	}

	return cdbResp.Hits, nil
}

func (dbc *DBConfig) QueryAnd(queryMap map[string]interface{}, size int) ([]*HitItem, error) {
	var query string
	var count int
	for k, v := range queryMap {
		query = fmt.Sprintf("%v:%v", k, v)
		count++
		if count < len(queryMap) {
			query = query + " AND "
		}
	}
	url := fmt.Sprintf("%v/search/%v?query=%v&size=%v", dbc.Addr, dbc.RaftTable, query, size)

	respData, err := SendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var cdbResp ResponseCDB
	if err := json.Unmarshal(respData, &cdbResp); err != nil {
		return nil, err
	}

	if cdbResp.Info.Success != 1 {
		err = errors.New(cdbResp.Info.Message)
		return nil, err
	}

	return cdbResp.Hits, nil
}

func (dbc *DBConfig) QuerySortTop(table string, queryMap map[string]interface{}, sortFiled string, seq Sequence) ([]byte, error) {
	var query string
	var count int
	for k, v := range queryMap {
		query = fmt.Sprintf("%v:%v", k, v)
		count++
		if count < len(queryMap) {
			query = query + " AND "
		}
	}
	url := fmt.Sprintf("%v/search/%v?query=%v&sort=%v:%v", dbc.Addr, table, query, sortFiled, seq)

	respData, err := SendRequest("GET", url, nil)
	if err != nil {
		LOG.Errorf("query chubaodb err: url[%v], err[%v]", url, err)
		return nil, err
	}

	var cdbResp ResponseCDB
	if err := json.Unmarshal(respData, &cdbResp); err != nil {
		LOG.Errorf("query chubaodb err: url[%v], err[%v]", url, err)
		return nil, err
	}

	if cdbResp.Info.Success != 1 {
		err = errors.New(cdbResp.Info.Message)
		LOG.Errorf("query chubaodb err: url[%v], err[%v]", url, err)
		return nil, err
	}

	if cdbResp.Total > 0 {
		if cdbResp.Total > 1 {
			LOG.Warningf("search cdb: url[%v] total[%v]", url, cdbResp.Total)
		}
		var data []byte
		if data, err = json.Marshal(cdbResp.Hits[0].Doc.Source); err != nil {
			LOG.Errorf("json marshal response data err: data[%v], err[%v]", cdbResp.Hits[0].Doc.Source, err)
			return nil, err
		}
		return data, nil
	} else {
		return nil, nil
	}
}
