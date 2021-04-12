package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/leboncoin/subot/pkg/globals"
)

// GetTeamMembers godoc
// @Summary Get all configured team members
// @Description TODO
// @Tags Team
// @ID get-team-members
// @Produce  json
// @Router /team [get]
func (a Analyser) GetTeamMembers(c *gin.Context) {
	answers, err := a.ESClient.GetTeamMembers()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, answers)
}

// AddTeamMember godoc
// @Summary Add a new member to the team
// @Description Saves the new team member to the database.
// @Description Authentication and admin access are required for this endpoint
// @Tags Team
// @ID add-team-member
// @Produce  json
// @Param slack_id body object true "ID of the user in slack"
// @Param name body object true "The name of the user only used in frontend display"
// @Router /team/new [post]
func (a Analyser) AddTeamMember(c *gin.Context) {
	var eventRequest globals.TeamMember
	if err := c.BindJSON(&eventRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := a.ESClient.AddTeamMember(eventRequest); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{})
}

// EditTeamMember godoc
// @Summary Modify a specific team member
// @Description Changes the name or the slack ID (or both)
// @Description for the document matching the given ID
// @Description Authentication and admin access are required for this endpoint
// @Tags Team
// @ID edit-team-member
// @Produce  json
// @Param team_member query string true "Team member id to update"
// @Param slack_id body object true "ID of the user in slack"
// @Param name body object true "The name of the user only used in frontend display"
// @Router /team/:team_member [put]
func (a Analyser) EditTeamMember(c *gin.Context) {
	documentID := c.Param("team_member")
	var eventRequest globals.TeamMember
	if err := c.BindJSON(&eventRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := a.ESClient.EditTeamMember(documentID, eventRequest); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{})
}

// DeleteTeamMember godoc
// @Summary Delete specified team member
// @Description Removes the document matching the team_member ID given in parameter
// @Description Authentication and admin access are required for this endpoint
// @Tags Team
// @ID delete-team-member
// @Produce  json
// @Param team_member query string true "Team member id to delete"
// @Router /team/:team_member [delete]
func (a Analyser) DeleteTeamMember(c *gin.Context) {
	teamMemberID := c.Param("team_member")
	if err := a.ESClient.DeleteTeamMember(teamMemberID); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(204, gin.H{})
}
