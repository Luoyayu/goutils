package main

// update AccountSelected
func findAccountLastUsed() error {
	return db.QueryRow(`
			SELECT uid,nikeName,loginUserName,accessToken,expire,SESSDATA,sid,DedeUserID__ckMd5,lastUsedTimestamp,blocked FROM Account
			WHERE lastUsedTimestamp IS (SELECT max(lastUsedTimestamp) FROM Account);`).Scan(
		&AccountSelected.Uid,
		&AccountSelected.NikeName,
		&AccountSelected.LoginUserName,
		&AccountSelected.AccessToken,
		&AccountSelected.Expire,
		&AccountSelected.SESSDATA,
		&AccountSelected.Sid,
		&AccountSelected.DedeUserID__ckMd5,
		&AccountSelected.LastUsedTimestamp,
		&AccountSelected.Blocked,
	)
}

func findAllAccounts() {
	cnts := 0
	if _ = db.QueryRow(`SELECT COUNT(uid) FROM Account`).Scan(&cnts); cnts != 0 {
		AccountsInDB = make([]*Account, cnts)

		rows, err := db.Query(`SELECT uid,nikeName,loginUserName,accessToken,expire,SESSDATA,sid,DedeUserID__ckMd5,lastUsedTimestamp,blocked FROM Account;`)
		if err == nil {
			defer rows.Close()
			for i := 0; i < cnts; i++ {
				rows.Next()
				AccountsInDB[i] = &Account{}
				if err := rows.Scan(&AccountsInDB[i].Uid,
					&AccountsInDB[i].NikeName,
					&AccountsInDB[i].LoginUserName,
					&AccountsInDB[i].AccessToken,
					&AccountsInDB[i].Expire,
					&AccountsInDB[i].SESSDATA,
					&AccountsInDB[i].Sid,
					&AccountsInDB[i].DedeUserID__ckMd5,
					&AccountsInDB[i].LastUsedTimestamp,
					&AccountsInDB[i].Blocked,
				); err != nil {
					Logger.Error(err)
				} else {

				}
			}
		} else {
			Logger.Error(err)
		}
	}
}
