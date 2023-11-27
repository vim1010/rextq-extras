package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

var errLog = log.New(os.Stderr, "", log.Ldate+log.Ltime)

func croak(err error) {
	if err != nil {
		panic(err)
	}
}

func logErr(err error) {
	if err != nil {
		errLog.Println(err)
	}
}

func inventoryFactory() map[string]any {
	return map[string]any{
		"_meta": map[string]any{
			"hostvars": map[string]any{},
		},
		"all": map[string]any{
			"hosts": []string{},
			"vars":  map[string]any{},
			"children": []string{
				"ungrouped",
			},
		},
		"ungrouped": map[string]any{
			"hosts": []string{},
			"vars":  map[string]any{},
		},
	}
}

func getInventory(client *Service, projectID string) (res map[string]any, err error) {
	res = inventoryFactory()
	d, err := client.Call("get_host_groups", map[string]any{
		"project_id": projectID,
	})
	croak(err)
	lockedHosts := make([]string, 0)
	for _, h := range d {
		hostGroupID, ok := h["host_group_id"]
		if !ok {
			croak(errors.New(fmt.Sprintf("bad host group ID: %v", h["host_group_id"])))
		}
		hostGroupName, ok := h["host_group_name"].(string)
		if !ok {
			croak(errors.New(fmt.Sprintf("bad host group name: %v", h["host_group_name"])))
		}
		if hostGroupName == "all" {
			logErr(errors.New("cannot have host group named 'all', skipping ..."))
			continue
		}
		var groupVars any
		if h["data"] != nil {
			groupVars = h["data"]
		}
		hosts := make([]string, 0)
		g, err := client.Call("get_host_group_ips", map[string]any{
			"host_group_id": hostGroupID,
			"project_id":    projectID,
		})
		croak(err)
		for _, x := range g {
			hostIP, ok := x["host_ip"].(string)
			if !ok {
				hostID := x["host_id"]
				logErr(errors.New(fmt.Sprintf("bad host_ip for host_id [%d]", hostID)))
				continue
			}
			hostLocked, ok := x["host_locked"].(bool)
			if !ok || hostLocked {
				lockedHosts = append(lockedHosts, hostIP)
				continue
			}
			hosts = append(hosts, hostIP)
		}
		if groupVars == nil {
			groupVars = map[string]any{}
		}
		res[hostGroupName] = map[string]any{
			"hosts": hosts,
			"vars":  groupVars,
		}
	}
	return res, err
}

func main() {
	baseURL := os.Getenv("REX_BASE_URL")
	user := os.Getenv("REX_USER")
	pass := os.Getenv("REX_PASS")
	projectID := os.Getenv("REX_PROJECT_ID")
	client := NewService(baseURL, user, pass)
	res, err := getInventory(client, projectID)
	croak(err)
	js, err := json.Marshal(res)
	croak(err)
	fmt.Println(string(js))
}
