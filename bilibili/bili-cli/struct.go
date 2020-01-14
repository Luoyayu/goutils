package main

import biliAPI "github.com/luoyayu/goutils/bilibili"

type User struct {
	Uid      int64
	Cid      int64
	NikeName string
}

type Account struct {
	Uid               int64
	NikeName          string
	LoginUserName     string
	AccessToken       string
	Expire            int64
	SESSDATA          string
	Sid               string
	DedeUserID__ckMd5 string
	LastUsedTimestamp int64
	Blocked           int64
	Info              *biliAPI.UserInfoDataStruct
}

type Following struct {
	Uid      int64 // local user id
	Fid      int64
	NikeName string // Fid
	Blocked  int64
}

type Live struct {
	Uid      int64
	Fid      int64
	NikeName string
	Cid      int64
	Title    string
	State    int
	Face     string
	Blocked  int64
}
