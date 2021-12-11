CREATE TABLE IF NOT EXISTS email_status(
	id SERIAL,
	email_address TEXT,
	gov_id TEXT,
	mailgun_msg TEXT,
	mailgun_id TEXT,
	error TEXT,
	ctime TIMESTAMP WITH TIME ZONE
);
