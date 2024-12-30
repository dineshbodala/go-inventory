package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dineshbodala/myinventory/constants"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct{
	Router *mux.Router
	db *sql.DB
}

func (app *App) initialize() error{
	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", constants.DBuser, constants.DBPassword, constants.DBName)
	var err error
	app.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) connect(address string) {
  	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}){
	response, _ := json.Marshal(payload)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string){
		error_message := map[string]string{"error": err}
		sendResponse(w, statusCode, error_message)
}

func (app *App)getProducts(w http.ResponseWriter, r *http.Request){
	fmt.Println("endpoint hit")
	products, err := getProducts(app.db)
	if err != nil{
		sendError(w, http.StatusInternalServerError, err.Error())
		return 
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint hit - getProduct")

    vars := mux.Vars(r)
    key, err := strconv.Atoi(vars["id"])
    if err != nil {
        sendError(w, http.StatusBadRequest, "Invalid product ID")
        return
    }

    p := product{ID: key}
    err = p.getProduct(app.db) 
    if err != nil {
        if err == sql.ErrNoRows {
            sendError(w, http.StatusNotFound, fmt.Sprintf("No product found with ID %d", key))
            return
        }
        sendError(w, http.StatusInternalServerError, "Error retrieving product")
        return
    }

    sendResponse(w, http.StatusOK, p)
}

func (app *App)addProduct(w http.ResponseWriter, r *http.Request){
	var p product 
	fmt.Println("endpoint hit - app.getProduct")
	fmt.Println(r.Body)
	
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil{
		fmt.Println(err)
		sendError(w, http.StatusBadRequest, "Bad request")
		return
	}
	err = p.addProduct(app.db)
	if err != nil{
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, p)

}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint hit - app.updateProduct")

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        sendError(w, http.StatusBadRequest, "Invalid product ID")
        return
    }

    var p product
    err = json.NewDecoder(r.Body).Decode(&p)
    if err != nil {
        fmt.Println(err)
        sendError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    p.ID = id

    err = p.updateProduct(app.db)
    if err != nil {
        sendError(w, http.StatusInternalServerError, err.Error())
        return
    }

    sendResponse(w, http.StatusOK, p)
}




func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit - app.deleteProduct")

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"]) 
    if err != nil {
        sendError(w, http.StatusBadRequest, "Invalid product ID")
        return
    }
	var p product
	p = product{ID: id}
	err = p.deletProduct(app.db)
	if err != nil{
		sendError(w, http.StatusInternalServerError, err.Error())
		return 
	}
	sendResponse(w, http.StatusAccepted, map[string]string{"result": "deleted"})


}






func (app *App) handleRoutes(){
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/products/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.addProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")


}
