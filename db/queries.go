package db

import (
	"errors"
	"log"
	"strconv"

	"github.com/okdv/wrench-turn/models"
)

// Auth Queries

// GetAuthInfoByUsername
// Take username, retrieve auth info: user ID, isAdmin and the hashed password
func GetAuthInfoByUsername(username string) (*int64, *string, *int, *[]byte, error) {
	var userId int64
	var isAdmin int
	var hashed []byte
	err := DB.QueryRow("SELECT id, is_admin, hashed_pw FROM user WHERE username = ?", username).Scan(&userId, &isAdmin, &hashed)
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return nil, nil, nil, nil, err
	}
	return &userId, &username, &isAdmin, &hashed, nil
}

// User Queries

// GetUserById
// Take user ID, return entire user
func GetUserById(userId int64) (*models.User, error) {
	var user models.User
	// query db, return any errors
	err := DB.QueryRow("SELECT * FROM user WHERE id=?", userId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Description,
		&user.Hashed_pw,
		&user.Is_admin,
		&user.Created_at,
		&user.Updated_at,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername
// Take username, return entire user
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	// query db, return any errors
	err := DB.QueryRow("SELECT * FROM user WHERE username=?", username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Description,
		&user.Hashed_pw,
		&user.Is_admin,
		&user.Created_at,
		&user.Updated_at,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	return &user, nil
}

// ListUsers
// Take filters as args, return User list
func ListUsers(jobId *string, vehicleId *string, isAdmin *string, searchStr *string, sort *string) ([]*models.User, error) {
	var joins []string
	var wheres []string
	var likes []Like
	// establish default sort if not provided
	var orderBy = "u.updated_at DESC"
	// establish basic query
	q := "SELECT * FROM user AS u"
	// if isAdmin provided, add where to query
	if isAdmin != nil && len(*isAdmin) > 0 {
		wheres = append(wheres, "u.is_admin="+*isAdmin)
	}
	// if job ID provided join by userID where jobID is present
	if jobId != nil && len(*jobId) > 0 {
		joins = append(joins, "JOIN user_job AS uj ON u.id = uj.user")
		wheres = append(wheres, "uj.job="+*jobId)
	}
	// if vehicle ID provided join by userID where vehicleID is present
	if vehicleId != nil && len(*vehicleId) > 0 {
		joins = append(joins, "JOIN user_vehicle AS uv ON u.id = uv.user")
		wheres = append(wheres, "uv.vehicle="+*vehicleId)
	}
	// if search string provided, construct likes to query username, description cols
	if searchStr != nil && len(*searchStr) > 0 {
		var fields []string
		fields = append(fields, "u.username")
		fields = append(fields, "u.description")
		likes = append(likes, Like{
			Fields: fields,
			Match:  *searchStr,
			Or:     true,
		})
	}
	// if sort provided, append appropriate sort based on query param
	if sort != nil && len(*sort) > 0 {
		switch *sort {
		case "az":
			orderBy = "u.username ASC"
		case "za":
			orderBy = "u.username DESC"
		case "oldest":
			orderBy = "u.created_at ASC"
		case "newest":
			orderBy = "u.created_at DESC"
		case "last_updated":
			orderBy = "u.updated_at DESC"
		default:
			orderBy = "u.updated_at DESC"
		}
	}
	// generate query with QueryBuilder
	query := QueryBuilder(q, &joins, &wheres, &likes, &orderBy)
	// retrieve all matching rows
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return nil, err
	}
	defer rows.Close()
	// create list of User
	users := make([]*models.User, 0)
	// loop through returned rows
	for rows.Next() {
		// attribute to User
		user := models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Description,
			&user.Hashed_pw,
			&user.Is_admin,
			&user.Created_at,
			&user.Updated_at,
		)
		if err != nil {
			log.Printf("Error scanning rows retrieved from DB: %s", err)
			return nil, err
		}
		// append User to list of User
		users = append(users, &user)
	}
	return users, nil
}

// CreateUser
// Take NewUser and hashed PW as arguments, insert them into db
func CreateUser(newUser models.NewUser, password *[]byte) (*int64, error) {
	// insert into db, return any errors
	res, err := DB.Exec("INSERT INTO user(Username, Email, Hashed_pw, Is_admin) VALUES (?,?,?,?)",
		newUser.Username,
		newUser.Email,
		password,
		newUser.Is_admin,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	// get inserted users id
	userId, err := res.LastInsertId()
	return &userId, err
}

// EditUser
// Take User as arg, update it in db
func EditUser(editedUser models.User) error {
	_, err := DB.Exec("UPDATE user SET username=?, email=?, description=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		editedUser.Username,
		editedUser.Email,
		editedUser.Description,
		editedUser.ID,
	)
	// throw SQL errors
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return err
	}
	return err
}

