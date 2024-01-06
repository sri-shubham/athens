package search

import (
	"context"

	elastic "github.com/olivere/elastic/v7"
	"github.com/sri-shubham/athens/models"
	"github.com/sri-shubham/athens/util"
	"golang.org/x/sync/errgroup"
)

type SyncHelper struct {
	client          *elastic.Client
	users           models.Users
	projects        models.Projects
	projectHashtags models.ProjectHashtags
	hashtags        models.Hashtags
	userProjects    models.UserProjects

	usersQueue           util.UpdateQueue
	projectsQueue        util.UpdateQueue
	projectHashtagsQueue util.UpdateQueue
	hashtagsQueue        util.UpdateQueue
	userProjectsQueue    util.UpdateQueue
}

func NewSyncHelper(client *elastic.Client,
	users models.Users,
	projects models.Projects,
	projectHashtags models.ProjectHashtags,
	hashtags models.Hashtags,
	userProjects models.UserProjects,
	usersQueue util.UpdateQueue,
	projectsQueue util.UpdateQueue,
	projectHashtagsQueue util.UpdateQueue,
	hashtagsQueue util.UpdateQueue,
	userProjectsQueue util.UpdateQueue) *SyncHelper {
	return &SyncHelper{
		client:          client,
		users:           users,
		projects:        projects,
		projectHashtags: projectHashtags,
		hashtags:        hashtags,
		userProjects:    userProjects,

		usersQueue:           usersQueue,
		projectsQueue:        projectsQueue,
		projectHashtagsQueue: projectHashtagsQueue,
		hashtagsQueue:        hashtagsQueue,
		userProjectsQueue:    userProjectsQueue,
	}
}

