package types

type MGetRsp struct {
	Uids    string `json:"uids" path:"uids"`
	OwnerId string `json:"owner_id" path:"owner_id"`
}

type MGetReply struct {
	Data []Profile `json:"data"`
}

type GetRsp struct {
	Uid string `json:"uid" path:"uid"`
	Id  string `json:"id,omitempty" path:"id"`
}

type GetReply struct {
	Data Profile `json:"data"`
}

type Profile struct {
	Uid      int64  `json:"uid"`
	Nick     string `json:"nick"`     // 昵称
	Gender   int8   `json:"gender"`   // 性别：1男2女
	Portrait string `json:"portrait"` // 头像
	Age      int8   `json:"age"`      // 年纪
}
