package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

type Response struct {
	Status  string
	Message string
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	// read in request body
	body, _ := ioutil.ReadAll(req.Body)
	var params map[string]json.RawMessage
	err := json.Unmarshal(body, &params)
	if err != nil {
		log.Println("Unable to parse body ", err)
		return
	}

	// convert function and groupName to strings
	function, _ := strconv.Unquote(string(params["function"]))
	groupName, _ := strconv.Unquote(string(params["group"]))

	var response Response

	switch function {
	case "ADD":
		// convert object to Member struct
		var member Member
		err = json.Unmarshal(params["object"], &member)
		if err != nil {
			log.Println("Unable to parse object ", err)
			return
		}
		response = add(groupName, member)
	case "GET":
		response = get(groupName)
	case "DELETE":
		memberName, _ := strconv.Unquote(string(params["object"]))
		response = deleteMember(groupName, memberName)
	case "DELETE_GROUP":
		response = deleteGroup(groupName)
	default:
		response = Response{Status: "error", Message: "not a valid function"}
	}

	// Convert Response struct into json and return response
	js, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func add(groupName string, member Member) Response {
	groups := loadGroups()

	_, found := find(groups, groupName)
	if found {
		// load group
		group := loadGroup(groupName)
		// add member to map
		group.Members[member.Name] = member
		saveGroup(group)
	} else {
		groups = append(groups, groupName)
		// new Member and Group struct
		members := map[string]Member{member.Name: member}
		group := Group{Name: groupName, Members: members}

		saveGroup(group)
		saveGroups(groups)
	}

	response := Response{Status: "success", Message: "member was added to group"}
	return response
}

func get(groupName string) Response {
	groups := loadGroups()
	_, found := find(groups, groupName)
	if found {
		// Read group file and convert to string
		file, _ := ioutil.ReadFile(groupName + ".json")
		groupStr := string(file)
		response := Response{Status: "success", Message: groupStr}
		return response
	} else {
		response := Response{Status: "error", Message: "group not found"}
		return response
	}
}

func deleteMember(groupName string, memberName string) Response {
	groups := loadGroups()

	_, found := find(groups, groupName)
	if found {
		// load group
		group := loadGroup(groupName)
		// Check if member exists, delete if so
		if _, found = group.Members[memberName]; found {
			delete(group.Members, memberName)
			// delete entire group if no more members
			if len(group.Members) <= 0 {
				_ = deleteGroup(groupName)
			} else {
				saveGroup(group)
			}
			response := Response{Status: "success", Message: "member was deleted from group"}
			return response

		} else {
			response := Response{Status: "error", Message: "member not found in the group"}
			return response
		}
	} else {
		response := Response{Status: "error", Message: "group not found"}
		return response
	}
}

func deleteGroup(groupName string) Response {
	groups := loadGroups()

	i, found := find(groups, groupName)
	if found {
		// replace group with last item and splice to be 1 shorter
		groups[i] = groups[len(groups)-1]
		groups[len(groups)-1] = ""
		groups = groups[:len(groups)-1]
		saveGroups(groups)
		_ = os.Remove(groupName + ".json")
		response := Response{Status: "success", Message: "group was deleted"}
		return response
	} else {
		response := Response{Status: "error", Message: "group not found"}
		return response
	}
}

func loadGroups() []string {
	file, _ := ioutil.ReadFile("groups.json")
	var groups []string
	_ = json.Unmarshal([]byte(file), &groups)
	return groups
}

func loadGroup(groupName string) Group {
	file, _ := ioutil.ReadFile(groupName + ".json")
	var group Group
	_ = json.Unmarshal([]byte(file), &group)
	return group
}

func saveGroup(group Group) {
	file, _ := json.MarshalIndent(group, "", " ")
	_ = ioutil.WriteFile(group.Name+".json", file, 0644)
}

func saveGroups(groups []string) {
	file, _ := json.MarshalIndent(groups, "", " ")
	_ = ioutil.WriteFile("groups.json", file, 0644)
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