// DeleteUser
// Take username as arg, delete User from user table where username present
func DeleteUser(username string) error {
	res, err := DB.Exec("DELETE FROM user WHERE username = ?", username)
	// throw SQL errors
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return err
	}
	// retrieve rows affected count
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return err
	}
	// throw error if no rows affected
	if rows == 0 {
		log.Printf("No rows deleted")
		return errors.New("No rows deleted")
	}

	return nil
}

// UpdatePassword
// Take username and hashed pw as args, update it in db
func UpdatePassword(username string, password *[]byte) error {
	res, err := DB.Exec("UPDATE user SET hashed_pw=? WHERE username=?", password, username)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return err
	}
	// retrieve rows affected count, error if 0
	rowCount, err := res.RowsAffected()
	if rowCount == 0 || err != nil {
		log.Printf("No rows updated: %v", err)
		return errors.New("No rows updated")
	}
	return nil
}

// Job Queries

// GetJob
// Takes job id, queries it in db, returns Job
func GetJob(jobId int64) (*models.Job, error) {
	var job models.Job
	// query db, return any errors
	err := DB.QueryRow("SELECT * FROM job WHERE id=?", jobId).Scan(
		&job.ID,
		&job.Name,
		&job.Description,
		&job.Instructions,
		&job.Is_template,
		&job.Is_complete,
		&job.Vehicle,
		&job.User,
		&job.Origin_job,
		&job.Repeats,
		&job.Odo_interval,
		&job.Time_interval,
		&job.Time_interval_unit,
		&job.Due_date,
		&job.Completed_at,
		&job.Created_at,
		&job.Updated_at,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	return &job, nil
}

// CreateJob
// Takes newJob, creates in db, returns id
func CreateJob(newJob models.NewJob) (*int64, error) {
	// insert into db, return any errors
	res, err := DB.Exec("INSERT INTO job(Name, Description, Instructions, Is_template, Vehicle, User, Origin_job, Repeats, Odo_interval, Time_interval, Time_interval_unit, Due_date) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
		newJob.Name,
		newJob.Description,
		newJob.Instructions,
		newJob.Is_template,
		newJob.Vehicle,
		newJob.User,
		newJob.Origin_job,
		newJob.Repeats,
		newJob.Odo_interval,
		newJob.Time_interval,
		newJob.Time_interval_unit,
		newJob.Due_date,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	// get inserted jobs id
	jobId, err := res.LastInsertId()
	return &jobId, err
}

// EditJob
// Take Job as arg, build update query with QueryBuilder, update it in db via generated query
func EditJob(editedJob models.Job) error {
	var wheres []string
	// setup query
	q := "UPDATE job SET name=?, description=?, instructions=?, is_template=?, repeats=?, odo_interval=?, time_interval=?, time_interval_unit=?, due_date=?, updated_at=CURRENT_TIMESTAMP"
	// add required wheres (ensures the job id and user id in the db match that of request body)
	wheres = append(wheres, "user=?")
	wheres = append(wheres, "id=?")
	// get generated query
	query := QueryBuilder(q, nil, &wheres, nil, nil)
	// exec query
	res, err := DB.Exec(query, editedJob.Name, editedJob.Description, editedJob.Instructions, editedJob.Is_template, editedJob.Repeats, editedJob.Odo_interval, editedJob.Time_interval, editedJob.Time_interval_unit, editedJob.Due_date, editedJob.User, editedJob.ID)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return err
	}
	// retrieve rows affected count, error if 0
	rowCount, err := res.RowsAffected()
	if rowCount == 0 || err != nil {
		log.Printf("No rows updated: %v", err)
		return errors.New("No rows updated")
	}
	return nil
}

// DeleteJob
// Take job id as arg, delete Job from job table where id present
func DeleteJob(jobId int64, userId *int64) error {
	var wheres []string
	q := "DELETE FROM job"
	wheres = append(wheres, "id="+strconv.FormatInt(jobId, 10))
	if userId != nil {
		wheres = append(wheres, "user="+strconv.FormatInt(*userId, 10))
	}
	query := QueryBuilder(q, nil, &wheres, nil, nil)
	res, err := DB.Exec(query)
	// throw SQL errors
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return err
	}
	// retrieve rows affected count
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return err
	}
	// throw error if no rows affected
	if rows == 0 {
		log.Printf("No rows deleted")
		return errors.New("No rows deleted")
	}

	return nil
}

