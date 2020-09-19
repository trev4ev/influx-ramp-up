package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Group struct {
	Name    string
	Members map[string]Member
}
type Member struct {
	Name string
	Year int
	Team string
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var params map[string]json.RawMessage
	err := json.Unmarshal(body, &params)
	if err != nil {
		log.Println("Unable to parse body ", err)
		return
	}
	var member Member
	err = json.Unmarshal(params["object"], &member)
	if err != nil {
		log.Println("Unable to parse object ", err)
		return
	}

	function, _ := strconv.Unquote(string(params["function"]))
	groupName, _ := strconv.Unquote(string(params["group"]))

	switch function {
	case "ADD":
		add(groupName, member)
	case "GET":
		fmt.Printf("Trying to GET")
	case "DELETE":
		fmt.Printf("Trying to DELETE")
	case "DELETE_GROUP":
		fmt.Printf("Trying to DELETE")
	}
}

func add(groupName string, member Member) {
	fmt.Println("Called ADD with", groupName, member)
	groups := loadGroups()

	_, found := find(groups, groupName)
	if found {
		// load group
		group := loadGroup(groupName)
		group.Members[member.Name] = member
		// add member
		// TODO save group
	} else {
		groups = append(groups, groupName)
		// new group struct
		members := map[string]Member{member.Name: member}
		group := Group{Name: groupName, Members: members}
		// new file
		fmt.Println("Group:", group)
		// save group
		// save list of groups
	}

	fmt.Println("Groups:", groups)
}

func loadGroups() []string {
	var groups []string
	return groups
}

func loadGroup(groupName string) Group {
	return Group{Name: groupName, Members: map[string]Member{}}
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func main() {
	http.HandleFunc("/", mainHandler)
	fmt.Printf("Listening on server at port 8000\n")
	http.ListenAndServe(":8000", nil)
}
