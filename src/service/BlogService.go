package service

import (
	"context"
	redisConfig "github.com/redis/go-redis/v9"
	"hmdp/src/config/redis"
	"hmdp/src/dto"
	"hmdp/src/model"
	"hmdp/src/utils"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type BlogService struct {
}

var BlogManager *BlogService

func (*BlogService) SaveBlog(userId int64, blog *model.Blog) (res int64, err error) {
	blog.CreateTime = time.Now()
	blog.UpdateTime = time.Now()

	id, err := blog.SaveBlog()
	if err != nil {
		logrus.Error("[Blog Service] failed to insert data!")
		return
	}
	var f model.Follow
	follows, err := f.GetFollowsByFollowId(userId)
	if err != nil {
		return
	}

	if follows == nil || len(follows) == 0 {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, value := range follows {
		followUserId := value.UserId

		redisKey := utils.FEED_KEY + strconv.FormatInt(followUserId, 10)
		redis.GetRedisClient().ZAdd(ctx, redisKey, redisConfig.Z{
			Member: blog.Id,
			Score:  float64(time.Now().Unix()),
		})
	}

	res = id
	return
}

func (*BlogService) LikeBlog(id int64, userId int64) (err error) {
	// var blog model.Blog
	// blog.Id = id
	// err = blog.IncreseLike()
	// return
	userStr := strconv.FormatInt(userId, 10)
	redisKey := utils.BLOG_LIKE_KEY + strconv.FormatInt(id, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = redis.GetRedisClient().ZScore(ctx, redisKey, userStr).Result()

	flag := false

	if err != nil {
		if err == redisConfig.Nil {
			flag = true
		} else {
			return err
		}
	}

	var blog model.Blog
	blog.Id = id

	if flag {
		// add like
		blog.IncrLike()
		// add the user
		err = redis.GetRedisClient().ZAdd(ctx, redisKey,
			redisConfig.Z{
				Score:  float64(time.Now().Unix()),
				Member: userStr,
			}).Err()
	} else {
		// have the data
		blog.DecrLike()
		err = redis.GetRedisClient().ZRem(ctx, redisKey, userStr).Err()
	}
	return err
}

func (*BlogService) QueryMyBlog(id int64, current int) ([]model.Blog, error) {
	var blog model.Blog
	blog.UserId = id
	blogs, err := blog.QueryBlogs(current)
	return blogs, err
}

func (*BlogService) QueryHotBlogs(current int) ([]model.Blog, error) {
	var blogUtils model.Blog
	blogs, err := blogUtils.QueryHots(current)
	if err != nil {
		return nil, err
	}
	for i := range blogs {
		id := blogs[i].UserId
		user, err := UserManager.GetUserById(id)
		if err != nil {
			logrus.Error(err.Error())
			continue
		}
		blogs[i].Icon = user.Icon
		blogs[i].Name = user.NickName
	}

	return blogs, nil
}

func (*BlogService) GetBlogById(id int64) (model.Blog, error) {
	var blog model.Blog
	err := blog.GetBlogById(id)
	if err != nil {
		return model.Blog{}, err
	}

	userId := blog.UserId
	user, err := UserManager.GetUserById(userId)

	if err != nil {
		return model.Blog{}, err
	}

	blog.Name = user.NickName
	blog.Icon = user.Icon

	return blog, err
}

func (*BlogService) QueryUserLike(id int64) ([]dto.UserDTO, error) {
	// get the redis key
	redisKey := utils.BLOG_LIKE_KEY + strconv.FormatInt(id, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	idStrs, err := redis.GetRedisClient().ZRange(ctx, redisKey, 0, 4).Result()
	if err != nil {
		return []dto.UserDTO{}, err
	}

	if idStrs == nil || len(idStrs) == 0 {
		return []dto.UserDTO{}, err
	}

	var ids []int64
	for _, value := range idStrs {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return []dto.UserDTO{}, err
		}
		ids = append(ids, id)
	}

	var userUtils model.User
	users, err := userUtils.GetUsersByIds(ids)
	if err != nil {
		return []dto.UserDTO{}, err
	}

	userDTOS := make([]dto.UserDTO, len(users))
	for i := range users {
		userDTOS[i].Id = users[i].Id
		userDTOS[i].Icon = users[i].Icon
		userDTOS[i].NickName = users[i].NickName
	}
	return userDTOS, nil
}

func (*BlogService) QueryBlogOfFollow(lastId int64, offset int, userId int64) (dto.ScrollResult[model.Blog], error) {
	redisKey := utils.FEED_KEY + strconv.FormatInt(userId, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := redis.GetRedisClient().ZRevRangeByScoreWithScores(ctx, redisKey,
		&redisConfig.ZRangeBy{
			Min:    "0",
			Max:    strconv.FormatInt(lastId, 10),
			Offset: int64(offset),
			Count:  2,
		}).Result()

	if err != nil {
		return dto.ScrollResult[model.Blog]{}, err
	}

	if result == nil || len(result) == 0 {
		return dto.ScrollResult[model.Blog]{}, err
	}

	var ids []int64

	var minTime int64 = 0
	os := 0 // the number of equal number

	for _, value := range result {
		// id , err := strconv.ParseInt(value.Member.(string) , 10 , 64)
		id := value.Member.(int64)
		if err != nil {
			return dto.ScrollResult[model.Blog]{}, err
		}
		ids = append(ids, id)

		if (int64)(value.Score) == minTime {
			os++
		} else {
			minTime = (int64)(value.Score)
			os = 1
		}
	}

	var blogUtils model.Blog
	blogs, err := blogUtils.QueryBlogByIds(ids)
	if err != nil {
		return dto.ScrollResult[model.Blog]{}, nil
	}

	for i := range blogs {
		createBlogUser(&blogs[i])
		isBlogLiked(userId, &blogs[i])
	}

	var r dto.ScrollResult[model.Blog]
	r.Data = blogs
	r.MinTime = minTime
	r.Offset = os

	return r, nil
}

func createBlogUser(blog *model.Blog) {
	userId := blog.UserId
	var userUtils model.User
	user, err := userUtils.GetUserById(userId)
	if err != nil {
		return
	}
	blog.Name = user.NickName
	blog.Icon = user.Icon
}

func isBlogLiked(userId int64, blog *model.Blog) {
	redisKey := utils.BLOG_LIKE_KEY + strconv.FormatInt(userId, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := redis.GetRedisClient().ZScore(ctx, redisKey, strconv.FormatInt(userId, 10)).Err()
	blog.IsLike = (err == redisConfig.Nil)
}
