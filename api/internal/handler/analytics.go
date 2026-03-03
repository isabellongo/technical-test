package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/driva/api/internal/repository"
)

// AnalyticsHandler trata endpoints /analytics/*.
type AnalyticsHandler struct {
	gold *repository.GoldRepository
}

// NewAnalyticsHandler cria um novo handler.
func NewAnalyticsHandler(gold *repository.GoldRepository) *AnalyticsHandler {
	return &AnalyticsHandler{gold: gold}
}

// GetOverview retorna KPIs agregados da Gold.
func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	result, err := h.gold.Overview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetEnrichments retorna Gold enrichments paginados com filtros opcionais.
func (h *AnalyticsHandler) GetEnrichments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	status := c.Query("status_processamento")
	categoria := c.Query("categoria_tamanho_job")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	result, err := h.gold.ListPaginated(c.Request.Context(), page, limit, status, categoria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
