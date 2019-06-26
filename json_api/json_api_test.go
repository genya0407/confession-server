package json_api

import (
	"encoding/json"
	"github.com/genya0407/confession-server/usecase"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GenerateMockUsecaseAccountInfoDTO(name string, expectID uuid.UUID) usecase.GetAccountInfo {
	return func(id uuid.UUID) (usecase.AccountInfoDTO, bool) {
		if id != expectID {
			return usecase.AccountInfoDTO{}, false
		}
		return usecase.AccountInfoDTO{
			AccountID: expectID,
			Name:      name,
		}, true
	}
}

func TestGetAccountInfo(t *testing.T) {
	mockname := "Mock name"
	accountID, err := uuid.NewUUID()
	uc := GenerateMockUsecaseAccountInfoDTO(mockname, accountID)

	if err != nil {
		panic(err.Error())
	}
	handler := GetAccountInfoGenerator(uc)
	router := httprouter.New()
	router.GET("/account/:account_id", handler)

	req := httptest.NewRequest("GET", "http://confession.com/account/"+accountID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("Invalid status")
	}

	result := &AccountJSON{}
	err = json.NewDecoder(w.Body).Decode(result)
	if err != nil {
		t.Error(err.Error())
	}

	if result.AccountID != accountID {
		t.Error("Invalid result id")
	}
	if result.Name != mockname {
		t.Error("Invalid result name")
	}
}
