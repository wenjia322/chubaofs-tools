package server

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/chubaofs/chubaofs-tools/audit-daemon/util"
	"net/http"
	"sync"

	"github.com/spf13/cast"
)

func StartServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc(PathForwardCmd, forwardCmdReq)
	mux.HandleFunc(PathSearchDB, searchDBReq)

	server := &http.Server{
		Addr:    ":" + cast.ToString(port),
		Handler: mux,
	}
	LOG.Debugf("start server on")
	LOG.Fatal(server.ListenAndServe())
}

func forwardCmdReq(w http.ResponseWriter, r *http.Request) {
	var (
		req RequestForwardCmdReq
		err error
	)
	if err = ReadReq(r, &req); err != nil {
		SendErr(w, err)
		return
	}
	cmdReq := RequestCommand{
		Command: req.Command,
		LimitMB: req.LimitMB,
	}
	results := make([]string, len(req.AddrList))
	var wg sync.WaitGroup
	wg.Add(len(req.AddrList))
	for i, addr := range req.AddrList {
		go func(i int, addr string) {
			respData, err := SendByteReq("POST", addr+PathCommand, &cmdReq)
			if err != nil {
				LOG.Errorf("forward cmd req err: addr[%v], req[%v]", addr, cmdReq)
				results[i] = "execute failed!"
				wg.Done()
				return
			}
			results[i] = string(respData)
			wg.Done()
		}(i, addr)
	}
	wg.Wait()

	resp, _ := json.Marshal(&ResponseForwardCmdReq{
		Code:    0,
		Msg:     "execute successfully",
		Results: results,
	})

	if _, err = w.Write(resp); err != nil {
		LOG.Errorf("write to server has err:[%s]", err.Error())
		SendErr(w, err)
		return
	}
}

func searchDBReq(w http.ResponseWriter, r *http.Request) {
	var (
		req     RequestSearch
		cdbResp ResponseCDB
		url     string
	)
	if err := ReadReq(r, &req); err != nil {
		SendErr(w, err)
		return
	}

	if req.Query == "" {
		url = fmt.Sprintf("%v/search/%v?size=%v", req.DBAddr, req.DBTable, req.Size)
	} else {
		url = fmt.Sprintf("%v/search/%v?query=%v&size=%v", req.DBAddr, req.DBTable, req.Query, req.Size)
	}

	respData, err := SendRequest("GET", url, nil)
	if err != nil {
		SendErr(w, err)
		return
	}

	if err := json.Unmarshal(respData, &cdbResp); err != nil {
		SendErr(w, err)
		return
	}

	if cdbResp.Info.Success != 1 {
		err = errors.New(cdbResp.Info.Message)
		SendErr(w, err)
		return
	}

	resp, _ := json.Marshal(&ResponseSearch{
		Code:  0,
		Msg:   "execute successfully",
		Total: cdbResp.Total,
		Hits:  cdbResp.Hits,
	})

	if _, err = w.Write(resp); err != nil {
		LOG.Errorf("write to server has err:[%s]", err.Error())
		SendErr(w, err)
		return
	}

}
