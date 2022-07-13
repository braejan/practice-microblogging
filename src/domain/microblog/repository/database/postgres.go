package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/braejan/practice-microblogging/src/domain/common/constants"
	commonEntities "github.com/braejan/practice-microblogging/src/domain/common/entities"
	"github.com/braejan/practice-microblogging/src/domain/common/interfaces"
	"github.com/braejan/practice-microblogging/src/domain/microblog/entities"
	userEntities "github.com/braejan/practice-microblogging/src/domain/user/entities"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type MicroBlogPostgres struct {
	*commonEntities.Postgres
}
type IMicroblogPostgres interface {
	CreatePost(owner *userEntities.User, text string) (err error)
	GetPostByID(ID string, internal bool) (post *entities.MicroBlog, err error)
	GetAllPosts() (posts []*entities.MicroBlog, err error)
	GetAllPostsByUserID(userID string) (posts []*entities.MicroBlog, err error)
	LikePost(ID string, user *userEntities.User) (err error)
	DislikePost(ID string, user *userEntities.User) (err error)
	interfaces.IPostgres
}

func NewMicroblogPostgres() (microblogPostgres IMicroblogPostgres, err error) {
	microblogPostgres = &MicroBlogPostgres{
		Postgres: commonEntities.NewPostgres(),
	}
	return
}

const (
	createPostDML = `INSERT INTO microblog (id, user_id, text) VALUES ($1, $2, $3)`
)

func (microblogPostgres *MicroBlogPostgres) CreatePost(owner *userEntities.User, text string) (err error) {

	if err != nil {
		return err
	}
	tx, err := microblogPostgres.BeginTx(context.Background())
	if err != nil {
		return
	}
	defer func() {
		err = microblogPostgres.DB.Close()
	}()
	postId := uuid.NewString()
	_, err = tx.ExecContext(context.Background(), createPostDML, postId, owner.ID, text)
	if err != nil {
		_ = microblogPostgres.RollbackTx(tx)
		return
	}

	err = microblogPostgres.CommitTx(tx)
	return
}

const (
	GetPostByIDDML            = `SELECT id, user_id, text, visit_count FROM "microblog"  WHERE id = $1`
	CountLikesDislikesByIDDML = `SELECT count(status) AS count FROM microblog_tracking WHERE microblog_id = $1 AND status = $2`
	IncrementViewDML          = `UPDATE microblog  SET visit_count = visit_count + 1 WHERE id = $1`
)

func (microblogPostgres *MicroBlogPostgres) GetPostByID(ID string, internal bool) (post *entities.MicroBlog, err error) {
	microblogPostgres.DB, err = microblogPostgres.Open()
	defer func() {
		_ = microblogPostgres.DB.Close()
	}()
	rows, err := microblogPostgres.DB.QueryContext(context.Background(), GetPostByIDDML, ID)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		post = &entities.MicroBlog{
			Detail: &entities.MicroBlogTracking{},
		}
		err = rows.Scan(&post.ID, &post.UserID, &post.Text, &post.Detail.VisitCount)
		if err != nil {
			post = nil
			return
		}

	}
	rows, err = microblogPostgres.DB.QueryContext(context.Background(), CountLikesDislikesByIDDML, ID, constants.Like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&post.Detail.LikesCount)
		if err != nil {
			post = nil
			return
		}
	}
	rows, err = microblogPostgres.DB.QueryContext(context.Background(), CountLikesDislikesByIDDML, ID, constants.Dislike)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&post.Detail.DislikesCount)
		if err != nil {
			post = nil
			return
		}
	}
	//update views
	if !internal {
		tx, err := microblogPostgres.BeginTx(context.Background())
		if err != nil {
			return nil, err
		}
		_, err = tx.Exec(IncrementViewDML, ID)
		if err != nil {
			microblogPostgres.RollbackTx(tx)
			return nil, err
		}
		err = microblogPostgres.CommitTx(tx)
		if err != nil {
			return nil, err
		}
		post.Detail.VisitCount++
	}
	return
}

const (
	GetAllPostDML = `SELECT id, user_id, text, visit_count FROM "microblog" order by creation_date DESC`
)

