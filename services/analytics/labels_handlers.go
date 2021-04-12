package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/leboncoin/subot/pkg/globals"
)

// GetLabels godoc
// @Summary Get all configured labels
// @Description Returns all the labels from the database.
// @Description No authentication required
// @Tags Labels
// @ID get-labels
// @Produce  json
// @Router /labels [get]
func (a Analyser) GetLabels(c *gin.Context) {
	hits, err := a.ESClient.GetLabels()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, hits)
}

// AddLabel godoc
// @Summary Create new label
// @Description Stores a new label into the database with the given information
// @Description Authentication and admin access are required for this endpoint
// @Tags Labels
// @ID add-label
// @Produce  json
// @Param name body string true "Name of the label to add"
// @Param query body object true "The percolator query it shall match"
// @Router /labels/new [post]
func (a Analyser) AddLabel(c *gin.Context) {
	var eventRequest globals.Perco
	if err := c.BindJSON(&eventRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := a.ESClient.AddLabel(eventRequest); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{})
}

// EditLabel godoc
// @Summary Edit specified label
// @Description Modify the data stored in the database
// @Description for the document having the given label ID.
// @Description Authentication and admin access are required for this endpoint
// @Tags Labels
// @ID edit-label
// @Produce  json
// @Param label query string true "Label id to update"
// @Param name body string true "Name of the label to add"
// @Param query body object true "The percolator query it shall match"
// @Router /labels/:label [put]
func (a Analyser) EditLabel(c *gin.Context) {
	label := c.Param("label")
	var eventRequest globals.Perco
	if err := c.BindJSON(&eventRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := a.ESClient.EditLabel(label, eventRequest); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{})
}

// DeleteLabel godoc
// @Summary Delete specified label
// @Description Deletes the entry matching the ID if not used.
// @Description Authentication and admin access are required for this endpoint
// @Tags Labels
// @ID delete-label
// @Produce  json
// @Param label query string true "Label id to update"
// @Router /labels/:label [delete]
func (a Analyser) DeleteLabel(c *gin.Context) {
	label := c.Param("label")
	if err := a.ESClient.DeleteLabel(label); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(204, gin.H{})
}
