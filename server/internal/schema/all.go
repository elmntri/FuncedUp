package schema

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// DeletedAt int  `json:"deletedAt" gorm:"index"`
}

type User struct {
	BaseModel

	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Points       int    `json:"points"`

	// Discussions      []Discussion     `json:"discussions" gorm:"foreignKey:OwnerID"`      // all discussions the user owns
	// DiscusionReplies []DiscusionReply `json:"discusionReplies" gorm:"foreignKey:OwnerID"` // all replies the user has made
	// Notes            []Note           `json:"notes" gorm:"foreignKey:OwnerID"`            // all notes the user owns
	// NoteReplies      []NoteReply      `json:"noteReplies" gorm:"foreignKey:OwnerID"`      // all replies the user has made
	// Content []Content `json:"posts" gorm:"foreignKey:OwnerID"`
}

type Discussion struct {
	BaseModel

	OwnerID   uint `json:"ownerId"`
	ContentID uint `json:"contentId"`

	// Owner   *User            `json:"owner" gorm:"foreignKey:OwnerID"`
	// Content *Content         `json:"content" gorm:"foreignKey:ContentID"`
	// Replies []DiscusionReply `json:"replies" gorm:"foreignKey:DiscussionID"`
}

type DiscusionReply struct {
	BaseModel

	OwnerID      uint `json:"ownerId"`
	DiscussionID uint `json:"discussionId"`
	ContentID    uint `json:"contentId"`

	// Owner   *User `json:"owner" gorm:"foreignKey:OwnerID"`
	// Discussion   *Discussion `json:"discussion" gorm:"foreignKey:DiscussionID"`
	// Content   *Content `json:"content" gorm:"foreignKey:ContentID"`
}

type Note struct {
	BaseModel

	OwnerID   uint `json:"ownerId"`
	ContentID uint `json:"contentId"`

	// Owner   *User `json:"owner" gorm:"foreignKey:OwnerID"`
	// Content   *Content `json:"content" gorm:"foreignKey:ContentID"`
	// Replies []NoteReply `json:"replies" gorm:"foreignKey:NoteID"`
}

type NoteReply struct {
	BaseModel

	OwnerID   uint `json:"ownerId"`
	NoteID    uint `json:"noteId"`
	ContentID uint `json:"contentId"`

	// Owner   *User `json:"owner" gorm:"foreignKey:OwnerID"`
	// Note   *Note `json:"note" gorm:"foreignKey:NoteID"`
	// Content   *Content `json:"content" gorm:"foreignKey:ContentID"`
}

type Content struct {
	BaseModel

	Title   string `json:"title"`
	Body    string `json:"body"`
	OwnerID uint   `json:"ownerId"`

	// optional
	// Owner          *User            `json:"owner" gorm:"foreignKey:OwnerID"`
	// Discussions    []Discussion     `json:"discussions" gorm:"foreignKey:ContentID"`
	// DiscusionReply []DiscusionReply `json:"discusionReply" gorm:"foreignKey:ContentID"`
	// Notes          []Note           `json:"notes" gorm:"foreignKey:ContentID"`
	// NoteReplies    []NoteReply      `json:"noteReplies" gorm:"foreignKey:ContentID"`

	Tags []Tag `json:"tags" gorm:"many2many:content_tags"`
}

type ContentTag struct {
	BaseModel
	ContentID    uint   `json:"contentId"`
	TagID        uint   `json:"tagId"`
	Relationship string `json:"relationship"` // potentially use enums for relationships
}

type Tag struct {
	BaseModel
	Name string `json:"name"`

	Content []Content `json:"contents" gorm:"many2many:content_tags"`
}
