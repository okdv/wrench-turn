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
func ListJobs(userId *string, vehicleId *string, isTemplate *string, labelId *string, searchStr *string, sort *string) ([]*models.Job, error) {
	var joins []string
	var wheres []string
	var likes []Like
	// establish default sort if not provided
	var orderBy = "j.updated_at DESC"
	// establish basic query
	q := "SELECT j.ID, j.Name, j.Description, j.Instructions, j.Is_template, j.Is_complete, j.Vehicle, j.User, j.Origin_job, j.Repeats, j.Odo_interval, j.Time_interval, j.Time_interval_unit, j.Due_date, j.Completed_at, j.Created_at, j.Updated_at FROM job AS j"
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
	// if label ID provided join by labelId where jobId is present
	if labelId != nil && len(*labelId) > 0 {
		joins = append(joins, "JOIN job_label AS jl ON j.id = jl.job")
		wheres = append(wheres, "jl.label="+*labelId)
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

// Task Queries

// GetTask
// Takes job id, queries it in db, returns Task
func GetTask(jobId int64, taskId int64) (*models.Task, error) {
	var task models.Task
	// query db, return any errors
	err := DB.QueryRow("SELECT * FROM task WHERE id=? AND job=?", taskId, jobId).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Is_complete,
		&task.Job,
		&task.Part_name,
		&task.Part_link,
		&task.Due_date,
		&task.Completed_at,
		&task.Created_at,
		&task.Updated_at,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	return &task, nil
}

// CreateTask
// Takes newTask, creates in db, returns id
func CreateTask(newTask models.NewTask, jobId int64) (*int64, error) {
	q := "INSERT INTO task(Name, Description, Job, Part_name, Part_link, Due_date) VALUES (?,?,?,?,?,?)"
	log.Printf(q)
	// insert into db, return any errors
	res, err := DB.Exec("INSERT INTO task(Name, Description, Job, Part_name, Part_link, Due_date) VALUES (?,?,?,?,?,?)",
		newTask.Name,
		newTask.Description,
		jobId,
		newTask.Part_name,
		newTask.Part_link,
		newTask.Due_date,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	// get inserted tasks id
	taskId, err := res.LastInsertId()
	return &taskId, err
}

// EditTask
// Take Task as arg, build update query with QueryBuilder, update it in db via generated query
func EditTask(editedTask models.Task, jobId int64) error {
	var wheres []string
	// setup query
	q := "UPDATE task SET name=?, description=?, part_name=?, part_link=?, due_date=?, updated_at=CURRENT_TIMESTAMP"
	// add required wheres (ensures the task id and user id in the db match that of request body)
	wheres = append(wheres, "job=?")
	wheres = append(wheres, "id=?")
	// get generated query
	query := QueryBuilder(q, nil, &wheres, nil, nil)
	// exec query
	res, err := DB.Exec(query, editedTask.Name, editedTask.Description, editedTask.Part_name, editedTask.Part_link, editedTask.Due_date, jobId, editedTask.ID)
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

// UpdateTaskStatus
// Take job id, task id, status as args, build update query with QueryBuilder, update it in db via generated query
func UpdateTaskStatus(jobId int64, taskId int64, status int) error {
	var wheres []string
	// setup query
	q := "UPDATE task SET is_complete=?, updated_at=CURRENT_TIMESTAMP"
	// if status is complete, updated completed_at also
	if status == 1 {
		q += ", completed_at=CURRENT_TIMESTAMP"
	}
	// add required wheres (ensures the task id and user id in the db match that of request body)
	wheres = append(wheres, "job=?")
	wheres = append(wheres, "id=?")
	// get generated query
	query := QueryBuilder(q, nil, &wheres, nil, nil)
	// exec query
	res, err := DB.Exec(query, status, jobId, taskId)
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

// DeleteTask
// Take task id as arg, delete Task from task table where id present
func DeleteTask(jobId int64, taskId int64) error {
	res, err := DB.Exec("DELETE FROM task WHERE id=? AND job=?", taskId, jobId)
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

// ListTasks
// Take filters as args, return Task list
func ListTasks(jobId int64, isComplete *string, searchStr *string, sort *string) ([]*models.Task, error) {
	var joins []string
	var wheres []string
	var likes []Like
	// establish default sort if not provided
	var orderBy = "t.updated_at DESC"
	// establish basic query
	q := "SELECT * FROM task AS t"
	// if isTemplate provided, add where to query
	if isComplete != nil && len(*isComplete) > 0 {
		wheres = append(wheres, "t.is_complete="+*isComplete)
	}
	// add wheres for job id
	wheres = append(wheres, "t.job="+strconv.FormatInt(jobId, 10))
	// if search string provided, construct likes to query username, description cols
	if searchStr != nil && len(*searchStr) > 0 {
		var fields []string
		fields = append(fields, "t.name")
		fields = append(fields, "t.description")
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
			orderBy = "t.name ASC"
		case "za":
			orderBy = "t.name DESC"
		case "completed":
			orderBy = "t.completed_at DESC"
		case "oldest":
			orderBy = "t.created_at ASC"
		case "newest":
			orderBy = "t.created_at DESC"
		case "last_updated":
			orderBy = "t.updated_at DESC"
		default:
			orderBy = "t.updated_at DESC"
		}
	}
	// generate query with QueryBuilder
	query := QueryBuilder(q, &joins, &wheres, &likes, &orderBy)
	log.Printf(query)
	// retrieve all matching rows
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("DB Query Error: %s", err)
		return nil, err
	}
	defer rows.Close()
	// create list of Task
	tasks := make([]*models.Task, 0)
	// loop through returned rows
	for rows.Next() {
		// attribute to Task
		task := models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.Is_complete,
			&task.Job,
			&task.Part_name,
			&task.Part_link,
			&task.Due_date,
			&task.Completed_at,
			&task.Created_at,
			&task.Updated_at,
		)
		if err != nil {
			log.Printf("Error scanning rows retrieved from DB: %s", err)
			return nil, err
		}
		// append Task to list of Task
		tasks = append(tasks, &task)
	}
	return tasks, nil
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

// Alert Queries

// GetAlert
// Takes alert id, queries it in db, returns Alert
func GetAlert(alertId int64) (*models.Alert, error) {
	var alert models.Alert
	// query db, return any errors
	err := DB.QueryRow("SELECT * FROM alert WHERE id=?", alertId).Scan(
		&alert.ID,
		&alert.Name,
		&alert.Description,
		&alert.Type,
		&alert.User,
		&alert.Vehicle,
		&alert.Job,
		&alert.Task,
		&alert.Is_read,
		&alert.Read_at,
		&alert.Alert_at,
		&alert.Created_at,
		&alert.Updated_at,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	return &alert, nil
}

// CreateAlert
// Takes newAlert, creates in db, returns id
func CreateAlert(newAlert models.NewAlert) (*int64, error) {
	// insert into db, return any errors
	res, err := DB.Exec("INSERT INTO alert(Name, Description, User, Vehicle, Job, Task, Alert_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		newAlert.Name,
		newAlert.Description,
		newAlert.User,
		newAlert.Vehicle,
		newAlert.Job,
		newAlert.Task,
		newAlert.Alert_at,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	// get inserted alerts id
	alertId, err := res.LastInsertId()
	return &alertId, err
}

// EditAlert
// Take Alert as arg, build update query with QueryBuilder, update it in db via generated query
func EditAlert(editedAlert models.Alert) error {
	var wheres []string
	// setup query
	q := "UPDATE alert SET name=?, description=?, type=?, user=?, vehicle=?, job=?, task=?, is_read=?, alert_at=?, updated_at=CURRENT_TIMESTAMP"
	// add required wheres (ensures the alert id and user id in the db match that of request body)
	wheres = append(wheres, "user=?")
	wheres = append(wheres, "id=?")
	// get generated query
	query := QueryBuilder(q, nil, &wheres, nil, nil)
	// exec query
	res, err := DB.Exec(query, editedAlert.Name, editedAlert.Description, editedAlert.Type, editedAlert.User, editedAlert.Vehicle, editedAlert.Job, editedAlert.Task, editedAlert.Is_read, editedAlert.Alert_at, editedAlert.User, editedAlert.ID)
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

// DeleteAlert
// Take alert id as arg, delete Alert from alert table where id present
func DeleteAlert(alertId int64, userId *int64) error {
	var wheres []string
	q := "DELETE FROM alert"
	wheres = append(wheres, "id="+strconv.FormatInt(alertId, 10))
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

// ListAlerts
// Take filters as args, return Alert list
func ListAlerts(userId *string, vehicleId *string, jobId *string, taskId *string, typeStr *string, isRead *string, alertDate *string, searchStr *string, sort *string) ([]*models.Alert, error) {
	var joins []string
	var wheres []string
	var likes []Like
	// establish default sort if not provided
	var orderBy = "a.updated_at DESC"
	// establish basic query
	q := "SELECT * FROM alert AS a"
	// if userId provided, add where to query
	if userId != nil && len(*userId) > 0 {
		wheres = append(wheres, "a.user="+*userId)
	}
	// if vehicleId provided, addawhere to query
	if vehicleId != nil && len(*vehicleId) > 0 {
		wheres = append(wheres, "a.vehicle="+*vehicleId)
	}
	// if jobId provided, addawhere to query
	if jobId != nil && len(*jobId) > 0 {
		wheres = append(wheres, "a.job="+*jobId)
	}
	// if taskId provided, addawhere to query
	if taskId != nil && len(*taskId) > 0 {
		wheres = append(wheres, "a.task="+*taskId)
	}
	// if typeStr provided, addawhere to query
	if typeStr != nil && len(*typeStr) > 0 {
		wheres = append(wheres, "a.type="+*typeStr)
	}
	// if isRead provided, add where to query
	if isRead != nil && len(*isRead) > 0 {
		wheres = append(wheres, "a.is_read="+*isRead)
	}
	// if alertDate provided, add where to query
	if alertDate != nil {
		wheres = append(wheres, "a.alert_at<='"+*alertDate+"'")
	}
	// if search string provided, construct likes to query username, description cols
	if searchStr != nil && len(*searchStr) > 0 {
		var fields []string
		fields = append(fields, "a.name")
		fields = append(fields, "a.description")
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
			orderBy = "a.name ASC"
		case "za":
			orderBy = "a.name DESC"
		case "completed":
			orderBy = "a.completed_at DESC"
		case "oldest":
			orderBy = "a.created_at ASC"
		case "newest":
			orderBy = "a.created_at DESC"
		case "last_updated":
			orderBy = "a.updated_at DESC"
		default:
			orderBy = "a.updated_at DESC"
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
	// create list of Alert
	alerts := make([]*models.Alert, 0)
	// loop through returned rows
	for rows.Next() {
		// attribute to Alert
		alert := models.Alert{}
		err := rows.Scan(
			&alert.ID,
			&alert.Name,
			&alert.Description,
			&alert.Type,
			&alert.User,
			&alert.Vehicle,
			&alert.Job,
			&alert.Task,
			&alert.Is_read,
			&alert.Read_at,
			&alert.Alert_at,
			&alert.Created_at,
			&alert.Updated_at,
		)
		if err != nil {
			log.Printf("Error scanning rows retrieved from DB: %s", err)
			return nil, err
		}
		// append Alert to list of Alert
		alerts = append(alerts, &alert)
	}
	return alerts, nil
}

// UpdatedAlertStatus
// Take alert id, user id, status as args, build update query with QueryBuilder, update it in db via generated query
func UpdatedAlertStatus(alertId int64, userId int64, status int) error {
	var wheres []string
	// setup query
	q := "UPDATE alert SET is_read=?, updated_at=CURRENT_TIMESTAMP"
	// if status is complete, updated completed_at also
	if status == 1 {
		q += ", read_at=CURRENT_TIMESTAMP"
	}
	// add required wheres (ensures the alert id and user id in the db match that of request body)
	wheres = append(wheres, "user=?")
	wheres = append(wheres, "id=?")
	// get generated query
	query := QueryBuilder(q, nil, &wheres, nil, nil)
	// exec query
	res, err := DB.Exec(query, status, userId, alertId)
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

// Label Queries

// GetLabel
// Takes label id, queries it in db, returns Label
func GetLabel(labelId int64) (*models.Label, error) {
	var label models.Label
	// query db, return any errors
	err := DB.QueryRow("SELECT * FROM label WHERE id=?", labelId).Scan(
		&label.ID,
		&label.Name,
		&label.Color,
		&label.User,
		&label.Created_at,
		&label.Updated_at,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	return &label, nil
}

// CreateLabel
// Takes newLabel, creates in db, returns id
func CreateLabel(newLabel models.NewLabel) (*int64, error) {
	// insert into db, return any errors
	res, err := DB.Exec("INSERT INTO label(Name, Color, User) VALUES (?,?,?)",
		newLabel.Name,
		newLabel.Color,
		newLabel.User,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	// get inserted labels id
	labelId, err := res.LastInsertId()
	return &labelId, err
}

// EditLabel
// Take Label as arg, build update query with QueryBuilder, update it in db via generated query
func EditLabel(editedLabel models.Label) error {
	var wheres []string
	// setup query
	q := "UPDATE label SET name=?, color=?, updated_at=CURRENT_TIMESTAMP"
	// add required wheres (ensures the label id and user id in the db match that of request body)
	wheres = append(wheres, "user=?")
	wheres = append(wheres, "id=?")
	// get generated query
	query := QueryBuilder(q, nil, &wheres, nil, nil)
	// exec query
	res, err := DB.Exec(query, editedLabel.Name, editedLabel.Color, editedLabel.User, editedLabel.ID)
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

// DeleteLabel
// Take label id as arg, delete Label from label table where id present
func DeleteLabel(labelId int64, userId *int64) error {
	var wheres []string
	q := "DELETE FROM label"
	wheres = append(wheres, "id="+strconv.FormatInt(labelId, 10))
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

// ListLabels
// Take filters as args, return Label list
func ListLabels(userId *string, jobId *string, searchStr *string, sort *string) ([]*models.Label, error) {
	var joins []string
	var wheres []string
	var likes []Like
	// establish default sort if not provided
	var orderBy = "l.updated_at DESC"
	// establish basic query
	q := "SELECT l.id, l.name, l.color, l.user, l.created_at, l.updated_at FROM label AS l"
	// if userId provided, add where to query
	if userId != nil && len(*userId) > 0 {
		wheres = append(wheres, "l.user="+*userId)
	}
	// if job ID provided join by jobID where vehicleID is present
	if jobId != nil && len(*jobId) > 0 {
		joins = append(joins, "JOIN job_label AS jl ON l.id = jl.label")
		wheres = append(wheres, "jl.job="+*jobId)
	}
	// if search string provided, construct likes to query username, description cols
	if searchStr != nil && len(*searchStr) > 0 {
		var fields []string
		fields = append(fields, "l.name")
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
			orderBy = "l.name ASC"
		case "za":
			orderBy = "l.name DESC"
		case "oldest":
			orderBy = "l.created_at ASC"
		case "newest":
			orderBy = "l.created_at DESC"
		case "last_updated":
			orderBy = "l.updated_at DESC"
		default:
			orderBy = "l.updated_at DESC"
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
	// create list of Label
	labels := make([]*models.Label, 0)

	// loop through returned rows
	for rows.Next() {
		// attribute to Label
		label := models.Label{}
		err := rows.Scan(
			&label.ID,
			&label.Name,
			&label.Color,
			&label.User,
			&label.Created_at,
			&label.Updated_at,
		)
		if err != nil {
			log.Printf("Error scanning rows retrieved from DB: %s", err)
			return nil, err
		}
		// append Label to list of Label
		labels = append(labels, &label)
	}
	return labels, nil
}

// AssignJobLabel
// Takes job id and task id, creates job_label in db, returns id
func AssignJobLabel(jobId int64, labelId int64) (*int64, error) {
	q := "INSERT INTO job_label(Job, Label) VALUES (?,?)"
	log.Printf(q)
	// insert into db, return any errors
	res, err := DB.Exec(q,
		jobId,
		labelId,
	)
	if err != nil {
		log.Printf("DB Execution Error: %s", err)
		return nil, err
	}
	// get inserted tasks id
	taskId, err := res.LastInsertId()
	return &taskId, err
}

// UnassignJobLabel
// Takes job id and task id, deletes job_label entry in db if its exists
func UnassignJobLabel(jobId int64, labelId int64) error {
	var wheres []string
	q := "DELETE FROM job_label"
	wheres = append(wheres, "job="+strconv.FormatInt(jobId, 10))
	wheres = append(wheres, "label="+strconv.FormatInt(labelId, 10))
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
