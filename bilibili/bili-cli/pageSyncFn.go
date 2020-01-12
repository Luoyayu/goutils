package main

import (
	biliAPI "github.com/luoyayu/goutils/bilibili"
	"log"
	"math"
	"sync"
)

func pageSyncCmdSyncFollowing(hintCoverFollowing bool) []*Following {
	stat, err := biliAPI.GetRelationStat(AccountSelected.Uid)
	if err != nil {
		log.Fatalln(err)
	}
	tot := stat.Data.Following

	var followings = make([]*Following, tot)
	followingsMutex := &sync.Mutex{}

	wg := &sync.WaitGroup{}
	times := int(math.Ceil(float64(tot) / 50.0))

	wg.Add(times)
	for i := 0; i < times; i++ {
		go func(idx int) {
			ret, err := biliAPI.GetRelationFollowings(AccountSelected.Uid, idx+1, "")
			if err != nil {
				Logger.Error(err)
			} else {

				followingsMutex.Lock()
				for j, user := range ret.Data.List {
					_, _ = db.Exec(`REPLACE INTO User(uid, nikeName) VALUES (?,?)`, user.Mid, user.Uname)
					//log.Println(idx*50+j+1, user.Mid, user.Uname)
					followings[idx*50+j] = &Following{
						Uid:      AccountSelected.Uid,
						Fid:      user.Mid,
						NikeName: user.Uname,
						Blocked:  0,
					}
				}
				followingsMutex.Unlock()
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	if hintCoverFollowing {
		/*isCoverPage := promptui.Prompt{
			Label:     "cover local following data",
			IsConfirm: true,
		}
		yn, _ := isCoverPage.Run()*/

		yn := promptConfirm("cover local following data")

		if yn == true {
			FollowingInDB = followings
			for _, ff := range followings {
				_, _ = db.Exec(`REPLACE INTO Following(uid, fid) VALUES (?,?);`, ff.Uid, ff.Fid)
			}
		} else if yn == false {
			for _, ff := range followings {
				r, err := db.Exec("SELECT COUNT(*) FROM Following WHERE uid=? AND fid=?;", ff.Uid, ff.Fid)
				if err != nil {
					log.Fatal(err)
				}

				if n, _ := r.RowsAffected(); n == 0 {
					_, _ = db.Exec(`INSERT INTO Following(uid, fid) VALUES (?,?);`, ff.Uid, ff.Fid)
				}
			}
		}
	}

	/*for i, ff := range followings {
		fmt.Println(i+1, ff.Uid, "->", ff.Fid, ff.NikeName)
	}*/
	return followings
}

func pageSyncCmdSyncLive(hintCoverLive, hintCoverFollowing, getAll bool) []*Live {
	followings := pageSyncCmdSyncFollowing(hintCoverFollowing)

	live := make([]*Live, 0, len(followings))
	times := int(math.Ceil(float64(len(followings)) / 20.0))

	liveMutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(times)

	for i := 0; i < times; i++ {
		go func(idx int) {
			//log.Println("No.", idx)
			liveFollowing, err := biliAPI.GetLiveMyFollowing(AccountSelected.Uid, AccountSelected.SESSDATA, idx+1, 20)
			if err != nil {
				Logger.Error(err)
			}
			liveMutex.Lock()
			for _, user := range liveFollowing.Data.List {
				//log.Println(user.UName)
				if user.LiveStatus == 0 && getAll == false {
					break
				}
				_, _ = db.Exec(`REPLACE INTO User(uid, cid, nikeName) VALUES (?,?,?)`, user.Uid, user.RoomId, user.UName)
				live = append(live, &Live{
					Uid:      AccountSelected.Uid,
					Fid:      user.Uid,
					NikeName: user.UName,
					Cid:      user.RoomId,
					Title:    user.Title,
					State:    user.LiveStatus,
					Face:     user.Face,
				})

			}
			liveMutex.Unlock()

			wg.Done()
		}(i)
	}
	wg.Wait()

	if hintCoverLive {
		/*	isCoverPage := promptui.Prompt{
				Label:     "cover local live data",
				IsConfirm: true,
			}

			yn, _ := isCoverPage.Run()
		*/

		yn := promptConfirm("cover local live data")

		if yn == true {
			LiveInDB = live
			for _, ff := range live {
				_, _ = db.Exec(`
				REPLACE INTO Live(uid,fid,cid,title,state,face) 
				VALUES (?,?,?,?,?,?)`, ff.Uid, ff.Fid, ff.Cid, ff.Title, ff.State, ff.Face)
			}
		} else if yn == false {
			for _, ff := range live {
				r, err := db.Exec("SELECT COUNT(*) FROM Live WHERE fid=? AND uid=?", ff.Fid, ff.Uid)
				if err != nil {
					log.Fatal(err)
				}
				if n, _ := r.RowsAffected(); n == 0 {
					_, _ = db.Exec(`
					INSERT INTO Live(uid,fid,cid,title,state,face) 
					VALUES (?,?,?,?,?,?);`, ff.Uid, ff.Fid, ff.Cid, ff.Title, ff.State, ff.Face)
				}
			}
		}
	}

	/*for i, ff := range live {
		fmt.Println(i+1, ff.Fid, ff.nikeName, ff.State == 1, ff.Title)
	}*/
	return live
}

func pageSyncCmdSyncBoth() {
	pageSyncCmdSyncLive(true, true, true)
}
