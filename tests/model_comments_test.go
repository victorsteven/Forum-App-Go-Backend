package tests

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorsteven/forum/api/models"
)

func TestCreateComment(t *testing.T) {
	err := refreshUserPostAndCommentTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and comment table %v\n", err)
	}
	user, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Cannot seed user and post %v\n", err)
	}
	newComment := models.Comment{
		ID:     1,
		Body:   "This is the comment body",
		UserID: user.ID,
		PostID: post.ID,
	}
	savedComment, err := newComment.SaveComment(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the comment: %v\n", err)
		return
	}
	assert.Equal(t, newComment.ID, savedComment.ID)
	assert.Equal(t, newComment.UserID, savedComment.UserID)
	assert.Equal(t, newComment.PostID, savedComment.PostID)
	assert.Equal(t, newComment.Body, "This is the comment body")

}

func TestCommentsForAPost(t *testing.T) {

	err := refreshUserPostAndCommentTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and comment table %v\n", err)
	}
	post, users, comments, err := seedUsersPostsAndComments()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment table %v\n", err)
	}
	//Where commentInstance is an instance of the post initialize in setup_test.go
	_, err = commentInstance.GetComments(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error getting the comments: %v\n", err)
		return
	}
	assert.Equal(t, len(comments), 2)
	assert.Equal(t, len(users), 2) //two users like the post
}

func TestDeleteAComment(t *testing.T) {

	err := refreshUserPostAndCommentTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and comment table %v\n", err)
	}
	_, _, comments, err := seedUsersPostsAndComments()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment table %v\n", err)
	}
	// Delete the first comment
	for _, v := range comments {
		if v.ID == 2 {
			continue
		}
		commentInstance.ID = v.ID //commentInstance is defined in setup_test.go
	}
	isDeleted, err := commentInstance.DeleteAComment(server.DB)
	if err != nil {
		t.Errorf("this is the error deleting the like: %v\n", err)
		return
	}
	assert.Equal(t, isDeleted, int64(1))
}

func TestDeleteCommentsForAPost(t *testing.T) {

	err := refreshUserPostAndCommentTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and comment table %v\n", err)
	}
	post, _, _, err := seedUsersPostsAndComments()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment table %v\n", err)
	}
	numberDeleted, err := commentInstance.DeletePostComments(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error deleting the like: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(2))
}

func TestDeleteCommentsForAUser(t *testing.T) {

	var userID uint32

	err := refreshUserPostAndCommentTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and comment table %v\n", err)
	}
	_, users, _, err := seedUsersPostsAndComments()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment table %v\n", err)
	}

	// get the first user. When you delete this user, also delete his comment
	for _, v := range users {
		if v.ID == 2 {
			continue
		}
		userID = v.ID
	}
	numberDeleted, err := commentInstance.DeleteUserComments(server.DB, userID)
	if err != nil {
		t.Errorf("this is the error deleting the comment: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(1))
}
