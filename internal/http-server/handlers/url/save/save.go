package save

import (
	"encoding/json"
	"errors"
	resplib "go_sql_test/internal/lib/api/response"
	"go_sql_test/internal/lib/logger/sl"
	"go_sql_test/internal/lib/random"
	"go_sql_test/internal/storage"
	//"io"
	"log/slog"
	//"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/render"
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

//TODO: move to config
const aliasLength = 6

type URLSaver interface{
	SaveURL(urlToSave string, alias string) (int, error)
}

	
func New(log *slog.Logger, urlSaver URLSaver) gin.HandlerFunc {
	return func(c *gin.Context){
		const fnNm = "handlers.url.save.New"
		requestid.RequestID(nil)
		log = log.With(
			slog.String("fnNm", fnNm),
			slog.String("request_id", requestid.GetRequestIDFromContext(c)),
		)

		var req Request
	
		dcded, derr := c.GetRawData()
		if dcded == nil{
			log.Error("request body is empty", sl.Err(derr))
			jsonresp, _ := json.Marshal(resplib.Error("request body is empty"))
			c.JSON(200, gin.H{
				"request":dcded,
				"error_code":derr,
				"error":jsonresp,
			})
			return
		}

		//is it needed?
		if derr != nil{
			log.Error("failed to get data from context", sl.Err(derr))
			jsonresp, _ := json.Marshal(resplib.Error("failed to get raw data"))
			c.JSON(200, gin.H{
				"request":dcded,
				"error_code":derr,
				"error":jsonresp,
			})
			return
		}
		//end of controversial code

		log.Info("Got raw data from context", slog.Any("RawData", dcded))

		
		err := json.Unmarshal(dcded, &req)
		if errors.Is(err, &json.InvalidUnmarshalError{}){
			log.Error("request body is empty", sl.Err(err))
			jsonresp, _ := json.Marshal(resplib.Error("failed to decode request"))
			c.JSON(200, gin.H{
				"request":dcded,
				"error_code":err,
				"error":jsonresp,
			})

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			jsonresp, _ := json.Marshal(resplib.Error("failed to decode request"))
			c.JSON(200, gin.H{
				"request":dcded,
				"error_code":err,
				"error":jsonresp,
			})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))


		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			jsonresp, _:= json.Marshal(resplib.ValidationError(validateErr))
			c.JSON(200, gin.H{
				"request":dcded,
				"error_code":err,
				"error":jsonresp,
			})

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		//make anti-collision mechanism 


		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrUrlExists){
			log.Info("url already exists", slog.String("url", req.URL))
			jsonresp, _ := json.Marshal(resplib.Error("url already exists"))
			c.JSON(200, gin.H{
				"request":   dcded,
				"error_code":err,
				"error":     jsonresp,
			})
			return
		}
		if err != nil{
			log.Info("failed to add url", sl.Err(err))
			jsonresp, _ := json.Marshal(resplib.Error("failed to add url"))
			c.JSON(200, gin.H{
				"request":   dcded,
				"error_code":err,
				"error":     jsonresp,
			})
		}

		log.Info("url added", slog.Int("id", id))
		c.JSON(200, gin.H{
			"response":resplib.OK(),
			"Alias":   alias,
		})

	}

	}

