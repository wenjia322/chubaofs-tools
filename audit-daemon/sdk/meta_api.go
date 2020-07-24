package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	. "github.com/chubaofs/chubaofs-tools/audit-daemon/util"
	"github.com/chubaofs/chubaofs-tools/audit-daemon/util/raft"
)

func GetAllDentry(metaAddr, partitionId string) (dentryMap map[string]*raft.Dentry, err error) {
	url := fmt.Sprintf("http://%v/getAllDentry?pid=%v", metaAddr, partitionId)
	respData, err := SendByteReq("GET", url, nil)
	if err != nil {
		return
	}

	dec := json.NewDecoder(bytes.NewBuffer(respData))
	dec.UseNumber()

	// It's the "items". We expect it to be an array
	if err = parseToken(dec, '['); err != nil {
		return
	}
	// Read items (large objects)
	dentryMap = make(map[string]*raft.Dentry)
	for dec.More() {
		// Read next item (large object)
		lo := &raft.Dentry{}
		if err = dec.Decode(lo); err != nil {
			return
		}
		if d, exist := dentryMap[strconv.FormatUint(lo.Inode, 10)]; exist {
			LOG.Warningf("inode exist: inode[%v], pid[%v], [%v] will be override", lo.Inode, partitionId, d)
		}
		dentryMap[strconv.FormatUint(lo.Inode, 10)] = lo
	}
	// Array closing delimiter
	if err = parseToken(dec, ']'); err != nil {
		return
	}
	return
}

func parseToken(dec *json.Decoder, expectToken rune) (err error) {
	t, err := dec.Token()
	if err != nil {
		return
	}
	if delim, ok := t.(json.Delim); !ok || delim != json.Delim(expectToken) {
		err = fmt.Errorf("expected token[%v], delim[%v], ok[%v]", string(expectToken), delim, ok)
		return
	}
	return
}
