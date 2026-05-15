package static

import (
	"os"
	"time"

	"github.com/robgtest/blog/internal/models"
)

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

var uri = func() string {
	if os.Getenv("ENV") == "production" {
		return "https://blog.bob-productions.dev/"
	}
	return "http://localhost:7000/"
}()

var ControlAndChoiceData = models.BlogMeta{
	Title:       "Weekly Stoic: Control & Choice",
	Url:         uri + "blog/control-and-choice",
	Description: "Quick testing wellbeing tidbits",
	ImageUri:    uri + "images/stoic1.webp",
	Published:   date(2025, 1, 1),
}

var ToBeSteadyData = models.BlogMeta{
	Title:       "Weekly Stoic: To Be Steady & Unsteady",
	Url:         uri + "blog/to-be-steady",
	Description: "Quick testing wellbeing tidbits",
	ImageUri:    uri + "images/stoic2.webp",
	Published:   date(2025, 1, 10),
}

var AWSServerlessData = models.BlogMeta{
	Description: "Achieving performant cloud architecture: AWS Lambdas",
	Url:         uri + "blog/serverless",
	Title:       "Software Performance Guide: AWS Lambdas",
	ImageUri:    uri + "images/lambda-serverless/AWS-Meta.webp",
	Published:   date(2025, 1, 1),
}

var IntroData = models.BlogMeta{
	Description: "Who is the mysterious sweaty fox",
	Url:         uri + "blog/intro",
	Title:       "An Introduction To Bob Productions",
	ImageUri:    uri + "images/Sweat.webp",
	Published:   date(2025, 1, 1),
}

var IsCopilotADudData = models.BlogMeta{
	Description: "I suggest you stop relying on suggestions",
	Url:         uri + "blog/ai-autocomplete",
	Title:       "The Code Suggestion Crisis",
	ImageUri:    uri + "images/copilot/skullpilot.webp",
	Published:   date(2025, 1, 10),
}

var PerformanceWorkshop = models.BlogMeta{
	Description: "Performance Workshop AD Quarterly",
	Url:         uri + "blog/perf-workshop",
	Title:       "AD Performance Workshop",
	ImageUri:    uri + "images/performance/performance.webp",
	Published:   date(2025, 4, 10),
}

var GrugAutomationData = models.BlogMeta{
	Description: "Why test automation projects fail, explained with simple Grug wisdom and practical engineering lessons.",
	Url:         uri + "blog/grug-automation",
	Title:       "Grug Guide to Why Test Automation Fails",
	ImageUri:    uri + "images/testing/grug-automation.webp",
	Published:   date(2026, 5, 4),
}

var QuietSkillsData = models.BlogMeta{
	Description: "Getting production systems started and running.",
	Url:         uri + "blog/quiet-skills",
	Title:       "Skills for New Production Systems",
	ImageUri:    uri + "images/docs.webp",
	Published:   date(2026, 5, 9),
}
