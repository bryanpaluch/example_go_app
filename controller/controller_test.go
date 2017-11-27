package controller

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bryanpaluch/example_go_app/example"
	"github.com/bryanpaluch/example_go_app/mocks"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func makeRequest(handler http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func NewSetup(ctrl *gomock.Controller) (*Router, *mocks.MockDB) {
	mockDB := mocks.NewMockDB(ctrl)
	router, _ := NewRouter(mockDB)
	return router, mockDB
}

func TestGetPerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Person Not returned by database
	{
		router, mockDB := NewSetup(ctrl)
		mockDB.EXPECT().GetPersonByID(gomock.Any(), 1).Times(1).Return(nil, nil)
		req := httptest.NewRequest("GET", "/person/1", nil)
		res := makeRequest(router, req)
		if res.Code != 404 {
			t.Fail()
		}
	}

	// ID not a number
	{
		router, _ := NewSetup(ctrl)
		req := httptest.NewRequest("GET", "/person/test", nil)
		res := makeRequest(router, req)
		if res.Code != 400 {
			t.Fail()
		}
	}

	// Error returned from database
	{
		router, mockDB := NewSetup(ctrl)
		mockDB.EXPECT().GetPersonByID(gomock.Any(), 1).Times(1).Return(nil, fmt.Errorf("some error"))
		req := httptest.NewRequest("GET", "/person/1", nil)
		res := makeRequest(router, req)
		if res.Code != 500 {
			t.Fail()
		}
	}

	// Person returned 200 OK
	{
		router, mockDB := NewSetup(ctrl)
		d := time.Date(1985, 10, 4, 23, 44, 20, 0, time.Local)
		p := &example.Person{ID: 3, Name: "Bryan", Birth: &d}
		mockDB.EXPECT().GetPersonByID(gomock.Any(), 1).Times(1).Return(p, nil)
		req := httptest.NewRequest("GET", "/person/1", nil)
		res := makeRequest(router, req)
		if res.Code != 200 {
			t.Fail()
		}
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println("body is ", string(body))

	}
}

func TestInsertPerson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Person Add
	{
		router, mockDB := NewSetup(ctrl)
		mockDB.EXPECT().AddPerson(gomock.Any(), gomock.Any()).Times(1).Return(nil).
			Do(func(c context.Context, p *example.Person) {
				if p.Name != "bryan" {
					t.Fail()
				}
				p.ID = 3
			})

		personJson := []byte(`{"name":"bryan"}`)
		req := httptest.NewRequest("POST", "/person", bytes.NewBuffer(personJson))
		res := makeRequest(router, req)
		if res.Code != 200 {
			t.Fail()
		}
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println("body is ", string(body))
	}
}
