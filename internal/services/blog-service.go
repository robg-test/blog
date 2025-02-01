package services

import (
	"log"

	"github.com/robgtest/blog/internal"
)

func UpdateBlogView(value int) {
	if _, err := internal.DB.Exec("UPDATE blog SET blog_views = blog_views + 1 WHERE ID = ?", value); err != nil {
		log.Printf("Error updating blog views: %v", err)
	}
}

func GetBlogView(value int) (int, error) {
	var blogViews int
	err := internal.DB.QueryRow("SELECT blog_views FROM blog WHERE ID = ?", value).Scan(&blogViews)
	if err != nil {
		return 0, err
	}
	return blogViews, nil
}
