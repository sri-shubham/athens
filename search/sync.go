package search

import (
	"context"

	elastic "github.com/olivere/elastic/v7"
	"github.com/sri-shubham/athens/models"
	"golang.org/x/sync/errgroup"
)

type SyncHelper struct {
	client          *elastic.Client
	users           models.Users
	projects        models.Projects
	projectHashtags models.ProjectHashtags
	hashtags        models.Hashtags
	userProjects    models.UserProjects
}

func NewSyncHelper(client *elastic.Client,
	users models.Users,
	projects models.Projects,
	projectHashtags models.ProjectHashtags,
	hashtags models.Hashtags,
	userProjects models.UserProjects) *SyncHelper {
	return &SyncHelper{
		client:          client,
		users:           users,
		projects:        projects,
		projectHashtags: projectHashtags,
		hashtags:        hashtags,
		userProjects:    userProjects,
	}
}

func (sh *SyncHelper) SyncAll(ctx context.Context) error {
	wg, _ := errgroup.WithContext(ctx)

	usersMap := make(map[int64]*models.User)
	projects := []*models.Project{}
	projectHashtagsMap := make(map[int64][]*models.ProjectHashtag)
	hashtagsMap := make(map[int64]*models.Hashtag)
	userProjectsMap := make(map[int64]*models.UserProject)

	wg.Go(func() (err error) {
		users, err := sh.users.GetAll()
		if err != nil {
			return err
		}

		for _, user := range users {
			usersMap[user.ID] = user
		}

		return nil
	})

	wg.Go(func() (err error) {
		projects, err = sh.projects.GetAll()
		return err
	})

	wg.Go(func() (err error) {
		projectHashtags, err := sh.projectHashtags.GetAll()
		if err != nil {
			return err
		}

		for _, ph := range projectHashtags {
			if _, exists := hashtagsMap[ph.ID]; !exists {
				projectHashtagsMap[ph.ID] = []*models.ProjectHashtag{}
			}
			projectHashtagsMap[ph.ID] = append(projectHashtagsMap[ph.ID], ph)
		}

		return nil
	})

	wg.Go(func() (err error) {
		hashtags, err := sh.hashtags.GetAll()
		if err != nil {
			return err
		}

		for _, h := range hashtags {
			hashtagsMap[h.ID] = h
		}

		return nil
	})

	wg.Go(func() (err error) {
		userProjects, err := sh.userProjects.GetAll()
		if err != nil {
			return err
		}

		for _, up := range userProjects {
			userProjectsMap[up.ProjectID] = up
		}

		return err
	})

	err := wg.Wait()
	if err != nil {
		return err
	}

	documents := make([]ProjectDocument, 0, len(projects))

	for _, project := range projects {
		projectHashtags := projectHashtagsMap[project.ID]
		userProject, ok := userProjectsMap[project.ID]
		if !ok {
			continue
		}

		user, ok := usersMap[userProject.UserID]
		if !ok {
			continue
		}

		hashtags := make([]Hashtag, 0, len(projectHashtags))
		for _, ph := range projectHashtags {
			hashtag := hashtagsMap[ph.HashtagID]
			hashtags = append(hashtags, Hashtag{
				ID:   hashtag.ID,
				Name: hashtag.Name,
			})
		}

		documents = append(documents, ProjectDocument{
			ProjectID:   int(project.ID),
			Name:        project.Name,
			Slug:        project.Slug,
			Description: project.Description,
			User: User{
				ID:   user.ID,
				Name: user.Name,
			},
			Hashtags: hashtags,
		})
	}

	return SyncData(ctx, sh.client, documents, true)
}
