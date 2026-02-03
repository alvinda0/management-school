package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"project-go/internal/model"
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
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleRepo.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Role not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Role retrieved successfully", role)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	role := &model.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.roleRepo.Create(role); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Role created successfully", role)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get existing role
	role, err := h.roleRepo.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Role not found")
		return
	}

	// Update fields
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := h.roleRepo.Update(role); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Role updated successfully", role)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	// Check if role exists
	_, err := h.roleRepo.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Role not found")
		return
	}

	if err := h.roleRepo.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Role deleted successfully", nil)
}