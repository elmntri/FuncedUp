package schema

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Points       int    `json:"points"`

	Discussions      []Discussion     `json:"discussions" gorm:"foreignKey:OwnerID"`      // all discussions the user owns
	DiscusionReplies []DiscusionReply `json:"discusionReplies" gorm:"foreignKey:OwnerID"` // all replies the user has made
	Notes            []Note           `json:"notes" gorm:"foreignKey:OwnerID"`            // all notes the user owns
	NoteReplies      []NoteReply      `json:"noteReplies" gorm:"foreignKey:OwnerID"`      // all replies the user has made

	/*
	 * Content field contains all the content the user owns.
	 * To Fetch all relationships on a user's content:
	 *   1. Fetch all content the user owns
	 *   2. Preload all discussions, notes, and their replies
	 *
	 * Example:
	 *	db.Preload("Content.Discussions.Replies").
	 *   	  Preload("Content.Notes.Replies").
	 *   	  Where("id = ?", userID).
	 *   	  Find(&user)
	 */
	Content []Content `json:"posts" gorm:"foreignKey:OwnerID"`
}

type Discussion struct {
	gorm.Model

	OwnerID uint  `json:"ownerId"`
	Owner   *User `json:"owner" gorm:"foreignKey:OwnerID"`

	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`

	Replies []DiscusionReply `json:"replies" gorm:"foreignKey:DiscussionID"`
}

type DiscusionReply struct {
	gorm.Model

	OwnerID uint  `json:"ownerId"`
	Owner   *User `json:"owner" gorm:"foreignKey:OwnerID"`

	DiscussionID uint        `json:"discussionId"`
	Discussion   *Discussion `json:"discussion" gorm:"foreignKey:DiscussionID"`

	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`
}

type Note struct {
	gorm.Model

	OwnerID uint  `json:"ownerId"`
	Owner   *User `json:"owner" gorm:"foreignKey:OwnerID"`

	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`

	Replies []NoteReply `json:"replies" gorm:"foreignKey:NoteID"`
}

type NoteReply struct {
	gorm.Model

	OwnerID uint  `json:"ownerId"`
	Owner   *User `json:"owner" gorm:"foreignKey:OwnerID"`

	NoteID uint  `json:"noteId"`
	Note   *Note `json:"note" gorm:"foreignKey:NoteID"`

	ContentID uint     `json:"contentId"`
	Content   *Content `json:"content" gorm:"foreignKey:ContentID"`
}

type Content struct {
	gorm.Model

	Title   string `json:"title"`
	Body    string `json:"body"`
	OwnerID uint   `json:"ownerId"`

	// optional
	Owner          *User            `json:"owner" gorm:"foreignKey:OwnerID"`
	Discussions    []Discussion     `json:"discussions" gorm:"foreignKey:ContentID"`
	DiscusionReply []DiscusionReply `json:"discusionReply" gorm:"foreignKey:ContentID"`
	Notes          []Note           `json:"notes" gorm:"foreignKey:ContentID"`
	NoteReplies    []NoteReply      `json:"noteReplies" gorm:"foreignKey:ContentID"`

	Tags []Tag `json:"tags" gorm:"many2many:content_tags"`
}

type ContentTag struct {
	gorm.Model
	ContentID    uint   `json:"contentId"`
	TagID        uint   `json:"tagId"`
	Relationship string `json:"relationship"` // potentially use enums for relationships
}

type Tag struct {
	gorm.Model
	Name string `json:"name"`

	Content []Content `json:"contents" gorm:"many2many:content_tags"`
}
