package repository

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"lang/pkg/models"
	"time"
)

type RatingPostgres struct {
	db *sqlx.DB
}

func NewRatingRepo(db *sqlx.DB) *RatingPostgres {
	return &RatingPostgres{db: db}
}

func (r *RatingPostgres) StartSession(userID int) (int, error) {
	var id int
	var startTime time.Time
	query := `INSERT INTO activity_sessions (user_id, start_time) 
              VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRow(query, userID, startTime).Scan(&id)
	if err != nil {
		return 0, errors.New("Activity session could not be created ")
	}

	return id, nil
}

func (r *RatingPostgres) EndSession(sessionID int, duration int) error {
	var endTime time.Time

	query := `UPDATE activity_sessions 
              SET end_time = $1, duration_minutes = $2 
              WHERE id = $3`
	err := r.db.QueryRow(query, endTime, duration, sessionID)
	if err != nil {
		return errors.New("Activity session could not be updated ")
	}
	return nil
}

func (r *RatingPostgres) UpdateDailyActivity(userID int, date time.Time, minutes int) error {
	// Проверяем, есть ли уже запись за этот день
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM daily_activity WHERE user_id = $1 AND date = $2)`
	err := r.db.QueryRow(query, userID, date).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		// Обновляем существующую запись
		query = `UPDATE daily_activity 
                 SET minutes_spent = minutes_spent + $1,
                     goal_achieved = (minutes_spent + $1) >= (SELECT daily_goal_minutes FROM users WHERE id = $2)
                 WHERE user_id = $2 AND date = $3`
		err = r.db.QueryRow(query, minutes, userID, date).Err()
	} else {
		// Создаем новую запись
		query = `INSERT INTO daily_activity (user_id, date, minutes_spent, goal_achieved)
                 VALUES ($1, $2, $3, $4 >= (SELECT daily_goal_minutes FROM users WHERE id = $1))`
		err = r.db.QueryRow(query, userID, date, minutes, minutes).Err()
	}

	return err
}

func (r *RatingPostgres) UpdateMonthlyStats(userID int, year, month, minutes int) error {
	// Проверяем, есть ли уже запись за этот месяц
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM monthly_stats WHERE user_id = $1 AND year = $2 AND month = $3)`
	err := r.db.QueryRow(query, userID, year, month).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		// Обновляем существующую запись
		query = `UPDATE monthly_stats 
                 SET total_minutes = total_minutes + $1,
                     days_activated = days_activated + 1
                 WHERE user_id = $2 AND year = $3 AND month = $4`
		err = r.db.QueryRow(query, minutes, userID, year, month).Err()
	} else {
		// Создаем новую запись
		query = `INSERT INTO monthly_stats (user_id, year, month, total_minutes, days_activated)
                 VALUES ($1, $2, $3, $4, 1)`
		err = r.db.QueryRow(query, userID, year, month, minutes).Err()
	}

	return err
}

func (r *RatingPostgres) GetUserRating(userID int) (*models.UserRating, error) {
	var rating models.UserRating
	var user models.User

	// Получаем основную информацию о пользователе
	query := `SELECT id, username, daily_goal_minutes FROM users WHERE id = $1`
	err := r.db.QueryRow(query, userID).Scan(&user.Id, &user.Username, &user.DailyGoal)
	if err != nil {
		return nil, err
	}

	rating.UserID = user.Id
	rating.Username = user.Username
	rating.DailyGoal = user.DailyGoal

	// Получаем общее время за последние 30 дней
	query = `SELECT COALESCE(SUM(minutes_spent), 0) 
             FROM daily_activity 
             WHERE user_id = $1 AND date >= CURRENT_DATE - INTERVAL '30 days'`
	err = r.db.QueryRow(query, userID).Scan(&rating.TotalMinutes)
	if err != nil {
		return nil, err
	}

	// Получаем текущую серию дней подряд с достижением цели
	query = `WITH dates AS (
                SELECT date 
                FROM daily_activity 
                WHERE user_id = $1 AND goal_achieved = TRUE
                ORDER BY date DESC
             )
             SELECT COUNT(*) 
             FROM (
                SELECT date, 
                       date - ROW_NUMBER() OVER (ORDER BY date) * INTERVAL '1 day' AS grp
                FROM dates
             ) t
             GROUP BY grp
             ORDER BY grp DESC
             LIMIT 1`
	err = r.db.QueryRow(query, userID).Scan(&rating.CurrentStreak)
	if err == sql.ErrNoRows {
		rating.CurrentStreak = 0
	} else if err != nil {
		return nil, err
	}

	// Получаем ежедневную статистику за последние 30 дней
	query = `SELECT date, minutes_spent, goal_achieved 
             FROM daily_activity 
             WHERE user_id = $1 AND date >= CURRENT_DATE - INTERVAL '30 days'
             ORDER BY date DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat models.DailyStat
		if err := rows.Scan(&stat.Date, &stat.MinutesSpent, &stat.GoalAchieved); err != nil {
			return nil, err
		}
		stat.Goal = user.DailyGoal
		rating.Last30Days = append(rating.Last30Days, stat)
	}

	return &rating, nil
}

func (r *RatingPostgres) GetMonthlyStats(userID int, year, month int) (*models.MonthlyStat, error) {
	var stat models.MonthlyStat

	query := `SELECT year, month, total_minutes, days_activated 
              FROM monthly_stats 
              WHERE user_id = $1 AND year = $2 AND month = $3`
	err := r.db.QueryRow(query, userID, year, month).Scan(
		&stat.Year, &stat.Month, &stat.TotalMinutes, &stat.DaysActivated)
	if err != nil {
		return nil, err
	}

	if stat.DaysActivated > 0 {
		stat.AvgDaily = float64(stat.TotalMinutes) / float64(stat.DaysActivated)
	}

	return &stat, nil
}

func (r *RatingPostgres) GetActivityTime(sessionId int) (time.Time, error) {
	var startTime time.Time
	query := `SELECT start_time FROM activity_sessions WHERE id = $1`
	err := r.db.QueryRow(query, sessionId).Scan(&startTime)
	if err != nil {
		return time.Time{}, err
	}

	return startTime, nil

}