func (sh *SyncHelper) SyncAll(ctx context.Context) error {
	wg, wgCtx := errgroup.WithContext(ctx)

	usersMap := make(map[int64]*models.User)
	projects := []*models.Project{}
	projectHashtagsMap := make(map[int64][]*models.ProjectHashtag)
	hashtagsMap := make(map[int64]*models.Hashtag)
	userProjectsMap := make(map[int64]*models.UserProject)

	wg.Go(func() (err error) {
		users, err := sh.users.GetAll(wgCtx)
		if err != nil {
			return err
		}

		for _, user := range users {
			usersMap[user.ID] = user
		}

		return nil
	})

	wg.Go(func() (err error) {
		projects, err = sh.projects.GetAll(wgCtx)
		return err
	})

	wg.Go(func() (err error) {
		projectHashtags, err := sh.projectHashtags.GetAll(wgCtx)
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
		hashtags, err := sh.hashtags.GetAll(wgCtx)
		if err != nil {
			return err
		}

		for _, h := range hashtags {
			hashtagsMap[h.ID] = h
		}

		return nil
	})

	wg.Go(func() (err error) {
		userProjects, err := sh.userProjects.GetAll(wgCtx)
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

	if len(projects) == 0 {
		return nil
	}

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

	if len(documents) == 0 {
		return nil
	}

	return SyncData(ctx, sh.client, documents, true)
}

func (sh *SyncHelper) BackgroundSync(ctx context.Context) error {
	wg, wgCtx := errgroup.WithContext(ctx)

	wg.Go(func() (err error) {
		for {
			item, ok := sh.usersQueue.Dequeue(wgCtx)
			if !ok {
				break
			}

			err = sh.SyncUser(ctx, item)
			if err != nil {
				return err
			}
		}
		return nil
	})

	wg.Go(func() (err error) {
		for {
			item, ok := sh.userProjectsQueue.Dequeue(wgCtx)
			if !ok {
				break
			}

			err = sh.SyncUserProject(ctx, item)
			if err != nil {
				return err
			}
		}
		return nil
	})

	wg.Go(func() (err error) {
		for {
			item, ok := sh.projectsQueue.Dequeue(wgCtx)
			if !ok {
				break
			}

			err = sh.SyncProject(ctx, item)
			if err != nil {
				return err
			}
		}
		return nil
	})

	wg.Go(func() (err error) {
		for {
			item, ok := sh.hashtagsQueue.Dequeue(wgCtx)
			if !ok {
				break
			}

			err = sh.SyncHashtag(ctx, item)
			if err != nil {
				return err
			}
		}
		return nil
	})

	wg.Go(func() (err error) {
		for {
			item, ok := sh.projectHashtagsQueue.Dequeue(wgCtx)
			if !ok {
				break
			}

			err = sh.SyncProjectHashtag(ctx, item)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return wg.Wait()
}

func (sh *SyncHelper) SyncUser(ctx context.Context, item *util.Item) error {
	// Specify the filter for documents to delete based on user ID
	filter := elastic.NewNestedQuery("users",
		elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("users.id", item.ID)),
	)

	switch item.Action {
	case util.ActionDelete:
		// Perform the delete by query
		_, err := sh.client.DeleteByQuery(indexName).
			Query(filter).
			Do(context.Background())
		if err != nil {
			return err
		}
	case util.ActionUpdate:
		value, _ := item.Value.(*models.User)
		oldValue, _ := item.OldValue.(*models.User)
		if value.Name == oldValue.Name {
			return nil
		}

		// Define the update script to update the 'name' field
		script := elastic.NewScript(`ctx._source.users.name = params.newName`).
			Param("newName", value.Name)

		// Perform the update by query
		_, err := sh.client.UpdateByQuery(indexName).
			Query(filter).
			Script(script).
			Do(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func (sh *SyncHelper) SyncUserProject(ctx context.Context, item *util.Item) error {
	value, _ := item.Value.(*models.UserProject)
	return sh.SyncProjects(ctx, []int64{value.ProjectID})
}

func (sh *SyncHelper) SyncProject(ctx context.Context, item *util.Item) error {
	value, _ := item.Value.(*models.Project)
	switch item.Action {
	case util.ActionDelete:
		// Perform the delete by query
		// Construct a TermQuery for the "project_id" field
		termQuery := elastic.NewTermQuery("project_id", value.ID)

		// Perform the delete by query
		_, err := sh.client.DeleteByQuery(indexName).
			Query(termQuery).
			Do(context.Background())
		if err != nil {
			return err
		}
	case util.ActionCreate, util.ActionUpdate:
		return sh.SyncProjects(ctx, []int64{item.ID})
	}

	return nil
}

func (sh *SyncHelper) SyncHashtag(ctx context.Context, item *util.Item) error {
	hashtag := item.Value.(*models.Hashtag)

	projectIds := []int64{}
	projectHashtags, err := sh.projectHashtags.GetByHashTag(ctx, hashtag.ID)
	if err != nil {
		return err
	}

	for _, projectHashtag := range projectHashtags {
		projectIds = append(projectIds, projectHashtag.ProjectID)
	}

	// Sync Project
	return sh.SyncProjects(ctx, projectIds)
}

func (sh *SyncHelper) SyncProjectHashtag(ctx context.Context, item *util.Item) error {
	value := item.Value.(*models.ProjectHashtag)

	// Sync projects
	return sh.SyncProjects(ctx, []int64{value.ProjectID})
}

func (sh *SyncHelper) SyncProjects(ctx context.Context, projectIDs []int64) error {

	projects, err := sh.projects.GetBulk(ctx, projectIDs)
	if err != nil {
		return err
	}

	hashtagsMap := make(map[int64]*models.Hashtag)
	documents := make([]ProjectDocument, 0, len(projects))

	// userProjectsMap := make(map[int64]*models.UserProject)

	for _, project := range projects {
		// Construct a TermQuery for the "project_id" field
		termQuery := elastic.NewTermQuery("project_id", project.ID)

		// Perform the delete by query
		_, err := sh.client.DeleteByQuery(indexName).
			Query(termQuery).
			Do(context.Background())
		if err != nil {
			return err
		}

		userProjects, err := sh.userProjects.GetByProjectID(ctx, project.ID)
		if err != nil {
			return err
		}

		if len(userProjects) == 0 {
			continue
		}

		projectHashtags, err := sh.projectHashtags.GetByProjectID(ctx, project.ID)
		if err != nil {
			return err
		}

		hashtagIds := []int64{}
		for _, ph := range projectHashtags {
			hashtagIds = append(hashtagIds, ph.HashtagID)
		}

		hashtags, err := sh.hashtags.GetBulk(ctx, hashtagIds)
		if err != nil {
			return err
		}

		for _, h := range hashtags {
			hashtagsMap[h.ID] = h
		}

		for _, up := range userProjects {
			user, err := sh.users.Get(ctx, up.UserID)
			if err != nil {
				return err
			}

			hashtags := make([]Hashtag, 0, len(projectHashtags))
			for _, ph := range projectHashtags {
				hashtag := hashtagsMap[ph.HashtagID]
				if hashtag == nil {
					continue
				}

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
	}

	if len(documents) == 0 {
		return nil
	}

	return SyncData(ctx, sh.client, documents, false)
}
