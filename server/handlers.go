package server

import (
	"context"
	"errors"
	"github.com/clearcodecn/carim/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (h *HttpServer) index(ctx *gin.Context) {
	// TODO,
}

type LoginRequest struct {
	Username string `json:"username" form:"password"`
	Password string `json:"password" form:"password"`
}

func (h *HttpServer) login(ctx *gin.Context) {
	req := new(LoginRequest)
	if err := ctx.Bind(req); err != nil {
		ctx.AbortWithError(400, errors.New("params error"))
		return
	}
	if req.Username == "" || req.Password == "" {
		ctx.AbortWithError(400, errors.New("params error"))
		return
	}

	sess := h.db.NewSession()
	defer sess.Close()
	var user = new(User)
	ok, err := sess.Where("username = ?", req.Username).Get(user)
	if err != nil {
		sql, arg := sess.LastSQL()
		logrus.WithError(err).Warnf("query failed: %s, %v", sql, arg)
		ctx.AbortWithError(400, errors.New("server error"))
		return
	}
	if !ok {
		ctx.AbortWithError(400, errors.New("user not exist or password error"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.AbortWithError(400, errors.New("user not exist or password error"))
		return
	}

	j := jwt.New(jwt.SigningMethodHS256)
	j.Claims = jwt.MapClaims{
		"username": req.Username,
		"car_no":   user.CarNo,
		"id":       user.Id,
	}
	token, err := j.SignedString([]byte(h.config.Key))
	if err != nil {
		logrus.WithError(err).Warn("failed to generate token")
		ctx.AbortWithError(400, errors.New("server error"))
		return
	}
	ctx.JSON(200, gin.H{
		"status": true,
		"data":   token,
	})
}

type RegisterRequest struct {
	Username string `json:"username" form:"password"`
	Password string `json:"password" form:"password"`
}

func (h *HttpServer) register(ctx *gin.Context) {
	req := new(RegisterRequest)
	if err := ctx.Bind(req); err != nil {
		ctx.AbortWithError(400, errors.New("params error"))
		return
	}
	if req.Username == "" || req.Password == "" {
		ctx.AbortWithError(400, errors.New("params error"))
		return
	}

	sess := h.db.NewSession()
	defer sess.Close()
	var user = new(User)
	ok, err := sess.Where("username = ?", req.Username).Get(user)
	if err != nil {
		sql, arg := sess.LastSQL()
		logrus.WithError(err).Warnf("query failed: %s, %v", sql, arg)
		ctx.AbortWithError(400, errors.New("server error"))
		return
	}
	if ok {
		ctx.AbortWithError(400, errors.New("user already exist"))
		return
	}
	user.Username = req.Username
	pass, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user.Password = string(pass)
	user.CarNo = h.uniqueCarNo()
	_, err = sess.Insert(user)
	if err != nil {
		sql, arg := sess.LastSQL()
		logrus.WithError(err).Warnf("create user failed: %s", sql, arg)
		ctx.AbortWithError(500, errors.New("server error"))
		return
	}
	j := jwt.New(jwt.SigningMethodHS256)
	j.Claims = jwt.MapClaims{
		"username": req.Username,
		"car_no":   user.CarNo,
		"id":       user.Id,
	}
	token, err := j.SignedString([]byte(h.config.Key))
	if err != nil {
		logrus.WithError(err).Warn("failed to generate token")
		ctx.AbortWithError(500, errors.New("server error"))
		return
	}
	ctx.JSON(200, gin.H{
		"status": true,
		"data":   token,
	})
}

// TODO
func (h *HttpServer) forget(ctx *gin.Context) {

}

// TODO
func (h *HttpServer) changePassword(ctx *gin.Context) {

}

// -----------
// 好友api
// ----------

func (h *HttpServer) searchCarNo(ctx *gin.Context) {
	carNo := ctx.Query("car_no")
	if carNo == "" {
		ctx.AbortWithError(400, errors.New("params error"))
		return
	}
	sess := h.db.NewSession()
	defer sess.Close()
	var user = new(User)
	ok, err := sess.Where("car_no = ?", carNo).Select("id,head,username,car_no").Get(user)
	if err != nil {
		sql, arg := sess.LastSQL()
		logrus.WithError(err).Warnf("query failed: %s", sql, arg)
		ctx.AbortWithError(500, errors.New("server error"))
		return
	}
	if !ok {
		ctx.JSON(200, gin.H{
			"status": false,
		})
		return
	}
	ctx.JSON(200, gin.H{
		"data": user,
	})
}

type AddFriendRequest struct {
	FriendId int64 `json:"friend_id" form:"friend_id"`
}

func (h *HttpServer) addFriend(ctx *gin.Context) {
	req := new(AddFriendRequest)
	if err := ctx.Bind(req); err != nil {
		ctx.JSON(400, errors.New("params error"))
		return
	}
	user := userFromContext(ctx)
	if user.Id == req.FriendId {
		ctx.JSON(400, errors.New("不能添加自己"))
		return
	}
	sess := h.db.NewSession()
	defer sess.Close()
	fri := new(User)
	ok, err := sess.Where(`id = ?`, req.FriendId).Get(fri)
	if err != nil {
		sql, arg := sess.LastSQL()
		logrus.WithError(err).Warnf("query failed: %s", sql, arg)
		ctx.AbortWithError(500, errors.New("server error"))
		return
	}
	if !ok {
		ctx.JSON(400, gin.H{
			"status":  false,
			"message": "用户不存在",
		})
		return
	}
	uf := new(UserFriend)
	uf.UserId = user.Id
	uf.FriendId = req.FriendId
	_, err = sess.Insert(uf)
	if err != nil {
		sql, arg := sess.LastSQL()
		logrus.WithError(err).Warnf("create user friend failed: %s", sql, arg)
		ctx.AbortWithError(500, errors.New("server error"))
		return
	}
	go func() {
		h.pushClient.AddFriend(context.Background(), &proto.AddFriendRequest{
			SenderId: user.Id,
			ToId:     fri.Id,
			SenderNo: user.CarNo,
			ToNo:     fri.CarNo,
			Head:     user.Head,
			Nickname: user.Username,
		})
	}()
	ctx.JSON(200, gin.H{
		"status":  true,
		"message": "已发送好友请求，等待对方同意",
	})
}
