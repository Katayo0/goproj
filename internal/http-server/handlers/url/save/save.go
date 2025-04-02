package save

import (
	"encoding/json"
	"errors"
	resplib "go_sql_test/internal/lib/api/response"
	"go_sql_test/internal/lib/logger/sl"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/go-playground/validator/v10"
	requestid "github.com/sumit-tembe/gin-requestid"
)

type Request struct{
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct{
	resplib.Response
	Alias  string `json:"alias,omitempty"`
}

type URLSaver interface{
	SaveURL(urlToSave string, alias string) (int, error)
}

/*
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		const fnNm = "handlers.url.save.New"
		log = log.With(
			slog.String("fnNm", fnNm),
			// joskiy kaloviy kod dlya requestid
			slog.String("request_id", requestid.GetRequestIDFromContext(&gin.Context{})),
		)
		var req Request

		err := gin.Dec
	}
}
*/	
func New(log *slog.Logger, urlSaver URLSaver) gin.HandlerFunc {
	return func(c *gin.Context){
		const fnNm = "handlers.url.save.New"
		log = log.With(
			slog.String("fnNm", fnNm),
			// joskiy kaloviy kod dlya requestid
			slog.String("request_id", requestid.GetRequestIDFromContext(c)),
		)

		var req Request
		
		//nihuya ne ponyal po http i gin.context NADO razobratsya
		//etot loger skoree vsego nahui ne rabotayet

		//1:20:48

		//UPD logger gotov, no kuchiy
		dcded, derr := c.GetRawData()
		if derr != nil{
			log.Error("failed to get data from context", sl.Err(derr))
			json.Marshal(resplib.Error("failed to get raw data"))
		}

		log.Info("Got raw data from context", slog.Any("RawData", dcded))

		err := json.Unmarshal(dcded, &req)
		if errors.Is(err, &json.InvalidUnmarshalError{}){
			log.Error("request body is empty", sl.Err(err))
			json.Marshal(resplib.Error("failed to decode request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			json.Marshal(resplib.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			json.Marshal(resplib.Error("invalid request"))
		}
	}

	}

