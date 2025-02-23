package controller

import (
	"databaseconnection/model"
	"databaseconnection/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Welcome")
}

func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	courses, err := service.FindAll()
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(courses)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusConflict)
		return
	}

}

func AddARecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var course model.Course
	var mymap map[string]interface{} = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&mymap)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusConflict)
		return
	}
	if len(mymap) != 3 {
		http.Error(w, "Mismatched no.of required json fields", http.StatusBadRequest)
		return
	}
	str, ok := mymap["name"].(string)
	if !ok {
		http.Error(w, "json input field \"name\" is not a string", http.StatusBadRequest) //-->using this will set content-type to text/plain
		// w.WriteHeader(http.StatusBadRequest)                                              //--using this will retain the content-type in header
		// w.Write([]byte("json input field \"name\" is not a string"))
		return
	}
	str2, ok := mymap["price"].(float64)
	if !ok {
		http.Error(w, "json input field \"price\" is not a float64", http.StatusBadRequest)
		return
	}
	str3, ok := mymap["videocount"].(float64)
	if !ok {
		http.Error(w, "json input field 'videocount' is not a float64"+fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	if str3 != float64(int(str3)) {
		http.Error(w, "json input field 'videocount' should be an integer"+fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	course.Name = str
	course.Price = str2
	course.Videocount = int(str3)
	err = service.AddOneRecord(course)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(course)
}

func DeleteByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var courses []model.Course
	count, err := service.DeleteRecord(params["name"], &courses)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusConflict)
		return
	}
	if count == 0 {
		http.Error(w, "No records found with provided courseName", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusFound)
	w.Write([]byte("Deleted Items Count : " + strconv.Itoa(int(count))))
}

func GetByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var course []model.Course
	params := mux.Vars(r)
	err := service.FindByName(params["name"], &course)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusConflict)
		return
	}
	if len(course) == 0 {
		http.Error(w, "No Records Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusFound) // http status code once written cannot be overriden by writing again.
	json.NewEncoder(w).Encode(course)
}

func UpdateByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var course_rbody model.Course
	var json_count int = 0
	var mp map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusConflict)
		return
	}
	if res, ok := mp["name"].(string); ok {
		course_rbody.Name = res
		json_count++
	} else if mp["name"] != nil {
		http.Error(w, "Please provide string value for name field in json", http.StatusBadRequest)
		return
	}
	if res, ok := mp["price"].(float64); ok {
		course_rbody.Price = res
		json_count++
	} else if mp["price"] != nil {
		http.Error(w, "Please provide float64 value for price field in json", http.StatusBadRequest)
		return
	}
	if res, ok := mp["videocount"].(float64); ok {
		if res != float64(int((res))) {
			http.Error(w, "Please provide int value for videocount field in json", http.StatusBadRequest)
			return
		}
		course_rbody.Videocount = int(res)
		json_count++
	} else if mp["videocount"] != nil {
		http.Error(w, "Please provide int value for videocount field in json", http.StatusBadRequest)
		return
	}
	if json_count != len(mp) {
		http.Error(w, "Json input consists unnecessary extra fields", http.StatusBadRequest)
		return
	}
	if course_rbody.Name == "" && course_rbody.Price == 0 && course_rbody.Videocount == 0 {
		http.Error(w, "Request body is empty for updation", http.StatusNotAcceptable)
		return
	}
	count, err := service.UpdateByName(params["name"], course_rbody)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusConflict)
		return
	}
	if count == 0 {
		http.Error(w, "No records Updated bcos either records don't exist or already updated", http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Updated count :" + strconv.Itoa(int(count))))
}
