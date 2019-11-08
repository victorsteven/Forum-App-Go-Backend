package tests

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victorsteven/forum/api/models"
)

func TestSaveALike(t *testing.T) {
	err := refreshUserPostAndLikeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and like table %v\n", err)
	}
	user, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Cannot seed user and post %v\n", err)
	}
	newLike := models.Like{
		ID:     1,
		UserID: user.ID,
		PostID: post.ID,
	}
	savedLike, err := newLike.SaveLike(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the like: %v\n", err)
		return
	}
	assert.Equal(t, newLike.ID, savedLike.ID)
	assert.Equal(t, newLike.UserID, savedLike.UserID)
	assert.Equal(t, newLike.PostID, savedLike.PostID)
}

func TestGetLikeInfoForAPost(t *testing.T) {

	err := refreshUserPostAndLikeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and like table %v\n", err)
	}
	post, users, likes, err := seedUsersPostsAndLikes()
	if err != nil {
		log.Fatalf("Error seeding user, post and like table %v\n", err)
	}
	//Where likeInstance is an instance of the post initialize in setup_test.go
	_, err = likeInstance.GetLikesInfo(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error getting the likes: %v\n", err)
		return
	}
	assert.Equal(t, len(likes), 2)
	assert.Equal(t, len(users), 2) //two users like the post
}

func TestDeleteALike(t *testing.T) {

	err := refreshUserPostAndLikeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and like table %v\n", err)
	}
	_, _, likes, err := seedUsersPostsAndLikes()
	if err != nil {
		log.Fatalf("Error seeding user, post and like table %v\n", err)
	}
	// Delete the first like
	for _, v := range likes {
		if v.ID == 2 {
			continue
		}
		likeInstance.ID = v.ID //likeInstance is defined in setup_test.go
	}
	deletedLike, err := likeInstance.DeleteLike(server.DB)
	if err != nil {
		t.Errorf("this is the error deleting the like: %v\n", err)
		return
	}
	assert.Equal(t, deletedLike.ID, likeInstance.ID)
}

// When a post is deleted, delete its likes
func TestDeleteLikesForAPost(t *testing.T) {

	err := refreshUserPostAndLikeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and like table %v\n", err)
	}
	post, _, _, err := seedUsersPostsAndLikes()
	if err != nil {
		log.Fatalf("Error seeding user, post and like table %v\n", err)
	}
	numberDeleted, err := likeInstance.DeletePostLikes(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error deleting the like: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(2))
}

// When a user is deleted, delete its likes
func TestDeleteLikesForAUser(t *testing.T) {
	var userID uint32
	err := refreshUserPostAndLikeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and like table %v\n", err)
	}
	_, users, likes, err := seedUsersPostsAndLikes()
	if err != nil {
		log.Fatalf("Error seeding user, post and like table %v\n", err)
	}
	for _, v := range likes {
		if v.ID == 2 {
			continue
		}
		likeInstance.ID = v.ID //likeInstance is defined in setup_test.go
	}
	// get the first user, this user has one like
	for _, v := range users {
		if v.ID == 2 {
			continue
		}
		userID = v.ID
	}
	numberDeleted, err := likeInstance.DeleteUserLikes(server.DB, userID)
	if err != nil {
		t.Errorf("this is the error deleting the like: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(1))
}
