package main

const ()

func initDB() error {
	_, err := db.Exec(`
		CREATE TABLE Preference ( 
			singleMode INTEGER DEFAULT 0
		);
		CREATE TABLE User (
			uid      INTEGER UNIQUE PRIMARY KEY,
			cid		 INTEGER DEFAULT 0,
			nikeName TEXT DEFAULT ''
		);
		CREATE TABLE Account (
			uid					INTEGER UNIQUE PRIMARY KEY,
			nikeName 			TEXT DEFAULT '',
			loginUserName		TEXT DEFAULT '',
			accessToken			TEXT DEFAULT '',
			expire				INTEGER DEFAULT 0,
			SESSDATA 			TEXT,
			sid 				TEXT,
			DedeUserID__ckMd5 	TEXT, 
			lastUsedTimestamp 	INTEGER DEFAULT 0,
			blocked				INTEGER DEFAULT 0
		);
		CREATE TABLE Following (
			uid      INTEGER NOT NULL ,
			fid 	 INTEGER NOT NULL ,
			blocked	 INTEGER DEFAULT 0,
		    CONSTRAINT Following_pk
			PRIMARY KEY (uid, fid)
		);
		CREATE TABLE Live (
			uid      INTEGER NOT NULL ,
			fid 	 INTEGER NOT NULL ,
			cid      INTEGER DEFAULT 0, 
			title 	 TEXT,
			state 	 INTEGER,
			face 	 TEXT,
			blocked	 INTEGER  DEFAULT 0,
		    CONSTRAINT Live_pk
			PRIMARY KEY (uid, fid)
		);`,
	)
	return err
}
