package repositories

import (
	models "github.com/Soup666/diss-api/model"
	"gorm.io/gorm"
)

type UserAnalyticsRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserAnalyticsRepository(db *gorm.DB) UserAnalyticsRepository {
	return &UserAnalyticsRepositoryImpl{DB: db}
}

func (r *UserAnalyticsRepositoryImpl) GetAnalytics(userID uint) (*models.UserAnalytics, error) {
	var analytics models.UserAnalytics

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Fetch user to ensure it exists
		var user models.User
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}

		// Fetch all analytics data in a single query
		rows := tx.Raw(`
			SELECT 
				(SELECT COUNT(*) FROM users WHERE id = ?) AS collection_total,
				(SELECT COUNT(*) FROM tasks WHERE user_id = ? AND completed = true) AS tasks_success,
				(SELECT COUNT(*) FROM tasks WHERE user_id = ? AND completed = false) AS tasks_failed,
				(SELECT COUNT(*) FROM tasks WHERE user_id = ?) AS tasks_total
			`, userID, userID, userID, userID).Row()

		// Scan the results into the analytics struct
		if err := rows.Scan(&analytics.CollectionTotal, &analytics.TasksSuccess, &analytics.TasksFailed, &analytics.TasksTotal); err != nil {
			return err
		}

		// Fetch weekly tasks data
		if err := tx.Raw(`
			SELECT to_char(created_at, 'dd.mm.YYYY') as date, COUNT(*) as count
			FROM tasks
			WHERE user_id = ?
			GROUP BY date
		`, userID).Scan(&analytics.WeekOfTasks).Error; err != nil {
			return err
		}

		// Fetch collections data
		if err := tx.Raw(`
				SELECT c.name as name, COUNT(ct.task_id) as count
				FROM collections c
				LEFT JOIN collection_tasks ct ON ct.collection_id = c.id
				WHERE c.user_id = ?
				GROUP BY c.name
		`, userID).Scan(&analytics.Collections).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &analytics, nil
}
