package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/antsrp/banner_service/internal/domain/models"
	"github.com/antsrp/banner_service/internal/domain/models/requests"
	"github.com/antsrp/banner_service/internal/service"
	rs "github.com/antsrp/banner_service/pkg/infrastructure/rest"
	"github.com/antsrp/banner_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine        *gin.Engine
	settings      rs.Settings
	logger        logger.Logger
	bannerService service.BannerServicer
	auth          authHandler
}

func NewHandler(settings rs.Settings, logger logger.Logger, bs service.BannerServicer, us service.UserStorager) Handler {
	h := Handler{
		engine:        gin.Default(),
		settings:      settings,
		logger:        logger,
		auth:          newAuthHandler(us, logger),
		bannerService: bs,
	}
	h.routes()
	return h
}

func (h Handler) routes() {
	group := h.engine.Group("/")
	group.GET("/user_banner", h.auth.authRequired, h.userBanner)
	group.GET("/banner", h.auth.adminAuthRequired, h.getBanner)
	group.POST("/banner", h.auth.adminAuthRequired, h.addBanner)
	group.PATCH("/banner/:id", h.auth.adminAuthRequired, h.updateBanner)
	group.DELETE("/banner/:id", h.auth.adminAuthRequired, h.deleteBanner)

	group.POST("/signin", h.auth.signIn)
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
	var req requests.UserBannerRequest
	if tagId, ok := c.GetQuery("tag_id"); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "tag id is not set"})
		return
	} else {
		if val, err := strconv.Atoi(tagId); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "tag id is not an integer type"})
			return
		} else {
			req.TagID = val
		}
	}
	if featureId, ok := c.GetQuery("feature_id"); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "feature id is not set"})
		return
	} else {
		if val, err := strconv.Atoi(featureId); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "feature id is not an integer type"})
			return
		} else {
			req.FeatureID = val
		}
	}
	if useLR, ok := c.GetQuery("use_last_revision"); ok { // not required
		if val, err := strconv.ParseBool(useLR); err != nil {
			h.logger.Info("can't parse use_last_revision parameter from %s: not a boolean type", useLR)
		} else {
			req.IsUseLastRevision = val
		}
	}
	data, _ := c.Get(authusertag)
	user := data.(models.User)

	banner, err := h.bannerService.GetOne(req, user.Name)
	if err != nil {
		h.logger.Error("can't get banner: %v", err.Cause().Error())
		if errors.Is(err.Cause(), service.ErrBannerNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": service.ErrDefaultInternalError.Error()})
		}
		return
	}
	if !(*banner.IsActive) && !user.IsAdmin {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.JSON(http.StatusOK, banner)
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
	var req requests.GetBannersRequest

	if tagId, ok := c.GetQuery("tag_id"); ok {
		if val, err := strconv.Atoi(tagId); err != nil {
			h.logger.Info("can't parse tag_id: %v", err.Error())
		} else {
			req.TagID = val
		}
	}
	if featureId, ok := c.GetQuery("feature_id"); ok {
		if val, err := strconv.Atoi(featureId); err != nil {
			h.logger.Info("can't parse feature_id: %v", err.Error())
		} else {
			req.FeatureID = val
		}
	}
	if limit, ok := c.GetQuery("limit"); ok {
		if val, err := strconv.Atoi(limit); err != nil {
			h.logger.Info("can't parse limit: %v", err.Error())
		} else {
			req.Limit = val
		}
	}
	if offset, ok := c.GetQuery("offset"); ok {
		if val, err := strconv.Atoi(offset); err != nil {
			h.logger.Info("can't parse limit: %v", err.Error())
		} else {
			req.Offset = val
		}
	}
	banners, err := h.bannerService.Get(req)
	if err != nil {
		h.logger.Error("can't get banners: %v", err.Cause().Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, service.ErrDefaultInternalError.Error())
		return
	}

	c.JSON(http.StatusOK, banners)
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
	var req requests.CreateBannerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("can't parse request body from json: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, service.ErrDefaultInternalError.Error())
		return
	}
	if req.Content == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "field content is empty"})
		return
	}
	if req.IsActive == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "field is_active is not set"})
		return
	}
	if req.FeatureID <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "field feature_id is not set or set wrong"})
		return
	}
	if req.TagIDS == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "field tag_ids is not set"})
		return
	}

	banner, err := h.bannerService.Create(req)
	if err != nil {
		h.logger.Error("can't parse request body from json: %v", err.Cause().Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, service.ErrDefaultInternalError.Error())
		return
	}

	c.JSON(http.StatusCreated, requests.CreateBannerResponse{BannerID: banner.ID})
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
	var req requests.UpdateBannerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("can't parse request body from json: %v", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, service.ErrDefaultInternalError.Error())
		return
	}
	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id parameter is not an integer type"})
		return
	} else {
		req.ID = id
	}

	err := h.bannerService.Update(req)
	if err != nil {
		h.logger.Error("can't update banner in database: %v", err.Cause().Error())
		if errors.Is(err.Cause(), service.ErrBannerNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, service.ErrDefaultInternalError.Error())
		}
		return
	}

	c.Status(http.StatusOK)
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
	var req requests.DeleteBannerRequest

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id parameter is not an integer type"})
		return
	}
	req.ID = id

	if err := h.bannerService.Delete(req); err != nil {
		if errors.Is(err.Cause(), service.ErrBannerNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, service.ErrDefaultInternalError.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}
