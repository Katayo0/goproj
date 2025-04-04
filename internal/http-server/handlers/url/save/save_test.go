package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go_sql_test/internal/http-server/handlers/url/save"
	"go_sql_test/internal/http-server/handlers/url/save/mocks"
	"go_sql_test/internal/lib/logger/handlers/slogdiscard"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

	

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
	  c.String(200, "pong")
	})
	return router
  }

  func postUser(router *gin.Engine, slog *slog.Logger, urlSaver save.URLSaver) *gin.Engine {
	
	router.POST("/user/add", save.New(slog, urlSaver))
	return router
  }




func TestSaveHandler(t *testing.T){
	cases := []struct {
		name string
		alias string
		url string
		respError string
		mockError error
		}{
		{
			name: "Success",
			alias: "test_alias",
			url: "https://google.com",
		},
		{
			name: "Empty Alias",
			alias: "",
			url: "https://google.com",
		},
		{
			name: "Empty url",
			url: "",
			alias: "some_alias",
			respError: "URL is empty",
		},
		{
			name: "Invalid url",
			url: "some invalid url",
			alias: "some_alias",
			respError: "url is not valid",
		},
		{
			name: "SaveURL error",
			alias: "test_alias",
			url: "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},

	}
	
	



	for _, tc := range cases{
		tc := tc

		t.Run(tc.name, func(t *testing.T){
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil{
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int(1), tc.mockError).
					Once()

			}
			
			//handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			handler := setupRouter()
			handler = postUser(handler, slogdiscard.NewDiscardLogger(), urlSaverMock)


			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

			
		})
	}
}

