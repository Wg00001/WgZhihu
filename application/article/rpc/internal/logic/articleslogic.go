package logic

import (
	"WgZhihu/application/article/rpc/internal/code"
	"WgZhihu/application/article/rpc/internal/model"
	"WgZhihu/application/article/rpc/internal/types"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
	"golang.org/x/exp/slices"
	"strconv"
	"time"

	"WgZhihu/application/article/rpc/internal/svc"
	"WgZhihu/application/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	prefixArticles = "biz#articles#%d#%d"
	//默认过期时间
	articlesExpire = 3600 * 24 * 2
)

type ArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticlesLogic {
	return &ArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ArticlesLogic) Articles(in *pb.ArticlesRequest) (*pb.ArticlesResponse, error) {
	if in.SortType != types.SortPublishTime && in.SortType != types.SortLikeCount {
		return nil, code.SortTypeInvalid
	}
	if in.UserId <= 0 {
		return nil, code.UserIdInvalid
	}
	if in.PageSize == 0 {
		in.PageSize = types.DefaultPageSize
	}
	if in.Cursor == 0 {
		if in.SortType == types.SortPublishTime {
			in.Cursor = time.Now().Unix()
		} else {
			in.Cursor = types.DefaultSortLikeCursor
		}
	}

	var (
		sortField       string
		sortLikeNum     int64
		sortPublishTime string
	)
	if in.SortType == types.SortLikeCount {
		sortField = "like_num"
		sortLikeNum = in.Cursor
	} else {
		sortField = "publish_time"
		sortPublishTime = time.Unix(in.Cursor, 0).Format("2006-01-02 15:04:05")
	}
	var (
		err error
		//使用标识符来进行状态转换
		isCache, isEnd bool
		lastId, cursor int64
		curPage        []*pb.ArticleItem
		articles       []*model.Article
	)
	//不用处理缓存错误，因为这是非主逻辑，缓存没有应该去sql里找
	articleIds, _ := l.cacheArticles(l.ctx, in.UserId, in.Cursor, in.PageSize, in.SortType)
	//长度大于0说明缓存里有
	if len(articleIds) > 0 {
		isCache = true
		if articleIds[len(articleIds)-1] == -1 {
			isEnd = true
		}
		//获取文章详情
		articles, err = l.articleByIds(l.ctx, articleIds)
		if err != nil {
			return nil, err
		}
		//根据选择的排序方式的不同，构造对应的排序函数
		cmpFunc := getSortFunc(sortField)
		slices.SortFunc(articles, cmpFunc)

		for _, article := range articles {
			curPage = append(curPage, &pb.ArticleItem{
				Id:           article.Id,
				Title:        article.Title,
				Content:      article.Content,
				LikeCount:    article.LikeNum,
				CommentCount: article.CommentNum,
				PublishTime:  article.PublishTime.Unix(),
			})
		}
	} else {
		v, err, _ := l.svcCtx.SingleFlightGroup.Do(fmt.Sprintf("ArticlesByUserId:%d:%d", in.UserId, in.SortType), func() (interface{}, error) {
			return l.svcCtx.ArticleModel.ArticleByUserId(l.ctx, in.UserId, types.ArticleStatusVisible, sortLikeNum, sortPublishTime, sortField, types.DefaultLimit)
		})
		if err != nil {
			logx.Errorf("ArticlesByUserId userId: %d sortField: %s error: %v", in.UserId, sortField, err)
			return nil, err
		}
		if v == nil {
			return &pb.ArticlesResponse{}, nil
		}
		articles = v.([]*model.Article)

		var firstPageArticles []*model.Article
		//查询时会多查数据，返回只返回第一页的数据
		//firstPageArticles = min(pageSize , articles)
		if len(articles) > int(in.PageSize) {
			firstPageArticles = articles[:int(in.PageSize)]
		} else {
			firstPageArticles = articles
			isEnd = true
		}

		for _, article := range firstPageArticles {
			curPage = append(curPage, &pb.ArticleItem{
				Id:           article.Id,
				Title:        article.Title,
				Content:      article.Content,
				LikeCount:    article.LikeNum,
				CommentCount: article.CommentNum,
				PublishTime:  article.PublishTime.Unix(),
			})
		}
	}

	//去重
	if len(curPage) > 0 {
		pageLast := curPage[len(curPage)-1]
		lastId = pageLast.Id
		//保存游标位置
		if in.SortType == types.SortPublishTime {
			cursor = pageLast.PublishTime
		} else {
			cursor = pageLast.LikeCount
		}
		if cursor < 0 {
			cursor = 0
		}
		for k, article := range curPage {
			//传入的文章id或时间相同的话，把后面的去掉
			if in.SortType == types.SortPublishTime {
				if article.PublishTime == in.Cursor && article.Id == in.ArticleId {
					curPage = curPage[k:]
					break
				}
			} else {
				if article.LikeCount == in.Cursor && article.Id == in.ArticleId {
					curPage = curPage[k:]
					break
				}
			}
		}
	}

	ret := &pb.ArticlesResponse{
		Articles:  curPage,
		IsEnd:     isEnd,
		Cursor:    cursor,
		ArticleId: lastId,
	}
	if !isCache {
		//不是从缓存拿的话，开个线程缓存一下查到的内容
		//缓存出错了也是不需要处理的
		threading.GoSafe(func() {
			//查询数量小于限制，说明已经查完了，插入一条id=-1的，表示最后一条数据
			if len(articles) < types.DefaultLimit && len(articles) > 0 {
				articles = append(articles, &model.Article{Id: -1})
			}
			err = l.addCacheArticles(context.Background(), articles, in.UserId, in.SortType)
			if err != nil {
				logx.Errorf("addCacheArticles error: %v", err)
			}
		})
	}
	return ret, nil
}

