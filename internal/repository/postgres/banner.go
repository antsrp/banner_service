package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/antsrp/banner_service/internal/domain/models"
	"github.com/antsrp/banner_service/internal/repository"
	mapper "github.com/antsrp/banner_service/pkg/presenters"
	"github.com/jackc/pgx/v5"
)

type BannerStorage struct {
	conn *Connection
}

func NewBannerStorage(conn *Connection) BannerStorage {
	return BannerStorage{
		conn: conn,
	}
}

func (s BannerStorage) Create(ctx context.Context, banner models.Banner) (models.Banner, repository.DatabaseError) {
	var result models.Banner
	//result.ID = id
	contentData, err := mapper.ToJSON(banner.Content, &mapper.DefaultIndent)
	if err != nil {
		return models.Banner{}, NewError("can't present banner's content to json", err)
	}

	tx, err := s.conn.PC.Begin(ctx)
	if err != nil {
		return models.Banner{}, NewError("can't create transaction", err)
	}
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, `INSERT INTO banners (feature_id, is_active, content) VALUES ($1, $2, $3) RETURNING id;`, banner.FeatureID, banner.IsActive, contentData).
		Scan(&result.ID); err != nil {
		return models.Banner{}, NewError("can't create banner", err)
	}
	values := make([]string, 0, len(banner.TagIDS))
	for _, tag := range banner.TagIDS {
		values = append(values, fmt.Sprintf("(%d, %d)", result.ID, tag))
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`INSERT INTO banners_tags (banner_id, tag_id) VALUES %s`, strings.Join(values, ","))); err != nil {
		return models.Banner{}, NewError("can't add tags for banner", err)
	}

	tx.Commit(ctx)

	return result, nil
}

func (s BannerStorage) Update(ctx context.Context, banner models.Banner) repository.DatabaseError {
	query := `UPDATE banners SET %s WHERE id = $1`
	errString := "can't update banner"

	tx, err := s.conn.PC.Begin(ctx)
	if err != nil {
		return NewError(errString, err)
	}
	defer tx.Rollback(ctx)

	setOpts := make([]string, 0, 5)
	if banner.FeatureID != 0 {
		setOpts = append(setOpts, fmt.Sprintf("feature_id = %d", banner.FeatureID))
	}
	if banner.IsActive != nil {
		setOpts = append(setOpts, fmt.Sprintf("is_active = %t", *banner.IsActive))
	}
	if banner.Content != nil {
		setOpts = append(setOpts, fmt.Sprintf("is_active = %v", banner.Content))
	}
	setOpts = append(setOpts, "updated_at = now()")
	query = fmt.Sprintf(query, strings.Join(setOpts, ","))

	tag, err := tx.Exec(ctx, query, banner.ID)
	if err != nil {
		return NewError(errString, err)
	}
	if tag.RowsAffected() == 0 {
		return NewError(errString, repository.ErrEntityNotFound)
	}

	if banner.TagIDS != nil {
		if _, err := tx.Exec(ctx, `DELETE FROM banners_tags WHERE banner_id = $1`, banner.ID); err != nil {
			return NewError("can't update tags for banner", err)
		}
		values := make([]string, 0, len(banner.TagIDS))
		for _, tag := range banner.TagIDS {
			values = append(values, fmt.Sprintf("(%d, %d)", banner.ID, tag))
		}
		if _, err := tx.Exec(ctx, fmt.Sprintf(`INSERT INTO banners_tags (banner_id, tag_id) VALUES %s`, strings.Join(values, ","))); err != nil {
			return NewError("can't add tags for banner", err)
		}
	}

	tx.Commit(ctx)
	return nil
}

