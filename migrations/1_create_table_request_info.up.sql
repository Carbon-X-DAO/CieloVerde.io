CREATE TABLE IF NOT EXISTS request_info(
	id SERIAL,
	acceptlanguage TEXT,
	cookie TEXT,
	useragent TEXT,
	cfconnectingip TEXT,
	xforwardedfor TEXT,
	cfray TEXT,
	cfipcountry TEXT,
	cfvisitor TEXT,
	url_value TEXT,
	ctime TIMESTAMP WITH TIME ZONE
);