func (microblogPostgres *MicroBlogPostgres) GetAllPosts() (posts []*entities.MicroBlog, err error) {
	microblogPostgres.DB, err = microblogPostgres.Open()
	defer func() {
		_ = microblogPostgres.DB.Close()
	}()
	rows, err := microblogPostgres.DB.Query(GetAllPostDML)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		post := &entities.MicroBlog{
			Detail: &entities.MicroBlogTracking{},
		}
		err = rows.Scan(&post.ID, &post.UserID, &post.Text, &post.Detail.VisitCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return
}

const (
	GetAllPostByUserIDDML = `SELECT m.id, m.user_id, m.text, m.visit_count FROM "microblog" WHERE m.user_id = $1 order by m.creation_date DESC`
)

func (microblogPostgres *MicroBlogPostgres) GetAllPostsByUserID(userID string) (posts []*entities.MicroBlog, err error) {
	microblogPostgres.DB, err = microblogPostgres.Open()
	defer func() {
		_ = microblogPostgres.DB.Close()
	}()
	rows, err := microblogPostgres.DB.QueryContext(context.Background(), GetAllPostByUserIDDML, userID)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		post := &entities.MicroBlog{}
		err = rows.Scan(&post.ID, &post.UserID, &post.Text, &post.Detail.VisitCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return
}

const (
	LikeOrDislikeDML = `UPDATE microblog_tracking  SET status = $1 WHERE microblog_id = $2 AND user_id = $3`
)

func (microblogPostgres *MicroBlogPostgres) LikePost(ID string, user *userEntities.User) (err error) {
	err = microblogPostgres.validateTrackingRelation(ID, user.ID)
	if err != nil {
		return
	}
	tx, err := microblogPostgres.BeginTx(context.Background())
	if err != nil {
		return err
	}
	_, err = tx.Exec(LikeOrDislikeDML, int(constants.Like), ID, user.ID)
	if err != nil {
		microblogPostgres.RollbackTx(tx)
		return err
	}
	err = microblogPostgres.CommitTx(tx)
	if err != nil {
		return err
	}
	return
}
func (microblogPostgres *MicroBlogPostgres) DislikePost(ID string, user *userEntities.User) (err error) {
	err = microblogPostgres.validateTrackingRelation(ID, user.ID)
	if err != nil {
		return
	}
	tx, err := microblogPostgres.BeginTx(context.Background())
	if err != nil {
		return err
	}
	_, err = tx.Exec(LikeOrDislikeDML, int(constants.Dislike), ID, user.ID)
	if err != nil {
		microblogPostgres.RollbackTx(tx)
		return err
	}
	err = microblogPostgres.CommitTx(tx)
	if err != nil {
		return err
	}
	return
}

func (microblogPostgres *MicroBlogPostgres) Open() (DB *sql.DB, err error) {
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		commonEntities.Host,
		commonEntities.Port,
		commonEntities.DBUser,
		commonEntities.DBSecret,
		commonEntities.DBName,
	)
	DB, err = sql.Open("postgres", sqlInfo)
	return
}

func (microblogPostgres *MicroBlogPostgres) BeginTx(ctx context.Context) (tx *sql.Tx, err error) {
	microblogPostgres.DB, err = microblogPostgres.Open()
	if err != nil {
		return
	}
	tx, err = microblogPostgres.DB.BeginTx(ctx, nil)
	return
}
func (microblogPostgres *MicroBlogPostgres) CommitTx(tx *sql.Tx) (err error) {
	err = tx.Commit()
	return
}
func (microblogPostgres *MicroBlogPostgres) RollbackTx(tx *sql.Tx) (err error) {
	err = tx.Rollback()
	return
}

const (
	createTrackingRelationDML   = `INSERT INTO microblog_tracking(microblog_id, user_id, status) VALUES ($1, $2, $3)`
	validateTrackingRelationDML = `SELECT microblog_id FROM microblog_tracking WHERE microblog_id = $1 AND user_id = $2`
)

func (microblogPostgres *MicroBlogPostgres) validateTrackingRelation(ID string, userID string) (err error) {
	microblogPostgres.DB, err = microblogPostgres.Open()
	if err != nil {
		return
	}
	row, err := microblogPostgres.DB.QueryContext(context.Background(), validateTrackingRelationDML, ID, userID)
	if err != nil {
		return
	}
	var postID string
	if row.Next() {
		err = row.Scan(&postID)
		if err != nil {
			return
		}
	}
	defer row.Close()
	if postID != "" {
		return
	}
	tx, err := microblogPostgres.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return
	}
	_, err = tx.ExecContext(context.Background(), createTrackingRelationDML, ID, userID, int(constants.Default))
	if err != nil {
		_ = microblogPostgres.RollbackTx(tx)
		return
	}
	err = microblogPostgres.CommitTx(tx)
	return
}
