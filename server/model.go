package server

import (
	"log"
	"time"
)

func (h *HttpServer) syncTable() {
	err := h.db.Sync2(new(User), new(Group), new(Message), new(GroupMessage))
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	Id       int64  `xorm:"id pk autoincr" json:"id"`
	Username string `xorm:"username" json:"username"`
	Password string `json:"password" xorm:"password"`
	Head     string `json:"head" xorm:"head"`
	CarNo    string `json:"car_no" xorm:"car_no unique"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}

type UserFriend struct {
	Id       int64 `xorm:"id pk autoincr" json:"id"`
	UserId   int64 `xorm:"user_id" json:"user_id"`
	FriendId int64 `xorm:"friend_id" json:"friend_id"`
	IsOk     bool  `xorm:"is_ok" json:"is_ok"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

type Group struct {
	Id           int64  `xorm:"id pk autoincr" json:"id"`
	Name         string `json:"name" xorm:"name"`
	CreatorId    int    `json:"creator_id" xorm:"creator_id"`
	IsBlockTalk  bool   `json:"is_block_talk" xorm:"is_block_talk"` // 禁言.
	MembersCount int64  `json:"members_count" xorm:"members_count"`
	GroupNumber  string `json:"group_number" xorm:"group_number unique"`
	Notification string `json:"notification" xorm:"notification"` // 公告
	MaxMember    int64  `json:"max_member" xorm:"max_member"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
}

func (Group) TableName() string {
	return "group"
}

type UserGroup struct {
	Id      int64 `json:"id" xorm:"id pk autoincr"`
	GroupId int   `json:"group_id" xorm:"group_id"`
	UserId  int   `json:"user_id" xorm:"user_id"`

	CreatedAt time.Time `xorm:"created" json:"created_at"`
}

func (UserGroup) TableName() string {
	return "user_group"
}

type Message struct {
	Id      int64 `xorm:"id pk autoincr" json:"id"`
	SendId  int64 `xorm:"send_id" json:"send_id"`
	ToId    int64 `xorm:"to_id" json:"to_id"`
	MsgType int   `json:"msg_type" xorm:"msg_type"`

	Content   string    `xorm:"content" json:"content"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
}

func (Message) TableName() string {
	return "message"
}

type GroupMessage struct {
	Id      int64 `xorm:"id pk autoincr" json:"id"`
	SendId  int64 `xorm:"send_id" json:"send_id"`
	ToId    int64 `xorm:"to_id" json:"to_id"`
	MsgType int   `json:"msg_type" xorm:"msg_type"`

	Content   string    `xorm:"content" json:"content"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
}

func (GroupMessage) TableName() string {
	return "group_message"
}
