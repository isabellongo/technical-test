package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/driva/api/internal/repository"
)

// PeopleHandler trata endpoints /people/v1/*.
type PeopleHandler struct {
	seed *repository.SeedRepository
}

// NewPeopleHandler cria um novo handler.
func NewPeopleHandler(seed *repository.SeedRepository) *PeopleHandler {
	return &PeopleHandler{seed: seed}
}

// GetEnrichments retorna enrichments paginados da fonte (api_enrichments_seed).
func (h *PeopleHandler) GetEnrichments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	result, err := h.seed.ListPaginated(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
