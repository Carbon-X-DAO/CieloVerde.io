CREATE TABLE IF NOT EXISTS qr_incoming_headers(
	id SERIAL,
	acceptlanguage TEXT,
	cookie TEXT,
	useragent TEXT,
	cfconnectingip TEXT,
	xforwardedfor TEXT,
	cfray TEXT,
	cfipcountry TEXT,
	cfvisitor TEXT,
	ctime TIMESTAMP WITH TIME ZONE
);