func (l *ArticlesLogic) cacheArticles(ctx context.Context, uid, cursor, ps int64, sortType int32) ([]int64, error) {
	key := articlesKey(uid, sortType)
	b, err := l.svcCtx.BizRedis.ExistsCtx(ctx, key)
	if err != nil {
		logx.Errorf("ExistsCtx key: %s error: %v", key, err)
	}
	if b {
		err = l.svcCtx.BizRedis.ExpireCtx(ctx, key, articlesExpire)
		if err != nil {
			logx.Errorf("ExpireCtx key: %s error: %v", key, err)
		}
	}
	pairs, err := l.svcCtx.BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx(ctx, key, 0, cursor, 0, int(ps))
	if err != nil {
		logx.Errorf("ZrevrangebyscoreWithScoresAndLimit key: %s error: %v", key, err)
		return nil, err
	}
	var ids []int64
	for _, pair := range pairs {
		id, err := strconv.ParseInt(pair.Key, 10, 64)
		if err != nil {
			logx.Errorf("strconv.ParseInt key: %s error: %v", pair.Key, err)
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// 使用MapReduce并行获取文章详情
func (l *ArticlesLogic) articleByIds(ctx context.Context, articleIds []int64) ([]*model.Article, error) {
	//int64, *model.Article, []*model.Article - 分别是generate生成的数据以及mapper的入参，mapper的出参以及reducer的入参，reducer的出参
	articles, err := mr.MapReduce[int64, *model.Article, []*model.Article](
		//generate 生成数据
		func(source chan<- int64) {
			for _, aid := range articleIds {
				if aid == -1 {
					continue
				}
				source <- aid
			}
		},
		//mapper 处理数据
		func(id int64, writer mr.Writer[*model.Article], cancel func(err error)) {
			p, err := l.svcCtx.ArticleModel.FindOne(ctx, id)
			if err != nil {
				cancel(err)
				return
			}
			writer.Write(p)
		},
		//reducer 聚合数据
		func(pipe <-chan *model.Article, writer mr.Writer[[]*model.Article], cancel func(error)) {
			var articles []*model.Article
			for article := range pipe {
				articles = append(articles, article)
			}
			writer.Write(articles)
		})
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func articlesKey(uid int64, sortType int32) string {
	return fmt.Sprintf(prefixArticles, uid, sortType)
}

func getSortFunc(sortField string) func(a, b *model.Article) int {
	if sortField == "like_num" {
		return func(a, b *model.Article) int {
			if b.LikeNum > a.LikeNum {
				return 1
			} else if b.LikeNum < a.LikeNum {
				return -1
			}
			return 0
		}
	} else {
		return func(a, b *model.Article) int {
			if b.PublishTime.Unix() > a.PublishTime.Unix() {
				return 1
			} else if b.PublishTime.Unix() < a.PublishTime.Unix() {
				return -1
			}
			return 0
		}
	}
}

func (l *ArticlesLogic) addCacheArticles(ctx context.Context, articles []*model.Article, userId int64, sortType int32) error {
	if len(articles) == 0 {
		return nil
	}
	key := articlesKey(userId, sortType)
	for _, article := range articles {
		var score int64
		if sortType == types.SortLikeCount {
			score = article.LikeNum
		} else if sortType == types.SortPublishTime && article.Id != -1 {
			score = article.PublishTime.Local().Unix()
		}
		if score < 0 {
			score = 0
		}
		_, err := l.svcCtx.BizRedis.ZaddCtx(ctx, key, score, strconv.Itoa(int(article.Id)))
		if err != nil {
			return err
		}
	}
	return l.svcCtx.BizRedis.ExpireCtx(ctx, key, articlesExpire)
}
