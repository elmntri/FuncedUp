package seeder

import (
	"fmt"

	"funcedup/internal/schema"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// -------------------------------------------------------------------------
// SeedIDs holds references to seeded objects so that subsequent seeding
// functions can reuse the same IDs rather than querying the DB multiple times.
// -------------------------------------------------------------------------
type SeedIDs struct {
	Users       map[string]uuid.UUID // Key = username,        Value = user ID
	Contents    map[string]uuid.UUID // Key = content title,   Value = content ID
	Discussions map[string]uuid.UUID // Key = "owner-title",   Value = discussion ID
	Notes       map[string]uuid.UUID // Key = "owner-title",   Value = note ID
	Tags        map[string]uuid.UUID // Key = tag name,        Value = tag ID
}

// -------------------------------------------------------------------------
// SeedAll runs all seed functions in a proper sequence, passing around a
// single SeedIDs struct to store/retrieve IDs as they’re created or loaded.
// -------------------------------------------------------------------------
func (d *Domain) SeedAll() error {
	seedIDs := &SeedIDs{
		Users:       make(map[string]uuid.UUID),
		Contents:    make(map[string]uuid.UUID),
		Discussions: make(map[string]uuid.UUID),
		Notes:       make(map[string]uuid.UUID),
		Tags:        make(map[string]uuid.UUID),
	}

	// Run each seeding function in sequence
	if err := d.seedUsers(seedIDs); err != nil {
		d.logger.Error("SeedAll: failed to seed users", zap.Error(err))
		return err
	}
	if err := d.seedTags(seedIDs); err != nil {
		d.logger.Error("SeedAll: failed to seed tags", zap.Error(err))
		return err
	}
	if err := d.seedContent(seedIDs); err != nil {
		d.logger.Error("SeedAll: failed to seed content", zap.Error(err))
		return err
	}

	if err := d.seedContentTags(seedIDs); err != nil {
		d.logger.Error("SeedAll: failed to attach tags to content", zap.Error(err))
		return err
	}

	if err := d.seedDiscussions(seedIDs); err != nil {
		d.logger.Error("SeedAll: failed to seed discussions", zap.Error(err))
		return err
	}
	if err := d.seedNotes(seedIDs); err != nil {
		d.logger.Error("SeedAll: failed to seed notes", zap.Error(err))
		return err
	}
	if err := d.seedDiscussionReplies(seedIDs); err != nil {
		d.logger.Error("SeedAll: failed to seed discussion replies", zap.Error(err))
		return err
	}
	if err := d.seedNoteReplies(seedIDs); err != nil {
		d.logger.Error("SeedAll: failed to seed note replies", zap.Error(err))
		return err
	}

	return nil
}

// -------------------------------------------------------------------------
// Seeding functions
// -------------------------------------------------------------------------

// seedUsers seeds user data and stores their IDs in seedIDs.
func (d *Domain) seedUsers(seedIDs *SeedIDs) error {
	db := d.params.DB.GetDB()

	users := []schema.User{
		{
			Username:     "michael",
			Email:        "michael.chen@elmntri.com",
			PasswordHash: "testtesttest",
			Points:       0,
		},
		{
			Username:     "alan",
			Email:        "vimalan.renganattan@elmntri.com",
			PasswordHash: "testtesttest",
			Points:       0,
		},
		{
			Username:     "jeff",
			Email:        "jeff.hsu@elmntri.com",
			PasswordHash: "testtesttest",
			Points:       0,
		},
	}

	for _, user := range users {
		err := db.
			Where("email = ?", user.Email).
			FirstOrCreate(&user).
			Error
		if err != nil {
			return fmt.Errorf("failed to seed user %s: %w", user.Email, err)
		}
		seedIDs.Users[user.Username] = user.ID
	}
	return nil
}

// seedTags seeds some tags to be reused by Content records.
func (d *Domain) seedTags(seedIDs *SeedIDs) error {
	db := d.params.DB.GetDB()

	tags := []schema.Tag{
		{Name: "charge"},
		{Name: "clock"},
		{Name: "fuel"},
		{Name: "methylation"},
		{Name: "oxidation"},
		{Name: "reduction"},
	}

	for _, tag := range tags {
		err := db.
			Where("name = ?", tag.Name).
			FirstOrCreate(&tag).
			Error
		if err != nil {
			return fmt.Errorf("failed to seed tag %s: %w", tag.Name, err)
		}
		seedIDs.Tags[tag.Name] = tag.ID
	}
	return nil
}

// seedContent seeds content data using the user IDs from seedIDs to assign ownership.
func (d *Domain) seedContent(seedIDs *SeedIDs) error {
	db := d.params.DB.GetDB()

	contents := []schema.Content{
		{
			OwnerID: seedIDs.Users["michael"],
			Title:   "Michael's Content 1",
			Body:    "lorem ipsum dolor sit amet",
		},
		{
			OwnerID: seedIDs.Users["michael"],
			Title:   "Michael's Content 2",
			Body:    "lorem ipsum dolor sit amet",
		},
		{
			OwnerID: seedIDs.Users["michael"],
			Title:   "Michael's Content 3",
			Body:    "lorem ipsum dolor sit amet",
		},
		{
			OwnerID: seedIDs.Users["alan"],
			Title:   "Alan's Content 1",
			Body:    "lorem ipsum dolor sit amet",
		},
		{
			OwnerID: seedIDs.Users["alan"],
			Title:   "Alan's Content 2",
			Body:    "lorem ipsum dolor sit amet",
		},
		{
			OwnerID: seedIDs.Users["alan"],
			Title:   "Alan's Content 3",
			Body:    "lorem ipsum dolor sit amet",
		},
		{
			OwnerID: seedIDs.Users["jeff"],
			Title:   "Jeff's Content 1",
			Body:    "lorem ipsum dolor sit amet",
		},
		{
			OwnerID: seedIDs.Users["jeff"],
			Title:   "Jeff's Content 2",
			Body:    "lorem ipsum dolor sit amet",
		},
		{
			OwnerID: seedIDs.Users["jeff"],
			Title:   "Jeff's Content 3",
			Body:    "lorem ipsum dolor sit amet",
		},

		// Example "reply" content. Feel free to remove or expand.
		{
			OwnerID: seedIDs.Users["michael"],
			Title:   "Michael's Reply to Alan's Discussion 1",
			Body:    "lorem ipsum reply",
		},
		{
			OwnerID: seedIDs.Users["michael"],
			Title:   "Michael's Reply to Alan's Note 2",
			Body:    "lorem ipsum reply",
		},
		{
			OwnerID: seedIDs.Users["michael"],
			Title:   "Michael's Reply to Jeff's Discussion 2",
			Body:    "lorem ipsum reply",
		},
		{
			OwnerID: seedIDs.Users["michael"],
			Title:   "Michael's Reply to Jeff's Note 1",
			Body:    "lorem ipsum reply",
		},
		{
			OwnerID: seedIDs.Users["alan"],
			Title:   "Alan's Reply to Michael's Discussion 1",
			Body:    "lorem ipsum reply",
		},
		{
			OwnerID: seedIDs.Users["alan"],
			Title:   "Alan's Reply to Michael's Note 2",
			Body:    "lorem ipsum reply",
		},
		{
			OwnerID: seedIDs.Users["alan"],
			Title:   "Alan's Reply to Jeff's Discussion 2",
			Body:    "lorem ipsum reply",
		},
		{
			OwnerID: seedIDs.Users["alan"],
			Title:   "Alan's Reply to Jeff's Note 1",
			Body:    "lorem ipsum reply",
		},
		{
			OwnerID: seedIDs.Users["jeff"],
			Title:   "Jeff's Reply to Michael's Discussion 1",
			Body:    "lorem ipsum reply",
		},
		{
			OwnerID: seedIDs.Users["jeff"],
			Title:   "Jeff's Reply to Michael's Note 2",
			Body:    "lorem ipsum reply",
		},
		{
			OwnerID: seedIDs.Users["jeff"],
			Title:   "Jeff's Reply to Alan's Discussion 2",
			Body:    "lorem ipsum reply",
		},
	}

	for _, content := range contents {
		err := db.
			Where("title = ?", content.Title).
			FirstOrCreate(&content).
			Error
		if err != nil {
			return fmt.Errorf("failed to seed content %s: %w", content.Title, err)
		}
		seedIDs.Contents[content.Title] = content.ID
	}
	return nil
}

// seedContentTags attaches tags to content. Uncomment the call in SeedAll()
// to use it.
func (d *Domain) seedContentTags(seedIDs *SeedIDs) error {
	db := d.params.DB.GetDB()

	contentTags := []schema.ContentTag{
		{
			ContentID: seedIDs.Contents["Michael's Content 1"],
			TagID:     seedIDs.Tags["charge"],
		},
		{
			ContentID: seedIDs.Contents["Michael's Content 1"],
			TagID:     seedIDs.Tags["clock"],
		},
		{
			ContentID: seedIDs.Contents["Michael's Content 1"],
			TagID:     seedIDs.Tags["fuel"],
		},
		{
			ContentID: seedIDs.Contents["Michael's Content 2"],
			TagID:     seedIDs.Tags["methylation"],
		},
		{
			ContentID: seedIDs.Contents["Michael's Content 2"],
			TagID:     seedIDs.Tags["oxidation"],
		},
		{
			ContentID: seedIDs.Contents["Michael's Content 2"],
			TagID:     seedIDs.Tags["reduction"],
		},
		{
			ContentID: seedIDs.Contents["Michael's Content 3"],
			TagID:     seedIDs.Tags["charge"],
		},
		{
			ContentID: seedIDs.Contents["Michael's Content 3"],
			TagID:     seedIDs.Tags["clock"],
		},
		{
			ContentID: seedIDs.Contents["Michael's Content 3"],
			TagID:     seedIDs.Tags["fuel"],
		},
		{
			ContentID: seedIDs.Contents["Alan's Content 1"],
			TagID:     seedIDs.Tags["methylation"],
		},
		{
			ContentID: seedIDs.Contents["Alan's Content 1"],
			TagID:     seedIDs.Tags["oxidation"],
		},
		{
			ContentID: seedIDs.Contents["Alan's Content 1"],
			TagID:     seedIDs.Tags["reduction"],
		},
		{
			ContentID: seedIDs.Contents["Alan's Content 2"],
			TagID:     seedIDs.Tags["charge"],
		},
		{
			ContentID: seedIDs.Contents["Alan's Content 2"],
			TagID:     seedIDs.Tags["clock"],
		},
		{
			ContentID: seedIDs.Contents["Alan's Content 2"],
			TagID:     seedIDs.Tags["fuel"],
		},
		{
			ContentID: seedIDs.Contents["Alan's Content 3"],
			TagID:     seedIDs.Tags["methylation"],
		},
		{
			ContentID: seedIDs.Contents["Alan's Content 3"],
			TagID:     seedIDs.Tags["oxidation"],
		},
		{
			ContentID: seedIDs.Contents["Alan's Content 3"],
			TagID:     seedIDs.Tags["reduction"],
		},
		{
			ContentID: seedIDs.Contents["Jeff's Content 1"],
			TagID:     seedIDs.Tags["charge"],
		},
		{
			ContentID: seedIDs.Contents["Jeff's Content 1"],
			TagID:     seedIDs.Tags["clock"],
		},
		{
			ContentID: seedIDs.Contents["Jeff's Content 1"],
			TagID:     seedIDs.Tags["fuel"],
		},
		{
			ContentID: seedIDs.Contents["Jeff's Content 2"],
			TagID:     seedIDs.Tags["methylation"],
		},
		{
			ContentID: seedIDs.Contents["Jeff's Content 2"],
			TagID:     seedIDs.Tags["oxidation"],
		},
		{
			ContentID: seedIDs.Contents["Jeff's Content 2"],
			TagID:     seedIDs.Tags["reduction"],
		},
		{
			ContentID: seedIDs.Contents["Jeff's Content 3"],
			TagID:     seedIDs.Tags["charge"],
		},
		{
			ContentID: seedIDs.Contents["Jeff's Content 3"],
			TagID:     seedIDs.Tags["clock"],
		},
		{
			ContentID: seedIDs.Contents["Jeff's Content 3"],
			TagID:     seedIDs.Tags["fuel"],
		},
	}

	for _, contentTag := range contentTags {
		err := db.
			Where("content_id = ? AND tag_id = ?", contentTag.ContentID, contentTag.TagID).
			FirstOrCreate(&contentTag).
			Error
		if err != nil {
			return fmt.Errorf("failed to seed content-tag relationship: %w", err)
		}
	}
	return nil
}

// seedDiscussions seeds discussion data, linking owners and content by IDs.
func (d *Domain) seedDiscussions(seedIDs *SeedIDs) error {
	db := d.params.DB.GetDB()

	discussions := []schema.Discussion{
		{
			OwnerID:   seedIDs.Users["michael"],
			ContentID: seedIDs.Contents["Michael's Content 1"],
		},
		{
			OwnerID:   seedIDs.Users["michael"],
			ContentID: seedIDs.Contents["Michael's Content 2"],
		},
		{
			OwnerID:   seedIDs.Users["michael"],
			ContentID: seedIDs.Contents["Michael's Content 3"],
		},
		{
			OwnerID:   seedIDs.Users["alan"],
			ContentID: seedIDs.Contents["Alan's Content 1"],
		},
		{
			OwnerID:   seedIDs.Users["alan"],
			ContentID: seedIDs.Contents["Alan's Content 2"],
		},
		{
			OwnerID:   seedIDs.Users["alan"],
			ContentID: seedIDs.Contents["Alan's Content 3"],
		},
		{
			OwnerID:   seedIDs.Users["jeff"],
			ContentID: seedIDs.Contents["Jeff's Content 1"],
		},
		{
			OwnerID:   seedIDs.Users["jeff"],
			ContentID: seedIDs.Contents["Jeff's Content 2"],
		},
		{
			OwnerID:   seedIDs.Users["jeff"],
			ContentID: seedIDs.Contents["Jeff's Content 3"],
		},
	}

	// We'll use these keys to store reference IDs in seedIDs.Discussions
	discussionKeys := []string{
		"michael-1", "michael-2", "michael-3",
		"alan-1", "alan-2", "alan-3",
		"jeff-1", "jeff-2", "jeff-3",
	}

	for i, discussion := range discussions {
		err := db.
			Where("owner_id = ? AND content_id = ?", discussion.OwnerID, discussion.ContentID).
			FirstOrCreate(&discussion).
			Error
		if err != nil {
			return fmt.Errorf("failed to seed discussion for contentID %v: %w", discussion.ContentID, err)
		}
		seedIDs.Discussions[discussionKeys[i]] = discussion.ID
	}
	return nil
}

// seedNotes seeds note data, linking owners and content by IDs.
func (d *Domain) seedNotes(seedIDs *SeedIDs) error {
	db := d.params.DB.GetDB()

	notes := []schema.Note{
		{
			OwnerID:   seedIDs.Users["michael"],
			ContentID: seedIDs.Contents["Michael's Content 1"],
		},
		{
			OwnerID:   seedIDs.Users["michael"],
			ContentID: seedIDs.Contents["Michael's Content 2"],
		},
		{
			OwnerID:   seedIDs.Users["michael"],
			ContentID: seedIDs.Contents["Michael's Content 3"],
		},
		{
			OwnerID:   seedIDs.Users["alan"],
			ContentID: seedIDs.Contents["Alan's Content 1"],
		},
		{
			OwnerID:   seedIDs.Users["alan"],
			ContentID: seedIDs.Contents["Alan's Content 2"],
		},
		{
			OwnerID:   seedIDs.Users["alan"],
			ContentID: seedIDs.Contents["Alan's Content 3"],
		},
		{
			OwnerID:   seedIDs.Users["jeff"],
			ContentID: seedIDs.Contents["Jeff's Content 1"],
		},
		{
			OwnerID:   seedIDs.Users["jeff"],
			ContentID: seedIDs.Contents["Jeff's Content 2"],
		},
		{
			OwnerID:   seedIDs.Users["jeff"],
			ContentID: seedIDs.Contents["Jeff's Content 3"],
		},
	}

	// We'll use these keys to store reference IDs in seedIDs.Notes
	noteKeys := []string{
		"michael-1", "michael-2", "michael-3",
		"alan-1", "alan-2", "alan-3",
		"jeff-1", "jeff-2", "jeff-3",
	}

	for i, note := range notes {
		err := db.
			Where("owner_id = ? AND content_id = ?", note.OwnerID, note.ContentID).
			FirstOrCreate(&note).
			Error
		if err != nil {
			return fmt.Errorf("failed to seed note for contentID %v: %w", note.ContentID, err)
		}
		seedIDs.Notes[noteKeys[i]] = note.ID
	}
	return nil
}

// seedDiscussionReplies seeds replies from different people to existing discussions.
func (d *Domain) seedDiscussionReplies(seedIDs *SeedIDs) error {
	db := d.params.DB.GetDB()

	// Example with just one reply. Add more as needed.
	replies := []schema.DiscusionReply{
		{
			OwnerID:      seedIDs.Users["michael"],
			DiscussionID: seedIDs.Discussions["michael-1"],
			ContentID:    seedIDs.Contents["Michael's Reply to Alan's Discussion 1"],
		},
	}

	for _, reply := range replies {
		err := db.
			Where("owner_id = ? AND discussion_id = ? AND content_id = ?",
				reply.OwnerID, reply.DiscussionID, reply.ContentID).
			FirstOrCreate(&reply).
			Error
		if err != nil {
			return fmt.Errorf("failed to seed discussion reply: %w", err)
		}
	}
	return nil
}

// seedNoteReplies seeds replies from different people to existing notes.
func (d *Domain) seedNoteReplies(seedIDs *SeedIDs) error {
	db := d.params.DB.GetDB()

	// Example with just one reply. Add more as needed.
	replies := []schema.NoteReply{
		{
			OwnerID:   seedIDs.Users["michael"],
			NoteID:    seedIDs.Notes["michael-1"],
			ContentID: seedIDs.Contents["Michael's Reply to Alan's Note 2"],
		},
	}

	for _, reply := range replies {
		err := db.
			Where("owner_id = ? AND note_id = ? AND content_id = ?",
				reply.OwnerID, reply.NoteID, reply.ContentID).
			FirstOrCreate(&reply).
			Error
		if err != nil {
			return fmt.Errorf("failed to seed note reply: %w", err)
		}
	}
	return nil
}