// ListJobs
// Take filters as args, return Job list
func ListJobs(userId *string, vehicleId *string, isTemplate *string, searchStr *string, sort *string) ([]*models.Job, error) {
	var joins []string
	var wheres []string
	var likes []Like
	// establish default sort if not provided
	var orderBy = "j.updated_at DESC"
	// establish basic query
	q := "SELECT * FROM job AS j"
	// if userId provided, add where to query
	if userId != nil && len(*userId) > 0 {
		wheres = append(wheres, "j.user="+*userId)
	}
	// if vehicleId provided, add where to query
	if vehicleId != nil && len(*vehicleId) > 0 {
		wheres = append(wheres, "j.vehicle="+*vehicleId)
	}
	// if isTemplate provided, add where to query
	if isTemplate != nil && len(*isTemplate) > 0 {
		wheres = append(wheres, "j.is_template="+*isTemplate)
	}
	// if search string provided, construct likes to query username, description cols
	if searchStr != nil && len(*searchStr) > 0 {
		var fields []string
		fields = append(fields, "j.name")
		fields = append(fields, "j.description")
		fields = append(fields, "j.instructions")
		likes = append(likes, Like{
			Fields: fields,
			Match:  *searchStr,
			Or:     true,
		})
	}
	// if sort provided, append appropriate sort based on query param
	if sort != nil && len(*sort) > 0 {
		switch *sort {
		case "az":
			orderBy = "j.name ASC"
		case "za":
			orderBy = "j.name DESC"
		case "completed":
			orderBy = "j.completed_at DESC"
		case "oldest":
			orderBy = "j.created_at ASC"
		case "newest":
			orderBy = "j.created_at DESC"
		case "last_updated":
			orderBy = "j.updated_at DESC"
		default:
			orderBy = "j.updated_at DESC"
		}
	}
	// generate query with QueryBuilder
	query := QueryBuilder(q, &joins, &wheres, &likes, &orderBy)
	// retrieve all matching rows
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return nil, err
	}
	defer rows.Close()
	// create list of Job
	jobs := make([]*models.Job, 0)
	// loop through returned rows
	for rows.Next() {
		// attribute to Job
		job := models.Job{}
		err := rows.Scan(
			&job.ID,
			&job.Name,
			&job.Description,
			&job.Instructions,
			&job.Is_template,
			&job.Is_complete,
			&job.Vehicle,
			&job.User,
			&job.Origin_job,
			&job.Repeats,
			&job.Odo_interval,
			&job.Time_interval,
			&job.Time_interval_unit,
			&job.Due_date,
			&job.Completed_at,
			&job.Created_at,
			&job.Updated_at,
		)
		if err != nil {
			log.Printf("Error scanning rows retrieved from DB: %s", err)
			return nil, err
		}
		// append Job to list of Job
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

// Vehicle Queries

// GetVehicle
// Takes vehicle id, queries it in db, returns vehicle
func GetVehicle(vehicleId int64) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	// query db, return any errors
	err := DB.QueryRow("SELECT * FROM vehicle WHERE id=?", vehicleId).Scan(
		&vehicle.ID,
		&vehicle.Name,
		&vehicle.Description,
		&vehicle.Type,
		&vehicle.Is_metric,
		&vehicle.Vin,
		&vehicle.Year,
		&vehicle.Make,
		&vehicle.Model,
		&vehicle.Trim,
		&vehicle.Odometer,
		&vehicle.User,
		&vehicle.Created_at,
		&vehicle.Updated_at,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	return &vehicle, nil
}

// ListVehicles
// Take filters as args, return Vehicle list
func ListVehicles(userId *string, jobId *string, searchStr *string, sort *string) ([]*models.Vehicle, error) {
	var joins []string
	var wheres []string
	var likes []Like
	// establish default sort if not provided
	var orderBy = "v.updated_at DESC"
	// establish basic query
	q := "SELECT * FROM vehicle AS v"
	// if userId provided, add where to query
	if userId != nil && len(*userId) > 0 {
		wheres = append(wheres, "v.user="+*userId)
	}
	// if job ID provided join by jobID where vehicleID is present
	if jobId != nil && len(*jobId) > 0 {
		joins = append(joins, "JOIN job AS j ON v.id = j.vehicle")
		wheres = append(wheres, "j.id="+*jobId)
	}
	// if search string provided, construct likes to query username, description cols
	if searchStr != nil && len(*searchStr) > 0 {
		var fields []string
		fields = append(fields, "v.name")
		fields = append(fields, "v.description")
		fields = append(fields, "v.make")
		fields = append(fields, "v.model")
		fields = append(fields, "v.trim")
		likes = append(likes, Like{
			Fields: fields,
			Match:  *searchStr,
			Or:     true,
		})
	}
	// if sort provided, append appropriate sort based on query param
	if sort != nil && len(*sort) > 0 {
		switch *sort {
		case "az":
			orderBy = "v.name ASC"
		case "za":
			orderBy = "v.name DESC"
		case "completed":
			orderBy = "v.completed_at DESC"
		case "oldest":
			orderBy = "v.created_at ASC"
		case "newest":
			orderBy = "v.created_at DESC"
		case "last_updated":
			orderBy = "v.updated_at DESC"
		default:
			orderBy = "v.updated_at DESC"
		}
	}
	// generate query with QueryBuilder
	query := QueryBuilder(q, &joins, &wheres, &likes, &orderBy)
	// retrieve all matching rows
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return nil, err
	}
	defer rows.Close()
	// create list of Vehicle
	vehicles := make([]*models.Vehicle, 0)
	// loop through returned rows
	for rows.Next() {
		// attribute to Vehicle
		vehicle := models.Vehicle{}
		err := rows.Scan(
			&vehicle.ID,
			&vehicle.Name,
			&vehicle.Description,
			&vehicle.Type,
			&vehicle.Is_metric,
			&vehicle.Vin,
			&vehicle.Year,
			&vehicle.Make,
			&vehicle.Model,
			&vehicle.Trim,
			&vehicle.Odometer,
			&vehicle.User,
			&vehicle.Created_at,
			&vehicle.Updated_at,
		)
		if err != nil {
			log.Printf("Error scanning rows retrieved from DB: %s", err)
			return nil, err
		}
		// append Job to list of Job
		vehicles = append(vehicles, &vehicle)
	}
	return vehicles, nil
}

