package jsonapi

// import (
// 	"encoding/json"
// 	"github.com/genya0407/confession-server/usecase"
// 	"github.com/google/uuid"
// 	"github.com/julienschmidt/httprouter"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )

// func generateMockUsecaseAccountInfoDTO(name string, expectID uuid.UUID) usecase.GetAccountInfo {
// 	return func(id uuid.UUID) (usecase.AccountInfoDTO, bool) {
// 		if id != expectID {
// 			return usecase.AccountInfoDTO{}, false
// 		}
// 		return usecase.AccountInfoDTO{
// 			AccountID: expectID,
// 			Name:      name,
// 		}, true
// 	}
// }

// func generateMockGetAccountInfoRouter(mockname string, accountID uuid.UUID) *httprouter.Router {
// 	uc := generateMockUsecaseAccountInfoDTO(mockname, accountID)
// 	handler := GenerateGetAccountInfo(uc)
// 	router := httprouter.New()
// 	router.GET("/account/:account_id", handler)
// 	return router
// }

// func TestGetAccountInfo(t *testing.T) {
// 	mockname := "Mock name"
// 	accountID, err := uuid.NewUUID()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	router := generateMockGetAccountInfoRouter(mockname, accountID)

// 	req := httptest.NewRequest("GET", "http://confession.com/account/"+accountID.String(), nil)
// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)
// 	if w.Code != http.StatusOK {
// 		t.Error("Invalid status")
// 	}

// 	result := &AccountJSON{}
// 	err = json.NewDecoder(w.Body).Decode(result)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	if result.AccountID != accountID {
// 		t.Error("Invalid result id")
// 	}
// 	if result.Name != mockname {
// 		t.Error("Invalid result name")
// 	}
// }

// func TestGetAccountInfoNotFound(t *testing.T) {
// 	mockname := "Mock name"
// 	accountID, err := uuid.NewUUID()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	router := generateMockGetAccountInfoRouter(mockname, accountID)

// 	notExistAccountID, err := uuid.NewUUID()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	req := httptest.NewRequest("GET", "http://confession.com/account/"+notExistAccountID.String(), nil)
// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)
// 	if w.Code != http.StatusNotFound {
// 		t.Error("Invalid status")
// 	}
// }
