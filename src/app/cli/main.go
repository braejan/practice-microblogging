package main

import (
	"fmt"

	commonUsecases "github.com/braejan/practice-microblogging/src/domain/common/usecases"
	microblogEntities "github.com/braejan/practice-microblogging/src/domain/microblog/entities"
	microblogUsecases "github.com/braejan/practice-microblogging/src/domain/microblog/usecases"
	userUsecases "github.com/braejan/practice-microblogging/src/domain/user/usecases"
)

func main() {
	userUsecases, err := userUsecases.NewUserUsecases()
	if err != nil {
		panic(err)
	}
	fmt.Println("Wellcome to microblogging CLI")
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Type your userID: ")
	//userID, _ := reader.ReadString('\n')
	userID := "braejan"
	user, err := userUsecases.FindUserByID(userID)

	if err != nil {
		panic(err)
	}
	if user == nil {
		err = userUsecases.CreateUser(userID, fmt.Sprintf("Default name for %s", userID))
		if err != nil {
			panic(err)
		}
		err = userUsecases.CreateUser("bruch", fmt.Sprintf("Default name for %s", "bruch"))
		if err != nil {
			panic(err)
		}
	}
	user, err = userUsecases.FindUserByID(userID)
	if err != nil {
		panic(err)
	}
	fmt.Println("user: ", user)
	var text string
	for i := 0; i < 203; i++ {
		text = fmt.Sprintf("%s%d", text, i)
	}
	commonUsecases := commonUsecases.NewCommonUsecases()
	err = commonUsecases.ValidatePostLength(text)
	if err == nil {
		panic("should be an error here")
	}
	text = ""
	err = commonUsecases.ValidatePostLength(text)
	if err == nil {
		panic("should be an error here")
	}

	microblogUsecases, err := microblogUsecases.NewMicroblogUsecases(userUsecases, commonUsecases)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 50; i++ {
		text = fmt.Sprintf("Entry number %d: this is just a test.", i)
		err = microblogUsecases.CreatePost(userID, text)
		if err != nil {
			panic(err)
		}
	}

	posts, err := microblogUsecases.GetAllPosts()
	if err != nil {
		panic(err)
	}
	if len(posts) != 50 {
		panic("len should be 50")
	}
	var last *microblogEntities.MicroBlog
	for i := 0; i < 5; i++ {
		last, err = microblogUsecases.GetPostByID(posts[0].ID, "bruch", false)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("item %v\n", last.Detail)
	if last.Detail.VisitCount != 5 {
		panic(fmt.Errorf("should be 5 not %d", last.Detail.VisitCount))
	}
	fmt.Printf("last before: %v\n", last)
	fmt.Printf("detail before: %v\n", last.Detail)
	err = microblogUsecases.LikePost(posts[0].ID, "bruch")
	if err != nil {
		panic(err)
	}
	err = microblogUsecases.DislikePost(posts[0].ID, "braejan")
	if err != nil {
		panic(err)
	}
	last, err = microblogUsecases.GetPostByID(posts[0].ID, "bruch", false)
	if err != nil {
		panic(err)
	}
	fmt.Printf("last: %v\n", last)
	fmt.Printf("detail: %v\n", last.Detail)

}
