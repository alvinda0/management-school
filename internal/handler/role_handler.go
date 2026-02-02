package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"project-go/internal/repository"
	"project-go/utils"
)

type RoleHandler struct {
	roleRepo repository.RoleRepository
}

func NewRoleHandler(roleRepo repository.RoleRepository) *RoleHandler {
	return &RoleHandler{
		roleRepo: roleRepo,
	}
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleRepo.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Roles retrieved successfully", roles)
}

func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleRepo.GetWithPermissions(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Role not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Role retrieved successfully", role)
}