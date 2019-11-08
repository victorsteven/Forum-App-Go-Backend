package tests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/victorsteven/forum/api/models"
)

func TestFindAllPosts(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}
	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Error seeding user and post  table %v\n", err)
	}
	//Where postInstance is an instance of the post initialize in setup_test.go
	posts, err := postInstance.FindAllPosts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*posts), 2)
}

func TestSavePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error user and post refreshing table %v\n", err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	newPost := models.Post{
		ID:       1,
		Title:    "This is the title",
		Content:  "This is the content",
		AuthorID: user.ID,
	}
	savedPost, err := newPost.SavePost(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	assert.Equal(t, newPost.ID, savedPost.ID)
	assert.Equal(t, newPost.Title, savedPost.Title)
	assert.Equal(t, newPost.Content, savedPost.Content)
	assert.Equal(t, newPost.AuthorID, savedPost.AuthorID)
}

func TestFindPostByID(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundPost, err := post.FindPostByID(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundPost.ID, post.ID)
	assert.Equal(t, foundPost.Title, post.Title)
	assert.Equal(t, foundPost.Content, post.Content)
}

func TestUpdateAPost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	postUpdate := models.Post{
		ID:       1,
		Title:    "modiUpdate",
		Content:  "modiupdate@example.com",
		AuthorID: post.AuthorID,
	}
	updatedPost, err := postUpdate.UpdateAPost(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedPost.ID, postUpdate.ID)
	assert.Equal(t, updatedPost.Title, postUpdate.Title)
	assert.Equal(t, updatedPost.Content, postUpdate.Content)
	assert.Equal(t, updatedPost.AuthorID, postUpdate.AuthorID)
}

func TestDeleteAPost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := post.DeleteAPost(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, isDeleted, int64(1))
}

func TestDeleteUserPosts(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	user, _, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}

	numberDeleted, err := postInstance.DeleteUserPosts(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error deleting the post: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(1))
}
