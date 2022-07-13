CREATE TABLE "user"(
	id VARCHAR(50) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
	creation_date TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE "microblog"(
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(50),
    visit_count BIGINT NOT NULL DEFAULT 0,
    text VARCHAR(200) NOT NULL, 
    creation_date TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT microblog_user_fk
      FOREIGN KEY(user_id) 
	  REFERENCES "user"(id)
);

CREATE TABLE "microblog_tracking"(
    microblog_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT microblog_tracking_user_fk
      FOREIGN KEY(user_id) 
	  REFERENCES "user" (id),
    CONSTRAINT microblog_tracking_microblog_fk
      FOREIGN KEY(microblog_id) 
	  REFERENCES microblog(id)
);

CREATE UNIQUE INDEX microblog_tracking_uq_1 ON microblog_tracking (microblog_id, user_id);
