package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/leboncoin/subot/pkg/globals"
)

// GetTools godoc
// @Summary Get all configured tools
// @Description TODO
// @Tags Tools
// @ID get-tools
// @Produce  json
// @Router /tools [get]
func (a Analyser) GetTools(c *gin.Context) {
	answers, err := a.ESClient.GetTools()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, answers)
}

// AddTool godoc
// @Summary Add new tool
// @Description Saves the new tool to the database.
// @Description Authentication and admin access are required for this endpoint
// @Tags Tools
// @ID add-tool
// @Produce  json
// @Param name body string true "Name of the tool to add"
// @Param query body object true "The percolator query it shall match"
// @Router /tools/new [post]
func (a Analyser) AddTool(c *gin.Context) {
	var eventRequest globals.Perco
	if err := c.BindJSON(&eventRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := a.ESClient.AddTool(eventRequest); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{})
}

// EditTool godoc
// @Summary Edit existing tool
// @Description Modifies the name or the query (or both)
// @Description of the specified tool in the database
// @Description Authentication and admin access are required for this endpoint
// @Tags Tools
// @ID edit-tool
// @Produce  json
// @Param tool query string true "Tool id to update"
// @Param name body string true "Name of the tool to add"
// @Param query body object true "The percolator query it shall match"
// @Router /tools/:tool [put]
func (a Analyser) EditTool(c *gin.Context) {
	tool := c.Param("tool")
	var eventRequest globals.Perco
	if err := c.BindJSON(&eventRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := a.ESClient.EditTool(tool, eventRequest); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{})
}

// DeleteTool godoc
// @Summary Delete existing tool
// @Description Removes the tool from the database
// @Description The tool to delete must not be used by any answer
// @Description Authentication and admin access are required for this endpoint
// @Tags Tools
// @ID delete-tool
// @Produce  json
// @Param tool query string true "Tool id to delete"
// @Router /tools/:tool [delete]
func (a Analyser) DeleteTool(c *gin.Context) {
	tool := c.Param("tool")
	if err := a.ESClient.DeleteTool(tool); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(204, gin.H{})
}
