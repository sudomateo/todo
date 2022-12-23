CREATE TABLE todos (
	id uuid,
	text text,
	priority text,
	completed bool,
	time_created timestamp,
	time_updated timestamp,

	PRIMARY KEY (id)
);

CREATE INDEX priority_index ON todos (priority);
CREATE INDEX completed_index ON todos (completed);
CREATE INDEX time_created_index ON todos (time_created);
