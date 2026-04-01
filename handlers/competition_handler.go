package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liyufei916/footballPaul/services"
)

type CompetitionHandler struct {
	competitionService *services.CompetitionService
}

func NewCompetitionHandler() *CompetitionHandler {
	return &CompetitionHandler{
		competitionService: services.NewCompetitionService(),
	}
}

func (h *CompetitionHandler) GetCompetitions(c *gin.Context) {
	competitions, err := h.competitionService.GetAllCompetitions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"competitions": competitions,
		"count":        len(competitions),
	})
}

func (h *CompetitionHandler) GetCompetition(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid competition id"})
		return
	}

	competition, err := h.competitionService.GetCompetitionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, competition)
}

type CreateCompetitionRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
	Code string `json:"code" binding:"required,min=1,max=20"`
	Logo string `json:"logo"`
}

func (h *CompetitionHandler) CreateCompetition(c *gin.Context) {
	var req CreateCompetitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	competition, err := h.competitionService.CreateCompetition(req.Name, req.Code, req.Logo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":     true,
		"competition": competition,
	})
}
