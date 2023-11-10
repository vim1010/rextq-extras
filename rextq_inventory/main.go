package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
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

type Meta struct {
	Hostvars map[string]any `json:"hostvars"`
}

type All struct {
	Hosts    []string       `json:"hosts"`
	Vars     map[string]any `json:"hosts"`
	Children []string       `json:"children"`
}

type Ungrouped struct {
	Hosts []string       `json:"hosts"`
	Vars  map[string]any `json:"vars"`
}

type Inventory struct {
	Meta      Meta      `json:"_meta"`
	All       All       `json:"all"`
	Ungrouped Ungrouped `json:"ungrouped"`
}

func getInventory(client, projectID) (res *Inventory, err error) {
	inventory = &Inventory{}
	res, err = client.POST("get_host_groups", map[string]any{
		"project_id": projectID,
	})
	croak(err)
	for _, h := range res {
		hostGroupID, ok := h["host_group_id"].(int64)
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
		var groupVars map[string]any
		if h["data"] != nil {
			groupVars = h["data"]
		}
	}
}

func main() {
	baseURL := os.Getenv("REX_BASE_URL")
	user := os.Getenv("REX_USER")
	pass = os.Getenv("REX_PASS")
	projectID = os.Getenv("REX_PROJECT_ID")
	client := NewService(baseURL, user, pass)

}
