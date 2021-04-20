package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const port = ":3000"

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", rootPage)
	router.HandleFunc("/productget/{fetchCount}", getWithParam).Methods("GET")
	router.HandleFunc("/productpost", postWithBody).Methods("POST")

	fmt.Println("Serving @ http://127.0.0.1" + port)
	log.Fatal(http.ListenAndServe(port, router))

}

func rootPage(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("This is root page"))
}

func postWithBody(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	// fmt.Println(r.Body)
	// fmt.Println(decoder)
	var t interface{}
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	fmt.Println("Data from post here, parse to Json")
	content, _ := json.Marshal(t)
	fmt.Println(string(content))

	w.Header().Set("content-type", "application/json")
	w.Write([]byte("Done post"))
}

func getWithParam(w http.ResponseWriter, r *http.Request) {

	sfetchCount, errInput := strconv.ParseFloat(mux.Vars(r)["fetchCount"], 1)//This reads params

	fetchCount := int(sfetchCount)

	if errInput != nil {
		fmt.Println(errInput.Error())
	} else {
		// fetchCount = int(float64(len(productList)) * fetchCountPercentage / 100)
		if fetchCount > len(productList) {
			fetchCount = len(productList)
		}
	}

	// write to response
	jsonList, err := json.Marshal(productList[0:fetchCount])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {
		w.Header().Set("content-type", "application/json")
		w.Write(jsonList)
	}

}

type product struct {
	Name  string
	Price float64
	Count int
}

var productList = []product{

	product{"p1", 25.0, 30},
	product{"p2", 20.0, 10},
	product{"p3", 250.0, 320},
	product{"p4", 256.0, 730},
	product{"p5", 24.0, 340},
	product{"p6", 10.0, 300},
	product{"p7", 100.0, 230},
	product{"p8", 2543.0, 120},
	product{"p9", 255.0, 10},
	product{"p10", 175.0, 20},
}