// CreateVehicle
// Takes newVehicle, creates in db, returns id
func CreateVehicle(newVehicle models.NewVehicle) (*int64, error) {
	// insert into db, return any errors
	res, err := DB.Exec("INSERT INTO vehicle(Name, Description, Type, Is_metric, Vin, Year, Make, Model, Trim, Odometer, User) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		newVehicle.Name,
		newVehicle.Description,
		newVehicle.Type,
		newVehicle.Is_metric,
		newVehicle.Vin,
		newVehicle.Year,
		newVehicle.Make,
		newVehicle.Model,
		newVehicle.Trim,
		newVehicle.Odometer,
		newVehicle.User,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	// get inserted vehicles id
	vehicleId, err := res.LastInsertId()
	return &vehicleId, err
}

// EditVehicle
// Take Vehicle as arg, build update query with QueryBuilder, update it in db via generated query
func EditVehicle(editedVehicle models.Vehicle) error {
	var wheres []string
	// setup query
	q := "UPDATE vehicle SET name=?, description=?, type=?, is_metric=?, vin=?, year=?, make=?, model=?, trim=?, odometer=?, user=?, updated_at=CURRENT_TIMESTAMP"
	// add required wheres (ensures the vehicle id and user id in the db match that of request body)
	wheres = append(wheres, "user=?")
	wheres = append(wheres, "id=?")
	// get generated query
	query := QueryBuilder(q, nil, &wheres, nil, nil)
	// exec query
	res, err := DB.Exec(query, editedVehicle.Name, editedVehicle.Description, editedVehicle.Type, editedVehicle.Is_metric, editedVehicle.Vin, editedVehicle.Year, editedVehicle.Make, editedVehicle.Model, editedVehicle.Trim, editedVehicle.Odometer, editedVehicle.User, editedVehicle.User, editedVehicle.ID)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return err
	}
	// retrieve rows affected count, error if 0
	rowCount, err := res.RowsAffected()
	if rowCount == 0 || err != nil {
		log.Printf("No rows updated: %v", err)
		return errors.New("No rows updated")
	}
	return nil
}

// DeleteVehicle
// Take vehicle id as arg, delete Vehicle from vehicle table where id present
func DeleteVehicle(vehicleId int64, userId *int64) error {
	var wheres []string
	q := "DELETE FROM vehicle"
	wheres = append(wheres, "id="+strconv.FormatInt(vehicleId, 10))
	if userId != nil {
		wheres = append(wheres, "user="+strconv.FormatInt(*userId, 10))
	}
	query := QueryBuilder(q, nil, &wheres, nil, nil)
	res, err := DB.Exec(query)
	// throw SQL errors
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return err
	}
	// retrieve rows affected count
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return err
	}
	// throw error if no rows affected
	if rows == 0 {
		log.Printf("No rows deleted")
		return errors.New("No rows deleted")
	}

	return nil
}
