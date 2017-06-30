package qnap

import (
	"encoding/base64"
	"net/url"
	"strings"
)

type TaskState uint8

const (
	appNAME    = "downloadstation"
	appVERSION = "V4"

	urlLogin  = "Misc/Login"
	urlLogout = "Misc/Logout"
	urlDir    = "Misc/Dir"
	urlEnv    = "Misc/Env"

	urlTaskStatus     = "Task/Status"
	urlTaskQuery      = "Task/Query"
	urlTaskDetail     = "Task/Detail"
	urlTaskAddUrl     = "Task/AddUrl"
	urlTaskAddTorrent = "Task/AddTorrent"
	urlTaskStart      = "Task/Start"
	urlTaskStop       = "Task/Stop"
	urlTaskPause      = "Task/Pause"
	urlTaskRemove     = "Task/Remove"
	urlTaskPriority   = "Task/Priority"
	urlTaskGetFile    = "Task/GetFile"
	urlTaskSetFile    = "Task/SetFile"

	TaskSeeding     TaskState = 100
	TaskDownloading           = 104
	TaskPaused                = 1
	TaskStopped               = 2
	TaskFinished              = 5

/*
rss:     [:add, :query, :update, :remove, :query_feed, :update_feed, :add_job, :query_job, :update_job, :remove_job],
config:  [:get, :set],
account: [:add, :query, :update, :remove],
addon:   [:query, :enable, :verify, :install, :uninstall, :search],
*/
)

func NewDownloadstation(uri, user, pass string) (*qnap, error) {
	q := &qnap{uri: uri, logout: urlLogout}
	m, err := q.fsdispatch(urlLogin,
		url.Values{"user": {user}, "pass": {base64.StdEncoding.EncodeToString([]byte(pass))}})
	if err != nil {
		return nil, err
	}
	if sid, ok := m["sid"]; ok {
		q.sid = sid.(string)
	}
	return q, nil
}

func (q *qnap) fsdispatch(path string, params url.Values) (map[string]interface{}, error) {
	return q.dispatch(strings.Join([]string{q.uri, appNAME, appVERSION, path}, "/"), "", params)
}

func (q *qnap) Dir() error {
	_, err := q.fsdispatch(urlDir, url.Values{})
	return err
}

func (q *qnap) Env() error {
	_, err := q.fsdispatch(urlEnv, url.Values{})
	return err
}

func (q *qnap) TaskStatus() error {
	_, err := q.fsdispatch(urlTaskStatus, url.Values{})
	return err
}

func (q *qnap) TaskQuery() (map[string]TaskState, error) {
	m, err := q.fsdispatch(urlTaskQuery, url.Values{})
	if err != nil {
		return nil, err
	}
	tasks := make(map[string]TaskState, int(m["total"].(float64)))
	for _, task_intf := range m["data"].([]interface{}) {
		task := task_intf.(map[string]interface{})
		//fmt.Printf("task %s, state: %v\n", task["source_name"], task["state"])

		switch TaskState(task["state"].(float64)) {
		case TaskPaused:
			if m, err = q.fsdispatch("Task/Start", url.Values{"hash": {task["hash"].(string)}}); err != nil {
				return nil, err
			}
		case TaskStopped, TaskFinished:
			if _, err = q.fsdispatch("Task/Remove", url.Values{"hash": {task["hash"].(string)}}); err != nil {
				return nil, err
			}
			continue
		}
		tasks[task["source_name"].(string)] = TaskState(task["state"].(float64))
	}
	return tasks, nil
}

func (q *qnap) TaskDetail() error {
	_, err := q.fsdispatch(urlTaskDetail, url.Values{})
	return err
}

func (q *qnap) TaskAddUrl(magnet string) error {
	_, err := q.fsdispatch(urlTaskAddUrl, url.Values{"temp": {"Download"}, "move": {"Download"}, "url": {magnet}})
	return err
}

func (q *qnap) TaskAddTorrent() error {
	_, err := q.fsdispatch(urlTaskAddTorrent, url.Values{})
	return err
}

func (q *qnap) TaskStart() error {
	_, err := q.fsdispatch(urlTaskStart, url.Values{})
	return err
}

func (q *qnap) TaskStop() error {
	_, err := q.fsdispatch(urlTaskStop, url.Values{})
	return err
}

func (q *qnap) TaskPause() error {
	_, err := q.fsdispatch(urlTaskPause, url.Values{})
	return err
}

func (q *qnap) TaskRemove() error {
	_, err := q.fsdispatch(urlTaskRemove, url.Values{})
	return err
}

func (q *qnap) TaskPriority() error {
	_, err := q.fsdispatch(urlTaskPriority, url.Values{})
	return err
}

func (q *qnap) TaskGetFile() error {
	_, err := q.fsdispatch(urlTaskGetFile, url.Values{})
	return err
}

func (q *qnap) TaskSetFile() error {
	_, err := q.fsdispatch(urlTaskSetFile, url.Values{})
	return err
}
