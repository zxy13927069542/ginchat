package service

import (
	"ginchat/models"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GroupService struct {
	gm *models.GroupBasicModel
	cm *models.ContactModel
}

func NewGroupService() *GroupService {
	return &GroupService{
		gm: models.NewGroupBasicModel(),
		cm: models.NewContactModel(),
	}
}

//	CreateGroup 创建群聊
func (s *GroupService) CreateGroup(c *gin.Context) {
	var req CreateGroupReq
	if err := c.Bind(&req); err != nil {
		c.IndentedJSON(http.StatusOK, JSONBResult{400, "参数错误", nil})
		return
	}

	var group models.GroupBasic
	group.Name = req.Name
	group.OwnerID = utils.Str2Uint(req.OwnerId)
	group.Icon = req.Icon
	group.Desc = req.Desc
	if err := s.gm.CreateGroup(group); err != nil {
		c.IndentedJSON(http.StatusOK, JSONBResult{400, err.Error(), nil})
		return
	}

	c.IndentedJSON(200, JSONBResult{200, "成功", nil})
}

//	AddGroup 加入群聊
func (s *GroupService) AddGroup(c *gin.Context) {
	var req JoinGroupReq
	if err := c.Bind(&req); err != nil {
		c.IndentedJSON(http.StatusOK, JSONBResult{400, "参数错误", nil})
		return
	}

	if err := s.cm.AddGroup(utils.Str2Uint(req.UserId), utils.Str2Uint(req.GroupId)); err != nil {
		c.IndentedJSON(http.StatusOK, JSONBResult{400, err.Error(), nil})
		return
	}
	c.IndentedJSON(200, JSONBResult{200, "成功", nil})
}
