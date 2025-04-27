package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	//"github.com/gorilla/mux"
)

// type Student struct {
// 	Id    int
// 	Name  string
// 	Age   int
// 	Score int
// }

type Student struct {
	Id    int    `json:"id,omitempty"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Score int    `json:"score,omitempty"`
}

var students map[int]Student
var lastId int

// func MakeWebHandler() http.Handler {
// 	mux := mux.NewRouter()
// 	mux.HandleFunc("/students", GetStudentListHandler).Methods("GET")
// 	mux.HandleFunc("/students/{id:[0-9]+}", GetStudentHandler).Methods("GET")
// 	mux.HandleFunc("/students", PostStudentHandler).Methods("GET")
// 	mux.HandleFunc("/students/{id:[0-9]+}", DeleteStudentHandler).Methods("DELETE")

// 	students = make(map[int]Student)
// 	students[1] = Student{1, "aaa", 16, 87}
// 	students[2] = Student{2, "bbb", 18, 98}
// 	lastId = 2

// 	return mux
// }

func SetupHandlers(g *gin.Engine) {
	// 웹핸들러를 셋팅
	g.GET("/students", GetStudentsHandler)
	g.GET("/student/:id", GetStudentHandler)
	g.POST("/student", PostStudentHandler)
	g.DELETE("/student/:id", DeleteStudentHandler)

	students = make(map[int]Student)
	students[1] = Student{1, "aaa", 16, 87}
	students[2] = Student{2, "bbb", 18, 98}
	lastId = 2
}

type Students []Student // Id로 정렬하는 인터페이스

func (s Students) Len() int {
	return len(s)
}
func (s Students) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Students) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

// func GetStudentListHandler(w http.ResponseWriter, r *http.Request) {
// 	list := make(Students, 0)
// 	for _, student := range students {
// 		list = append(list, student)
// 		sort.Sort(list)
// 		w.WriteHeader(http.StatusOK)
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(list)
// 	}
// }

func GetStudentsHandler(c *gin.Context) {
	list := make(Students, 0)
	for _, student := range students {
		list = append(list, student)

		sort.Sort(list)
		c.JSON(http.StatusOK, list) // JSON 포맷으로 변환
	}
}

// func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r) // id를 가져옵니다.
// 	id, _ := strconv.Atoi(vars["id"])
// 	student, ok := students[id]
// 	if !ok {
// 		w.WriteHeader(http.StatusNotFound)
// 		// id에 해당하는 학생이 없으면 에러
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(student)
// }

func GetStudentHandler(c *gin.Context) {
	idstr := c.Param("id")
	if idstr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	student, ok := students[id]
	if !ok {
		c.AbortWithStatus(http.StatusNotFound) // 에러 반환
		return
	}
	c.JSON(http.StatusOK, student)
}

// func PostStudentHandler(w http.ResponseWriter, r *http.Request) {
// 	var student Student
// 	err := json.NewDecoder(r.Body).Decode(&student) // JSON 데이터 반환
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	lastId++
// 	student.Id = lastId // id를 증가시킨 후 앱에 등록
// 	students[lastId] = student
// 	w.WriteHeader(http.StatusCreated)
// }

func PostStudentHandler(c *gin.Context) {
	var student Student

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	lastId++
	student.Id = lastId
	students[lastId] = student
	c.String(http.StatusCreated, "Success to add id: %d", lastId)
}

// func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r) // id를 가져옵니다.
// 	id, _ := strconv.Atoi(vars["id"])
// 	_, ok := students[id]
// 	if !ok {
// 		w.WriteHeader(http.StatusNotFound)
// 		// id에 해당하는 학생이 없으면 에러
// 		return
// 	}
// 	delete(students, id)
// 	w.WriteHeader(http.StatusOK) // statusOK 반환
// }

func DeleteStudentHandler(c *gin.Context) {
	idstr := c.Param("id")
	if idstr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	delete(students, id)
	c.String(http.StatusOK, "success to delete")
}

func main() {
	// http.ListenAndServe(":3000", MakeWebHandler())
	r := gin.Default()
	SetupHandlers(r)
	r.Run(":3000")
}
