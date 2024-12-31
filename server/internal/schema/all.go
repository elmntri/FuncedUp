// temporary file to hold all the schema structs, will be moved to individual files later

package schema

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Points       int    `json:"points"`

	Discussions      []Discussion     `json:"discussions" gorm:"foreignKey:UserID"`      // all discussions the user owns
	DiscusionReplies []DiscusionReply `json:"discusionReplies" gorm:"foreignKey:UserID"` // all replies the user has made
	Notes            []Note           `json:"notes" gorm:"foreignKey:UserID"`            // all notes the user owns
	NoteReplies      []NoteReply      `json:"noteReplies" gorm:"foreignKey:UserID"`      // all replies the user has made

	/*
		*Content field contains all the content the user owns.
		To Fetch all relationships on a user's content:
			1. Fetch all content the user owns
			2. Preload all discussions, notes, and their replies
		Example:
		```
		db.Preload("Content.Discussions.Replies"). Preload("Content.Notes.Replies"). Where("id = ?", userID). Find(&user)
		```
	*/
	Content []Content `json:"posts" gorm:"foreignKey:UserID"`
}

type Discussion struct {
	gorm.Model
	UserID uint  `json:"userId"`
	User   *User `json:"user" gorm:"foreignKey:UserID"`

	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`

	Replies []DiscusionReply `json:"replies" gorm:"foreignKey:DiscussionID"`
}

type DiscusionReply struct {
	gorm.Model
	UserID uint  `json:"userId"`
	User   *User `json:"user" gorm:"foreignKey:UserID"`

	DiscussionID uint        `json:"discussionId"`
	Discussion   *Discussion `json:"discussion" gorm:"foreignKey:DiscussionID"`

	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`
}

type Note struct {
	gorm.Model
	UserID uint  `json:"userId"`
	User   *User `json:"user" gorm:"foreignKey:UserID"`

	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`

	Replies []NoteReply `json:"replies" gorm:"foreignKey:NoteID"`
}

type NoteReply struct {
	gorm.Model
	UserID uint  `json:"userId"`
	User   *User `json:"user" gorm:"foreignKey:UserID"`

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

	// optional
	User           *User            `json:"user" gorm:"foreignKey:UserID"`
	Discussions    []Discussion     `json:"discussions" gorm:"foreignKey:ContentID"`
	DiscusionReply []DiscusionReply `json:"discusionReply" gorm:"foreignKey:ContentID"`
	Notes          []Note           `json:"notes" gorm:"foreignKey:ContentID"`
	NoteReplies    []NoteReply      `json:"noteReplies" gorm:"foreignKey:ContentID"`
	Tags           []Tag            `json:"tags" gorm:"many2many:content_tags"`
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
