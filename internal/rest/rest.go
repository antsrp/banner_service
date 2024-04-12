package rest

import (
	"fmt"
	"net/http"

	//"github.com/antsrp/warehouse/internal/models"
	rs "github.com/antsrp/banner_service/pkg/infrastructure/rest"
	"github.com/antsrp/banner_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	statusCode int
	realError  error
	Message    string `json:"message"`
}

func (r ResponseError) Error() string {
	return r.Message
}

type Handler struct {
	engine   *gin.Engine
	settings rs.Settings
	logger   logger.Logger
}

func NewHandler(settings rs.Settings, logger logger.Logger) Handler {
	h := Handler{
		engine:   gin.Default(),
		settings: settings,
		logger:   logger,
	}
	h.engine.Use(h.errorHandler())
	h.routes()
	return h
}

func (h Handler) errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case ResponseError:
				c.AbortWithStatusJSON(e.statusCode, e)
				h.logger.Info(fmt.Sprintf("error within request: %v", e.realError.Error()))
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseError{Message: "Service Unavailable"})
				h.logger.Error(fmt.Sprintf("internal error: %v", e.Error()))
			}
		}
	}
}

func (h Handler) routes() {
	h.engine.GET("/user_banner", h.userBanner)
	h.engine.GET("/banner", h.getBanner)
	h.engine.POST("/banner", h.addBanner)
	h.engine.PATCH("/banner/:id", h.updateBanner)
	h.engine.DELETE("/banner/:id", h.deleteBanner)
}

func (h Handler) Run() error {
	if err := h.engine.Run(fmt.Sprintf("%s:%s", h.settings.Host, h.settings.Port)); err != nil {
		return fmt.Errorf("can't run server: %w", err)
	}
	return nil
}

/*
summary: Получение баннера для пользователя

	parameters:
	  - in: query
	    name: tag_id
	    required: true
	    schema:
	      type: integer
	      description: Тэг пользователя
	  - in: query
	    name: feature_id
	    required: true
	    schema:
	      type: integer
	      description: Идентификатор фичи
	  - in: query
	    name: use_last_revision
	    required: false
	    schema:
	      type: boolean
	      default: false
	      description: Получать актуальную информацию
	  - in: header
	    name: token
	    description: Токен пользователя
	    schema:
	      type: string
	      example: "user_token"
	responses:
	  '200':
	    description: Баннер пользователя
	    content:
	      application/json:
	        schema:
	          description: JSON-отображение баннера
	          type: object
	          additionalProperties: true
	          example: '{"title": "some_title", "text": "some_text", "url": "some_url"}'
	  '400':
	    description: Некорректные данные
	    content:
	      application/json:
	        schema:
	          type: object
	          properties:
	            error:
	              type: string
	  '401':
	    description: Пользователь не авторизован
	  '403':
	    description: Пользователь не имеет доступа
	  '404':
	    description: Баннер для не найден
	  '500':
	    description: Внутренняя ошибка сервера
	    content:
	      application/json:
	        schema:
	          type: object
	          properties:
	            error:
	              type: string
*/
func (h Handler) userBanner(c *gin.Context) { // GET /user_banner
	_ = c.GetHeader("token") // parse token
	c.GetQuery("tag_id")
	c.GetQuery("feature_id")
	c.GetQuery("use_last_revision") // not required
}

/*
summary: Получение всех баннеров c фильтрацией по фиче и/или тегу

	parameters:
	  - in: header
	    name: token
	    description: Токен админа
	    schema:
	      type: string
	      example: "admin_token"
	  - in: query
	    name: feature_id
	    required: false
	    schema:
	      type: integer
	      description: Идентификатор фичи
	  - in: query
	    name: tag_id
	    required: false
	    schema:
	      type: integer
	      description: Идентификатор тега
	  - in: query
	    name: limit
	    required: false
	    schema:
	      type: integer
	      description: Лимит
	  - in: query
	    name: offset
	    required: false
	    schema:
	      type: integer
	      description: Оффсет
	responses:
	  '200':
	    description: OK
	    content:
	      application/json:
	        schema:
	          type: array
	          items:
	            type: object
	            properties:
	              banner_id:
	                type: integer
	                description: Идентификатор баннера
	              tag_ids:
	                type: array
	                description: Идентификаторы тэгов
	                items:
	                  type: integer
	              feature_id:
	                type: integer
	                description: Идентификатор фичи
	              content:
	                type: object
	                description: Содержимое баннера
	                additionalProperties: true
	                example: '{"title": "some_title", "text": "some_text", "url": "some_url"}'
	              is_active:
	                type: boolean
	                description: Флаг активности баннера
	              created_at:
	                type: string
	                format: date-time
	                description: Дата создания баннера
	              updated_at:
	                type: string
	                format: date-time
	                description: Дата обновления баннера
	  '401':
	    description: Пользователь не авторизован
	  '403':
	    description: Пользователь не имеет доступа
	  '500':
	    description: Внутренняя ошибка сервера
	    content:
	      application/json:
	        schema:
	          type: object
	          properties:
	            error:
	              type: string
*/
func (h Handler) getBanner(c *gin.Context) { // GET /banner

}

