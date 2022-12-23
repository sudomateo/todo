CREATE TABLE todos (
	id uuid,
	text text,
	priority text,
	completed bool,
	time_created timestamp,
	time_updated timestamp,

	PRIMARY KEY (id)
);

CREATE INDEX priority_idx ON todos (priority);
CREATE INDEX completed_idx ON todos (completed);
