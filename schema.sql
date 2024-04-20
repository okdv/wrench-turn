-- TABLE
CREATE TABLE alert ( 
  id INTEGER PRIMARY KEY AUTOINCREMENT, 
  name TEXT, 
  description TEXT, 
  type TEXT NOT NULL DEFAULT 'notification',
  user INTEGER NOT NULL, 
  vehicle INTEGER,
  job INTEGER,  
  task INTEGER, 
  is_read INTEGER NOT NULL DEFAULT 0, 
  read_at DATETIME,
  alert_at DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE job(
  id INTEGER PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  instructions TEXT,
  is_template INTEGER DEFAULT 0 NOT NULL,
  is_complete INTEGER DEFAULT 0 NOT NULL,
  vehicle INTEGER,
  user INTEGER,
  origin_job INTEGER,
  repeats INTEGER DEFAULT 0 NOT NULL, 
  odo_interval INTEGER, 
  time_interval INTEGER,
  time_interval_unit VARCHAR(10),
  due_date DATETIME, 
  completed_at DATETIME,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
  );
CREATE TABLE job_label ( id INTEGER PRIMARY KEY AUTOINCREMENT, job INTEGER NOT NULL, label INTEGER NOT NULL );
CREATE TABLE label ( 
  id INTEGER PRIMARY KEY AUTOINCREMENT, 
  name TEXT NOT NULL, 
  color TEXT,
  user INTEGER, 
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE task ( 
  id INTEGER PRIMARY KEY AUTOINCREMENT, 
  name TEXT, 
  description TEXT, 
  is_complete INTEGER NOT NULL DEFAULT 0, 
  job INTEGER NOT NULL, 
  part_name TEXT,
  part_link TEXT,
  due_date DATETIME,
  completed_at DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE user(
  id INTEGER PRIMARY KEY NOT NULL,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE,
  description TEXT,
  hashed_pw BLOB,
  is_admin INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );
CREATE TABLE vehicle ( 
  id INTEGER PRIMARY KEY AUTOINCREMENT, 
  name TEXT, 
  description TEXT, 
  type TEXT, 
  is_metric INTEGER NOT NULL DEFAULT 0, 
  vin TEXT, 
  year INTEGER, 
  make TEXT,
  model TEXT, 
  trim TEXT,
  odometer INTEGER, 
  user INTEGER NOT NULL, 
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
 
-- INDEX
CREATE INDEX alert_at_user_idx ON alert (user, alert_at);
CREATE INDEX alert_user_idx ON alert (user);
CREATE INDEX job_label_job_idx ON job_label (job);
CREATE INDEX job_user_idx ON job (user);
CREATE INDEX job_vehicle_idx ON job (vehicle);
CREATE INDEX label_user_idx ON label (user);
CREATE INDEX task_job_idx ON task (job);
CREATE INDEX username_idx ON user (username);
CREATE INDEX vehicle_user_idx ON vehicle (user);
 
-- TRIGGER
 
-- VIEW
 
