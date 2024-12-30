package schema

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`

	Discussions      []Discussion     `json:"discussions" gorm:"foreignKey:UserID"`
	DiscusionReplies []DiscusionReply `json:"discusionReplies" gorm:"foreignKey:UserID"`
	Notes            []Note           `json:"notes" gorm:"foreignKey:UserID"`
	NoteReplies      []NoteReply      `json:"noteReplies" gorm:"foreignKey:UserID"`
	Content          []Content        `json:"posts" gorm:"foreignKey:UserID"`
}

type Discussion struct {
	gorm.Model
	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`

	Replies []DiscusionReply `json:"replies" gorm:"foreignKey:DiscussionID"`
}

type DiscusionReply struct {
	gorm.Model

	DiscussionID uint        `json:"discussionId"`
	Discussion   *Discussion `json:"discussion" gorm:"foreignKey:DiscussionID"`

	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`
}

type Note struct {
	gorm.Model
	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`

	Replies []NoteReply `json:"replies" gorm:"foreignKey:NoteID"`
}

type NoteReply struct {
	gorm.Model

	NoteID uint  `json:"noteId"`
	Note   *Note `json:"note" gorm:"foreignKey:NoteID"`

	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`
}

type Content struct {
	gorm.Model

	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID uint   `json:"userId"`
	Type   string `json:"type"` //forumContent, note, etc.

	// optional
	User        *User        `json:"user" gorm:"foreignKey:UserID"`
	Discussions []Discussion `json:"discussions" gorm:"foreignKey:ContentID"`
	Notes       []Note       `json:"notes" gorm:"foreignKey:ContentID"`
}

type ContentTag struct {
	gorm.Model
	ContentID    uint   `json:"contentId"`
	TagID        uint   `json:"tagId"`
	Relationship string `json:"relationship"` //potentially use enums for relationships
}

type Tag struct {
	gorm.Model
	Name string `json:"name"`

	Content []Content `json:"contents" gorm:"many2many:content_tags"`
}
