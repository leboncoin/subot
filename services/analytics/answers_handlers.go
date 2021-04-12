package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/leboncoin/subot/pkg/globals"
)

// GetAnswers godoc
// @Summary Get all configured answers
// @Description Returns the list of existing answers.
// @Description No authentication required.
// @Tags Answers
// @ID get-answers
// @Produce  json
// @Router /answers [get]
func (a Analyser) GetAnswers(ctx *gin.Context) {
	answers, err := a.ESClient.GetAnswers()
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, answers)
}

// AddAnswer godoc
// @Summary Add a answer
// @Description Saves the new answer to the database,
// @Description ensuring the specified tools and labels exist.
// @Description Authentication and admin access are required for this endpoint
// @Tags Answers
// @ID add-answer
// @Produce  json
// @Param tool body string false "Tool to match for this answer"
// @Param label body string false "Label to match for this answer"
// @Param answer body string true "The answer to reply when matched"
// @Param feedback body bool true "Whether or not the bot shall ask for a user feedback"
// @Router /answers/new [post]
func (a Analyser) AddAnswer(c *gin.Context) {
	var eventRequest globals.Answer
	if err := c.BindJSON(&eventRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := a.ESClient.AddAnswer(eventRequest); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{})
}

// EditAnswer godoc
// @Summary Modify the answer
// @Description Updates the document in the database that matches
// @Description the given documentID and stores the new information.
// @Description Authentication and admin access are required for this endpoint
// @Tags Answers
// @ID edit-answer
// @Produce  json
// @Param documentID query string true "Answer id to update"
// @Param tool body string false "Tool to match for this answer"
// @Param label body string false "Label to match for this answer"
// @Param answer body string true "The answer to reply when matched"
// @Param feedback body bool true "Whether or not the bot shall ask for a user feedback"
// @Router /answers/:documentID [put]
func (a Analyser) EditAnswer(c *gin.Context) {
	var eventRequest globals.Answer
	documentID := c.Param("documentID")
	if err := c.BindJSON(&eventRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := a.ESClient.EditAnswer(documentID, eventRequest); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{})
}

// DeleteAnswer godoc
// @Summary Delete specified answer
// @Description Removes the entry at the given documentID from the database.
// @Description Authentication and admin access are required for this endpoint
// @Tags Answers
// @ID delete-answer
// @Produce  json
// @Param documentID query string true "Answer id to delete"
// @Router /answers/:documentID [delete]
func (a Analyser) DeleteAnswer(c *gin.Context) {
	documentID := c.Param("documentID")
	if err := a.ESClient.DeleteAnswer(documentID); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(204, gin.H{})
}
