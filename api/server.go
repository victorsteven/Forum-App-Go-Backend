package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/victorsteven/fullstack/api/controllers"
)

var server = controllers.Server{}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), 	os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	//seed.Load(server.DB)

	//from := mail.NewEmail("Example User", "chodsteven@gmail.com")
	//subject := "Sending with SendGrid is Fun"
	//to := mail.NewEmail("Example User", "chikodi543@gmail.com")
	//plainTextContent := "and easy to do anywhere, even with Go"
	//htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	//message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	//client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	//_, err = client.Send(message)
	//if err != nil {
	//	log.Println(err)
	//} else {
	//	fmt.Println("Email sent")
	//	//fmt.Println(response.StatusCode)
	//	//fmt.Println(response.Body)
	//	//fmt.Println(response.Headers)
	//}

	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s", apiPort)

	server.Run(apiPort)

}
