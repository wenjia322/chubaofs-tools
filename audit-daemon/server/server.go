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
			respData, err := SendDaemonReq(addr+PathCommand, &cmdReq)
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
		req RequestSearch
		url string
	)
	if err := ReadReq(r, &req); err != nil {
		SendErr(w, err)
		return
	}

	domain := fmt.Sprintf("%v/search/%v", req.DBAddr, req.DBTable)
	if req.Fields == "" {
		url = fmt.Sprintf("%v?query=%v", domain, req.Query)
	} else {
		url = fmt.Sprintf("%v?query=%v&def_fields=%v", domain, req.Query, req.Fields)
	}

	respData, err := SendRequest(url, &req)
	if err != nil {
		SendErr(w, err)
		return
	}

	var cdbResp ResponseCDB
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
		Code: 0,
		Msg:  "execute successfully",
		Hits: cdbResp.Hits,
	})

	if _, err = w.Write(resp); err != nil {
		LOG.Errorf("write to server has err:[%s]", err.Error())
		SendErr(w, err)
		return
	}

}