/*
summary: Создание нового баннера

	parameters:
	  - in: header
	    name: token
	    description: Токен админа
	    schema:
	      type: string
	      example: "admin_token"
	requestBody:
	  required: true
	  content:
	    application/json:
	      schema:
	        type: object
	        properties:
	          tag_ids:
	            type: array
	            description: Идентификаторы тэгов
	            items:
	              type: integer
	          feature_id:
	            type: integer
	            description: Идентификатор фичи
	          content:
	            type: object
	            description: Содержимое баннера
	            additionalProperties: true
	            example: '{"title": "some_title", "text": "some_text", "url": "some_url"}'
	          is_active:
	            type: boolean
	            description: Флаг активности баннера
	responses:
	  '201':
	    description: Created
	    content:
	      application/json:
	        schema:
	          type: object
	          properties:
	            banner_id:
	              type: integer
	              description: Идентификатор созданного баннера
	  '400':
	    description: Некорректные данные
	    content:
	      application/json:
	        schema:
	          type: object
	          properties:
	            error:
	              type: string
	  '401':
	    description: Пользователь не авторизован
	  '403':
	    description: Пользователь не имеет доступа
	  '500':
	    description: Внутренняя ошибка сервера
	    content:
	      application/json:
	        schema:
	          type: object
	          properties:
	            error:
	              type: string
*/
func (h Handler) addBanner(c *gin.Context) { // POST /banner

}

/*
summary: Обновление содержимого баннера

	parameters:
	  - in: path
	    name: id
	    required: true
	    schema:
	      type: integer
	      description: Идентификатор баннера
	  - in: header
	    name: token
	    description: Токен админа
	    schema:
	      type: string
	      example: "admin_token"
	requestBody:
	  required: true
	  content:
	    application/json:
	      schema:
	        type: object
	        properties:
	          tag_ids:
	            nullable: true
	            type: array
	            description: Идентификаторы тэгов
	            items:
	              type: integer
	          feature_id:
	            nullable: true
	            type: integer
	            description: Идентификатор фичи
	          content:
	            nullable: true
	            type: object
	            description: Содержимое баннера
	            additionalProperties: true
	            example: '{"title": "some_title", "text": "some_text", "url": "some_url"}'
	          is_active:
	            nullable: true
	            type: boolean
	            description: Флаг активности баннера
	responses:
	  '200':
	    description: OK
	  '400':
	    description: Некорректные данные
	    content:
	      application/json:
	        schema:
	          type: object
	          properties:
	            error:
	              type: string
	  '401':
	    description: Пользователь не авторизован
	  '403':
	    description: Пользователь не имеет доступа
	  '404':
	    description: Баннер не найден
	  '500':
	    description: Внутренняя ошибка сервера
	    content:
	      application/json:
	        schema:
	          type: object
	          properties:
	            error:
	              type: string
*/
func (h Handler) updateBanner(c *gin.Context) { // PATCH /banner/{id}

}

/*
summary: Удаление баннера по идентификатору

	parameters:
	  - in: path
	    name: id
	    required: true
	    schema:
	      type: integer
	      description: Идентификатор баннера
	  - in: header
	    name: token
	    description: Токен админа
	    schema:
	      type: string
	      example: "admin_token"
	responses:
	  '204':
	    description: Баннер успешно удален
	  '400':
	    description: Некорректные данные
	    content:
	      application/json:
	        schema:
	          type: object
	          properties:
	            error:
	              type: string
	  '401':
	    description: Пользователь не авторизован
	  '403':
	    description: Пользователь не имеет доступа
	  '404':
	    description: Баннер для тэга не найден
	  '500':
	    description: Внутренняя ошибка сервера
	    content:
	      application/json:
	        schema:
	          type: object
	          properties:
	            error:
	              type: string
*/
func (h Handler) deleteBanner(c *gin.Context) { // DELETE /banner/{id}

}

/*func (h Handler) reserveGoods(c *gin.Context) { // POST ReserveGoods ([]models.Reserve)
	var reserves []models.Reserve
	if err := c.ShouldBindJSON(&reserves); err != nil {
		c.Error(ResponseError{realError: err, statusCode: http.StatusBadRequest, Message: "bad request: input should be json"})
		return
	}
	c.JSON(http.StatusOK, reserves)
}

func (h Handler) releaseGoods(c *gin.Context) { // POST ReleaseGoods ([]models.Release)
	var releases []models.Release
	if err := c.ShouldBindJSON(&releases); err != nil {
		c.Error(ResponseError{realError: err, statusCode: http.StatusBadRequest, Message: "bad request: input should be json"})
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h Handler) getGoodsOfWarehouse(c *gin.Context) { // POST GetGoodsOfWarehouse ([]models.GetGoodsRequest)
	var ggrs []models.GetGoodsRequest
	if err := c.ShouldBindJSON(&ggrs); err != nil {
		c.Error(ResponseError{realError: err, statusCode: http.StatusBadRequest, Message: "bad request: input should be json"})
		return
	}
	c.JSON(http.StatusOK, nil)
} */
