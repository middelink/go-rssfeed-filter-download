package qnap

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type qnap struct {
	uri    string
	sid    string
	logout string
}

func (q *qnap) dispatch(path, function string, params url.Values) (map[string]interface{}, error) {
	path = strings.Join([]string{q.uri, path}, "/")
	if q.sid != "" {
		params.Set("sid", q.sid)
	}
	if function != "" {
		params.Set("func", function)
	}
	resp, err := http.PostForm(path, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("res=%+v\n", res)
	m := res.(map[string]interface{})
	if s, ok := m["status"]; ok && s.(float64) != 1 {
		return nil, fmt.Errorf("reason: %v", m["status"].(float64))
	}
	return m, nil
}

func (q *qnap) Close() error {
	if q.sid == "" {
		return nil
	}
	_, err := q.dispatch(q.logout, "", url.Values{})
	q.sid = ""
	return err
}

func GetDefaults() (baseuri, user, pass string) {
	home, _ := os.LookupEnv("HOME")
	data, err := ioutil.ReadFile(path.Join(home, ".local/etc/qnap-downloader"))
	if err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		i := 0
		for scanner.Scan() {
			if i == 0 {
				baseuri = scanner.Text()
			} else if i == 1 {
				user = scanner.Text()
			} else if i == 2 {
				pass = scanner.Text()
			}
			i++
		}
	}
	return
}
