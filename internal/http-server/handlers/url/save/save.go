package save

import (
	//"encoding/json"
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
//go:generate go run github.com/vektra/mockery/v2@v2.53.3 --name=URLSaver
type URLSaver interface{
	SaveURL(urlToSave string, alias string) (int, error)
}

	
func New(log *slog.Logger, urlSaver URLSaver) gin.HandlerFunc {
	return func(c *gin.Context){
		const fnNm = "handlers.url.save.New"
		log = log.With(
			slog.String("fnNm", fnNm),
			slog.String("request_id", requestid.GetRequestIDFromContext(c)),
		)

		var req Request
	
		err := c.Bind(&req)
		if err != nil{
			log.Error("unable to bind body to struct")
			errresp := resplib.Error("request body is empty")
			c.JSON(200, gin.H{
				"request":     req,
				"error_code" : err,
				"error":       errresp,
			})
			return
		}
		

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			errresp:= resplib.ValidationError(validateErr)
			c.JSON(200, gin.H{
				"request":    req,
				"error_code": err,
				"error":      errresp,
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
			errresp := resplib.Error("url already exists")
			c.JSON(200, gin.H{
				"request":    req,
				"error_code": err,
				"error":      errresp,
			})
			return
		}
		if err != nil{
			log.Info("failed to add url", sl.Err(err))
			errresp := resplib.Error("failed to add url")
			c.JSON(200, gin.H{
				"request":    req,
				"error_code": err,
				"error":      errresp,
			})
		}

		log.Info("url added", slog.Int("id", id))
		c.JSON(200, gin.H{
			"response":resplib.OK(),
			"Alias":   alias,
		})

	}

	}

