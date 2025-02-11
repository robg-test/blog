package static

import (
	"os"

	"github.com/robgtest/blog/internal/models"
)

var uri = func() string {
	if os.Getenv("ENV") == "production" {
		return "https://blog.bob-productions.dev/"
	}
	return "http://localhost:8080/"
}()

var ControlAndChoiceData = models.BlogMeta{
	Title:       "Weekly Stoic: Control & Choice",
	Url:         uri + "blog/control-and-choice",
	Description: "Quick testing wellbeing tidbits",
	ImageUri:    uri + "images/stoic1.jpeg",
}

var ToBeSteadyData = models.BlogMeta{
	Title:       "Weekly Stoic: To Be Steady & Unsteady",
	Url:         uri + "blog/to-be-steady",
	Description: "Quick testing wellbeing tidbits",
	ImageUri:    uri + "images/stoic2.jpg",
}

var AWSServerlessData = models.BlogMeta{
	Description: "A Guide for performant cloud architecture: AWS Lambdas",
	Url:         uri + "blog/serverless",
	Title:       "Software Performance Guide: AWS Lambdas",
	ImageUri:    uri + "images/AWS-Meta.png",
}

var IntroData = models.BlogMeta{
	Description: "Who is the mysterious sweaty fox",
	Url:         uri + "blog/intro",
	Title:       "An Introduction To Bob Productions",
	ImageUri:    uri + "images/Sweat.webp",
}

var IsCopilotADudData = models.BlogMeta{
	Description: "An opinion piece on what one could consider a disastrous implementation.",
	Url:         uri + "blog/copilot-a-dud",
	Title:       "The AI Auto-complete Headache",
	ImageUri:    uri + "images/copilot.webp",
}