func (s BannerStorage) Get(ctx context.Context, opts repository.GetBannerLimited) ([]models.Banner, repository.DatabaseError) {
	var whereConditions, limitConditions []string
	if opts.FeatureID != 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("feature_id = %d", opts.FeatureID))
	}
	if opts.TagID != 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("tag_id = %d", opts.TagID))
	}
	if opts.Limit > 0 {
		limitConditions = append(limitConditions, fmt.Sprintf("LIMIT %d", opts.Limit))
	}
	if opts.Offset > 0 {
		limitConditions = append(limitConditions, fmt.Sprintf("OFFSET %d", opts.Offset))
	}
	var wheres string
	if len(whereConditions) != 0 {
		wheres = fmt.Sprintf("WHERE %s", strings.Join(whereConditions, " AND "))
	}

	query := fmt.Sprintf(`SELECT banners.id, feature_id, content, created_at, updated_at, is_active, array_agg(tag_id) FROM banners
	JOIN banners_tags ON banners.id = banners_tags.banner_id
	%s
	GROUP BY(banners.id)
	%s`, wheres, strings.Join(limitConditions, " "))

	rows, err := s.conn.PC.Query(ctx, query)
	if err != nil {
		return nil, NewError("can't get banners from database", err)
	}
	defer rows.Close()
	var banners []models.Banner
	for rows.Next() {
		banner := models.Banner{
			BannerCommon: models.BannerCommon{
				Content: make(models.BannerContent),
			},
		}
		var (
			createdAt, updatedAt sql.NullTime
			isActive             sql.NullBool
		)
		if err := rows.Scan(&banner.ID, &banner.FeatureID, &banner.Content, &createdAt, &updatedAt, &isActive, &banner.TagIDS); err != nil {
			return nil, NewError("can't scan banner from row", err)
		}
		if createdAt.Valid {
			banner.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			banner.UpdatedAt = updatedAt.Time
		}
		if isActive.Valid {
			banner.IsActive = &isActive.Bool
		}
		banners = append(banners, banner)
	}

	return banners, nil
}
func (s BannerStorage) GetOne(ctx context.Context, opts repository.GetBanner, userName string) (models.Banner, repository.DatabaseError) {
	/*query := `SELECT banners.id, feature_id, content, created_at, updated_at, is_active, array_agg(tag_id) AS tags FROM banners
	JOIN banners_tags ON banners.id = banners_tags.banner_id
	WHERE feature_id = $1
	GROUP BY (banners.id) HAVING $2=ANY(array_agg(tag_id))` */

	query := `SELECT b.id, feature_id, content, created_at, updated_at, is_active, array_agg(tag_id) AS tags FROM banners b 
	JOIN banners_tags bt ON b.id = bt.banner_id WHERE EXISTS 
	(SELECT u.id from users u JOIN users_tags ut ON u.id = ut.user_id where name = $1 AND (is_admin OR tag_id = $2))
	AND feature_id = $3 
	GROUP BY(b.id)`

	banner := models.Banner{
		BannerCommon: models.BannerCommon{
			Content: make(models.BannerContent),
		},
	}
	var (
		createdAt, updatedAt sql.NullTime
		isActive             sql.NullBool
	)
	if err := s.conn.PC.QueryRow(ctx, query, userName, opts.TagID, opts.FeatureID).Scan(&banner.ID, &banner.FeatureID, &banner.Content, &createdAt, &updatedAt, &isActive, &banner.TagIDS); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = repository.ErrEntityNotFound
		}
		return models.Banner{}, NewError("can't scan banner from row", err)
	}
	if createdAt.Valid {
		banner.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		banner.UpdatedAt = updatedAt.Time
	}
	if isActive.Valid {
		banner.IsActive = &isActive.Bool
	}
	return banner, nil
}

func (s BannerStorage) Delete(ctx context.Context, id int) repository.DatabaseError {
	tag, err := s.conn.PC.Exec(ctx, `DELETE FROM banners WHERE id = $1`, id)
	errString := fmt.Sprintf("can't delete banner with id %d", id)
	if err != nil {
		return NewError(errString, err)
	}
	if tag.RowsAffected() == 0 {
		return NewError(errString, repository.ErrNoRowsAffected)
	}
	return nil
}

var _ repository.BannerStorage = BannerStorage{}
