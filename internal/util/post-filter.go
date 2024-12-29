package util

import "bandicute-server/internal/storage/repository/post"

func FilterRecentPost(latestPost *post.Model, posts []*post.Model) []*post.Model {
	filteredPosts := make([]*post.Model, 0)
	if (latestPost == nil || latestPost == &post.Model{}) {
		filteredPosts = append(filteredPosts, posts[0])
		return filteredPosts
	}

	for _, p := range posts {
		if p.PublishedAt.Before(latestPost.PublishedAt) || p.GUID == latestPost.GUID {
			continue
		}
		filteredPosts = append(filteredPosts, p)
	}
	return filteredPosts
}
