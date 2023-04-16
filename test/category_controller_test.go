package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"golang-restful-api/app"
	"golang-restful-api/controller"
	"golang-restful-api/helper"
	"golang-restful-api/middleware"
	"golang-restful-api/model/domain"
	"golang-restful-api/repository"
	"golang-restful-api/service"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func setUpDB() *sql.DB {
	db, err := sql.Open("mysql", "devtest:root@tcp(localhost:3306)/devtest")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxIdleConns(50)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}

func setUpRouter(db *sql.DB) http.Handler {
	validate := validator.New()

	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)

	router := app.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB) {
	_, err := db.Exec("TRUNCATE category")
	helper.PanicIfError(err)
}

func generateData(db *sql.DB) domain.Category {
	tx, errBegin := db.Begin()
	helper.PanicIfError(errBegin)

	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "name_test",
	})
	errCommit := tx.Commit()
	helper.PanicIfError(errCommit)

	return category
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setUpDB()
	defer truncateCategory(db)
	router := setUpRouter(db)

	requestBody := strings.NewReader(`{"name" : "name_test"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, errReadAll := io.ReadAll(response.Body)
	helper.PanicIfError(errReadAll)

	var responseBody map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(errUnmarshal)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, http.StatusText(http.StatusOK), responseBody["status"])
	assert.Equal(t, "name_test", responseBody["data"].(map[string]interface{})["name"])
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setUpDB()
	defer truncateCategory(db)
	router := setUpRouter(db)

	requestBody := strings.NewReader(`{"name" : ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, errReadAll := io.ReadAll(response.Body)
	helper.PanicIfError(errReadAll)

	var responseBody map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(errUnmarshal)

	assert.Equal(t, http.StatusBadRequest, int(responseBody["code"].(float64)))
	assert.Equal(t, http.StatusText(http.StatusBadRequest), responseBody["status"])
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setUpDB()
	defer truncateCategory(db)
	category := generateData(db)
	router := setUpRouter(db)

	requestBody := strings.NewReader(`{"name" : "name_test_2"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, errReadAll := io.ReadAll(response.Body)
	helper.PanicIfError(errReadAll)

	var responseBody map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(errUnmarshal)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, http.StatusText(http.StatusOK), responseBody["status"])
	assert.Equal(t, "name_test_2", responseBody["data"].(map[string]interface{})["name"])
}

func TestUpdateCategoryFailed(t *testing.T) {
	db := setUpDB()
	defer truncateCategory(db)
	category := generateData(db)
	router := setUpRouter(db)

	requestBody := strings.NewReader(`{"name" : ""}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, errReadAll := io.ReadAll(response.Body)
	helper.PanicIfError(errReadAll)

	var responseBody map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(errUnmarshal)

	assert.Equal(t, http.StatusBadRequest, int(responseBody["code"].(float64)))
	assert.Equal(t, http.StatusText(http.StatusBadRequest), responseBody["status"])
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setUpDB()
	defer truncateCategory(db)
	category := generateData(db)
	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("X-API-Key", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, errReadAll := io.ReadAll(response.Body)
	helper.PanicIfError(errReadAll)

	var responseBody map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(errUnmarshal)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, http.StatusText(http.StatusOK), responseBody["status"])
	assert.Equal(t, nil, responseBody["data"])
}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setUpDB()
	defer truncateCategory(db)
	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/404", nil)
	request.Header.Add("X-API-Key", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, errReadAll := io.ReadAll(response.Body)
	helper.PanicIfError(errReadAll)

	var responseBody map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(errUnmarshal)

	assert.Equal(t, http.StatusNotFound, int(responseBody["code"].(float64)))
	assert.Equal(t, http.StatusText(http.StatusNotFound), responseBody["status"])
}

func TestGetCategoryByIdSuccess(t *testing.T) {
	db := setUpDB()
	defer truncateCategory(db)
	category := generateData(db)
	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("X-API-Key", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, errReadAll := io.ReadAll(response.Body)
	helper.PanicIfError(errReadAll)

	var responseBody map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(errUnmarshal)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, http.StatusText(http.StatusOK), responseBody["status"])
	assert.Equal(t, category.Name, responseBody["data"].(map[string]interface{})["name"])
}

func TestGetCategoryByIdFailed(t *testing.T) {
	db := setUpDB()
	defer truncateCategory(db)
	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/404", nil)
	request.Header.Add("X-API-Key", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, errReadAll := io.ReadAll(response.Body)
	helper.PanicIfError(errReadAll)

	var responseBody map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(errUnmarshal)

	assert.Equal(t, http.StatusNotFound, int(responseBody["code"].(float64)))
	assert.Equal(t, http.StatusText(http.StatusNotFound), responseBody["status"])
}

func TestGetAllCategorySuccess(t *testing.T) {
	db := setUpDB()
	defer truncateCategory(db)
	category1 := generateData(db)
	category2 := generateData(db)
	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	request.Header.Add("X-API-Key", "RAHASIA")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, errReadAll := io.ReadAll(response.Body)
	helper.PanicIfError(errReadAll)

	var responseBody map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(errUnmarshal)

	var categories = responseBody["data"].([]interface{})
	categoryResponse1 := categories[0].(map[string]interface{})
	categoryResponse2 := categories[1].(map[string]interface{})

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, http.StatusText(http.StatusOK), responseBody["status"])
	assert.Equal(t, category1.Id, int(categoryResponse1["id"].(float64)))
	assert.Equal(t, category1.Name, categoryResponse1["name"])
	assert.Equal(t, category2.Id, int(categoryResponse2["id"].(float64)))
	assert.Equal(t, category2.Name, categoryResponse2["name"])
}

func TestUnauthorized(t *testing.T) {
	db := setUpDB()
	defer truncateCategory(db)
	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	request.Header.Add("X-API-Key", "SALAH")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)

	body, errReadAll := io.ReadAll(response.Body)
	helper.PanicIfError(errReadAll)

	var responseBody map[string]interface{}
	errUnmarshal := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(errUnmarshal)

	assert.Equal(t, http.StatusUnauthorized, int(responseBody["code"].(float64)))
	assert.Equal(t, http.StatusText(http.StatusUnauthorized), responseBody["status"])
}
