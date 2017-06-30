package qnap

import (
	"encoding/base64"
	"net/url"
)

func NewFilestation(uri, user, pass string) (*qnap, error) {
	q := &qnap{uri: uri, logout: "/cgi-bin/filemanager/wfm2Logout.cgi"}
	m, err := q.dispatch("/cgi-bin/filemanager/wfm2Login.cgi", "", url.Values{"user": {user}, "pwd": {base64.StdEncoding.EncodeToString([]byte(pass))}})
	if err != nil {
		return nil, err
	}
	if sid, ok := m["sid"]; ok {
		q.sid = sid.(string)
	}
	return q, nil
}

// get_list: [:dir, :filename, :flv_720, :hidden_file, :is_iso, :limit, :list_mode, :mp4_360, :mp4_720, :path, :sort, :start, :type],
func (q *qnap) GetList(path string) (map[string]bool, error) {
	m, err := q.dispatch("/cgi-bin/filemanager/utilRequest.cgi", "get_list", url.Values{"start": {"0"}, "limit": {"10000"}, "path": {path}})
	if err != nil {
		return nil, err
	}
	items := make(map[string]bool, int(m["total"].(float64)))
	for _, item_intf := range m["datas"].([]interface{}) {
		item := item_intf.(map[string]interface{})
		//fmt.Printf("item %s, folder: %v\n", item["filename"], item["isfolder"])
		items[item["filename"].(string)] = item["isfolder"].(float64) != 0
	}
	return items, nil
}
