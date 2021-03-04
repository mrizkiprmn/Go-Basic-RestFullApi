package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Product struct {
	ID string `json:id`
	Name string `json:name`
	Quantity int `json:quantity`
}

var (
	database = make(map[string]Product)
)

func SetJSONRespon(res http.ResponseWriter, message []byte, httpCode int) {
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(httpCode)
	res.Write(message)
}

func main() {
	
	database["001"] = Product{ID: "001", Name: "Iphone X", Quantity: 5}
	database["002"] = Product{ID: "002", Name: "Samsung S12", Quantity: 10}



	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
		message := []byte(`{"message": "Server Connected"}`)
		SetJSONRespon(res, message, http.StatusOK)
	})

	http.HandleFunc("/products", func(res http.ResponseWriter, req *http.Request){
		if req.Method != "GET" {
			message := []byte(`{"message": "invalid http method"}`)
			SetJSONRespon(res, message, http.StatusMethodNotAllowed)
			return
		}

		var products []Product

		for _, product := range database {
			products = append(products, product)
		}

		productJSON, err := json.Marshal(&products)
		if err != nil {
			message := []byte(`{"message": "error parsing data"}`)
			SetJSONRespon(res, message, http.StatusInternalServerError)
			return
		}	

		SetJSONRespon(res, productJSON, http.StatusOK)
	})

	http.HandleFunc("/addproducts", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			message := []byte(`{"message": "invalid http method"}`)
			SetJSONRespon(res, message, http.StatusMethodNotAllowed)
			return
		}

		var product Product

		payload := req.Body

		defer req.Body.Close()

		err := json.NewDecoder(payload).Decode(&product)
		if err != nil {
			message := []byte(`{"message": "error parsing data"}`)
			SetJSONRespon(res, message, http.StatusInternalServerError)
			return
		}

		database[product.ID] =  product

		message := []byte (`{"message": "add product success"}`)
		SetJSONRespon(res, message, http.StatusCreated)
	})

	http.HandleFunc("/product", func(res http.ResponseWriter, req *http.Request){
		if req.Method != "GET" {
			message := []byte(`{"message": "invalid http method}`)
			SetJSONRespon(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte (`{"message": "required product id"}`)
			SetJSONRespon(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]
		product, ok := database[id]
		if !ok {
			message := []byte (`{"message": "product not found"}`)
			SetJSONRespon(res, message, http.StatusOK)
			return
		}

		productJSON, err := json.Marshal(&product)
		if err != nil {
			message := []byte (`{"message": "error when parsing data"}`)
			SetJSONRespon(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONRespon(res, productJSON, http.StatusOK)
	})

	http.HandleFunc("/deleteproducts", func(res http.ResponseWriter, req *http.Request){

		if req.Method != "DELETE" {
			message := []byte(`{"message": "invalid http method}`)
			SetJSONRespon(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte (`{"message": "required product id"}`)
			SetJSONRespon(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]
		product, ok := database[id]
		if !ok {
			message := []byte (`{"message": "product not found"}`)
			SetJSONRespon(res, message, http.StatusOK)
			return
		}

		delete(database, id)

		productJSON, err := json.Marshal(&product)
		if err != nil {
			message := []byte(`{"message": "error when parsing data}`)
			SetJSONRespon(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONRespon(res, productJSON, http.StatusOK)
	})

	http.HandleFunc("/updateproducts", func(res http.ResponseWriter, req *http.Request) {

		if req.Method != "PUT" {
			message := []byte(`{"message": "invalid http method}`)
			SetJSONRespon(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte (`{"message": "required product id"}`)
			SetJSONRespon(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]
		product, ok := database[id]
		if !ok {
			message := []byte (`{"message": "product not found"}`)
			SetJSONRespon(res, message, http.StatusOK)
			return
		}

		var newProduct Product

		payload := req.Body

		defer req.Body.Close()

		err := json.NewDecoder(payload).Decode(&newProduct)
		if err != nil {
			message := []byte (`{"message": "error when parsing product"}`)
			SetJSONRespon(res, message, http.StatusInternalServerError)
			return
		}

		product.Name = newProduct.Name
		product.Quantity = newProduct.Quantity

		database[product.ID] = product

		productJSON, err := json.Marshal(&product)
		if err != nil {
			message := []byte (`{"message": "error when parsing product"}`)
			SetJSONRespon(res, message, http.StatusInternalServerError)
			return
		}

		SetJSONRespon(res, productJSON, http.StatusOK)
	})

	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}