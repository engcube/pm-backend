package PrivateMessageModel

import (
	"fmt"
	"strconv"
	"time"

	"pm-backend/public"

	"golang.org/x/crypto/bcrypt"
)

const (
	DIRECTION_SENT     = "sent"
	DIRECTION_RECEIVED = "received"
)

// User 用户信息
type User struct {
	UserID     int
	Email      string
	Username   string
	Password   string
	InsertTime int64
	UpdateTime int64
	IsDeleted  bool
	SessionID  string
}

// Get 获取用户信息
func (u *User) Get() error {
	if u.UserID == 0 {
		return fmt.Errorf("No UserID provided")
	}
	rows, err := PrivateMessageBackendPublic.Select(SQL_GET_USER, u.UserID)
	if err != nil {
		return err
	}
	if len(rows) != 1 {
		return fmt.Errorf("No User existed")
	}
	res := rows[0]
	userid, _ := strconv.ParseInt(string(res[0]), 10, 64)
	insertime, _ := strconv.ParseInt(string(res[4]), 10, 64)
	updatetime, _ := strconv.ParseInt(string(res[5]), 10, 64)
	u.Email = res[1]
	u.UserID = int(userid)
	u.Username = res[2]
	u.InsertTime = int64(insertime)
	u.UpdateTime = int64(updatetime)
	return nil
}

// Register 用户注册
func (u *User) Register() error {
	if u.Email == "" {
		return fmt.Errorf("No Email provided")
	}
	bExist, err := u.GetUserByEmail()
	if err != nil {
		return err
	}
	if bExist {
		return fmt.Errorf("User register with same email has already existed")
	}
	if u.Password == "" || len(u.Password) < 6 {
		return fmt.Errorf("Password should not be empty or less than 6 bytes")
	}
	if u.Username == "" {
		return fmt.Errorf("Empty username")
	}
	u.IsDeleted = false
	if u.InsertTime == 0 {
		u.InsertTime = time.Now().Unix()
	}
	if u.UpdateTime == 0 {
		u.UpdateTime = time.Now().Unix()
	}
	err = u.EncodePassword()
	if err != nil {
		return err
	}
	userid, err := PrivateMessageBackendPublic.Insert(SQL_NEW_USER, u.Email, u.Username, u.Password, u.InsertTime, u.UpdateTime)
	if err != nil {
		return err
	}
	u.UserID = int(userid)
	u.Password = ""
	return nil
}

// Update 更新用户信息（只能更新用户名）
func (u *User) Update() error {
	if u.Username == "" {
		return fmt.Errorf("No Username provided")
	}
	cnt, err := PrivateMessageBackendPublic.Update(SQL_UPDATE_USERNAME, u.Username, time.Now().Unix(), u.UserID)
	if cnt == 0 {
		return fmt.Errorf("no row updated")
	}
	return err
}

// Delete 删除指定用户
func (u *User) Delete() error {
	if u.UserID == 0 {
		return fmt.Errorf("No UserID provided")
	}
	cnt, err := PrivateMessageBackendPublic.Update(SQL_DELETE_USER, time.Now().Unix(), u.UserID)
	if cnt == 0 {
		return fmt.Errorf("no row updated")
	}
	u.IsDeleted = true
	return err
}

// GetUserByEmail 根据邮箱获取用户信息
func (u *User) GetUserByEmail() (bool, error) {
	rows, err := PrivateMessageBackendPublic.Select(SQL_GET_USER_BY_EMAIL, u.Email)
	if err != nil {
		return false, err
	}
	if len(rows) == 0 {
		return false, nil
	}
	if len(rows) != 1 {
		return false, fmt.Errorf("Multiple User with same email")
	}
	res := rows[0]
	userid, _ := strconv.ParseInt(string(res[0]), 10, 64)
	insertime, _ := strconv.ParseInt(string(res[4]), 10, 64)
	updatetime, _ := strconv.ParseInt(string(res[5]), 10, 64)
	u.Password = string(res[3])
	u.UserID = int(userid)
	u.Username = res[2]
	u.InsertTime = int64(insertime)
	u.UpdateTime = int64(updatetime)
	return true, nil
}

// Validate 验证用户名密码
func (u *User) Validate() (bool, error) {
	if u.Password == "" {
		return false, fmt.Errorf("No Password Provided")
	}
	if u.Email == "" {
		return false, fmt.Errorf("No Email Provided")
	}
	passwd := u.Password
	bExisted, err := u.GetUserByEmail()
	if err != nil {
		return false, err
	}
	if !bExisted {
		return false, fmt.Errorf("User Not Existed")
	}
	if u.ValidatePassword(passwd) {
		u.Password = ""
		return true, nil
	}
	return false, fmt.Errorf("Wrong Password")
}

// EncodePassword 加密密码
func (u *User) EncodePassword() error {
	if u.Password == "" {
		return fmt.Errorf("No Password Provided")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	encodePW := string(hash) // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即可
	u.Password = encodePW
	return nil
}

// ValidatePassword 验证密码
func (u *User) ValidatePassword(passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwd))
	if err != nil {
		return false
	}
	return true
}

// GetAllFriends 获取所有联系人信息
func (u *User) GetAllFriends() ([]Friend, error) {
	return u.GetFriend([]int{})
}

// GetFriend 获取联系人信息
func (u *User) GetFriend(userids []int) ([]Friend, error) {
	rows, err := PrivateMessageBackendPublic.Select(SQL_GET_FRIENDS, u.UserID)
	if err != nil {
		return nil, err
	}
	tmpFriends := make(map[int]Friend)
	for _, row := range rows {
		friend := Friend{}
		fid, _ := strconv.ParseInt(string(row[0]), 10, 64)
		friend.FriendID = int(fid)
		fuid, _ := strconv.ParseInt(string(row[1]), 10, 64)
		friend.FriendUserID = int(fuid)
		friend.Nickname = string(row[2])
		tmpFriends[friend.FriendUserID] = friend
	}
	if len(userids) == 0 {
		friends := make([]Friend, 0)
		for _, friend := range tmpFriends {
			friends = append(friends, friend)
		}
		return friends, nil
	}
	friends := make([]Friend, 0)
	for _, user := range userids {
		if f, ok := tmpFriends[user]; ok {
			friends = append(friends, f)
		}
	}
	return friends, nil
}

// AddFriend 添加联系人
func (u *User) AddFriend(friend *Friend) error {
	if u.UserID == 0 {
		return fmt.Errorf("userid not provided")
	}
	if friend.Email == "" {
		return fmt.Errorf("friend email not provided")
	}
	friendUser := User{Email: friend.Email}
	bExist, err := friendUser.GetUserByEmail()
	if err != nil {
		return err
	}
	if !bExist {
		return fmt.Errorf("friend not exist")
	}
	if friendUser.UserID == u.UserID {
		return fmt.Errorf("can not add self as friend")
	}
	rows, err := PrivateMessageBackendPublic.Select(SQL_GET_FRIEND, u.UserID, friendUser.UserID)
	if err != nil {
		return err
	}
	if len(rows) > 0 {
		return fmt.Errorf("already been friend")
	}
	fid, err := PrivateMessageBackendPublic.Insert(SQL_ADD_FRIEND, u.UserID, friendUser.UserID, friendUser.Username, time.Now().Unix())
	friend.FriendID = int(fid)
	friend.FriendUserID = friendUser.UserID
	friend.Nickname = friendUser.Username
	friend.IsDeleted = false
	friend.InsertTime = time.Now().Unix()
	friend.UnreadCount = 0
	return err
}

// DeleteFriend 删除联系人
func (u *User) DeleteFriend(friend *Friend) error {
	if u.UserID == 0 {
		return fmt.Errorf("userid not provided")
	}
	if friend.FriendID == 0 {
		return fmt.Errorf("friendid not provided")
	}
	cnt, err := PrivateMessageBackendPublic.Update(SQL_DELETE_FRIEND, time.Now().Unix(), friend.FriendID)
	if err != nil {
		return err
	}
	if cnt != 1 {
		return fmt.Errorf("no rows affected")
	}
	friend.IsDeleted = true
	friend.UpdateTime = time.Now().Unix()
	return nil
}

// GetMessages 获取所有联系人信息数
func (u *User) GetMessages(userids []int) ([]Friend, error) {
	sent, err := u.GetMessagesByDirection([]int{}, DIRECTION_SENT)
	if err != nil {
		return nil, err
	}
	recieved, err := u.GetMessagesByDirection([]int{}, DIRECTION_RECEIVED)
	if err != nil {
		return nil, err
	}
	friends := make(map[int]*Friend)
	for _, f := range sent {
		if friend, ok := friends[f.FriendUserID]; ok {
			friend.TotalCount += f.TotalCount
			friend.UnreadCount += f.UnreadCount
			friend.SentMsgs = append(friend.SentMsgs, f.SentMsgs...)
		} else {
			friends[f.FriendUserID] = &Friend{FriendUserID: f.FriendUserID, TotalCount: f.TotalCount, UnreadCount: f.UnreadCount, SentMsgs: f.SentMsgs}
		}
	}
	for _, f := range recieved {
		if friend, ok := friends[f.FriendUserID]; ok {
			friend.TotalCount += f.TotalCount
			friend.UnreadCount += f.UnreadCount
			friend.RecieveMsgs = append(friend.RecieveMsgs, f.RecieveMsgs...)
		} else {
			friends[f.FriendUserID] = &Friend{FriendUserID: f.FriendUserID, TotalCount: f.TotalCount, UnreadCount: f.UnreadCount, RecieveMsgs: f.RecieveMsgs}
		}
	}
	if len(userids) == 0 {
		friendlists := make([]Friend, 0)
		for _, friend := range friends {
			friendlists = append(friendlists, *friend)
		}
		return friendlists, nil
	}
	friendlists := make([]Friend, 0)
	for _, friend := range friends {
		bFound := false
		for _, id := range userids {
			if id == friend.FriendUserID {
				bFound = true
				break
			}
		}
		if bFound {
			friendlists = append(friendlists, *friend)
		}
	}
	return friendlists, nil
}

// GetMessagesByDirection 获取联系人信息数
func (u *User) GetMessagesByDirection(userids []int, direction string) ([]Friend, error) {
	sql := SQL_GET_MESSAGE_RECIEVED
	if direction == DIRECTION_SENT {
		sql = SQL_GET_MESSAGE_SENT
	}
	// receiver to是自己
	rows, err := PrivateMessageBackendPublic.Select(sql, u.UserID)
	if err != nil {
		return nil, err
	}
	tmpFriends := make(map[int]*Friend)
	for _, row := range rows {
		message := Message{}
		mid, _ := strconv.ParseInt(string(row[0]), 10, 64)
		message.MessageID = int(mid)
		uid, _ := strconv.ParseInt(string(row[1]), 10, 64)
		message.Sender = int(uid)
		touid, _ := strconv.ParseInt(string(row[2]), 10, 64)
		message.Reciever = int(touid)
		message.Content = string(row[3])
		isviewed, _ := strconv.ParseInt(string(row[4]), 10, 32)
		if int(isviewed) == 1 {
			message.IsViewed = true
		} else {
			message.IsViewed = false
		}
		inserttime, _ := strconv.ParseInt(string(row[5]), 10, 64)
		message.InsertTime = int64(inserttime)
		updatetime, _ := strconv.ParseInt(string(row[6]), 10, 64)
		message.UpdateTime = int64(updatetime)
		message.IsDeleted = false
		if direction == DIRECTION_SENT {
			if friend, ok := tmpFriends[message.Reciever]; ok {
				friend.SentMsgs = append(friend.SentMsgs, message)
			} else {
				tmpFriends[message.Reciever] = &Friend{FriendUserID: message.Reciever}
				tmpFriends[message.Reciever].SentMsgs = append(tmpFriends[message.Reciever].SentMsgs, message)
			}
		} else {
			if friend, ok := tmpFriends[message.Sender]; ok {
				friend.RecieveMsgs = append(friend.RecieveMsgs, message)
				if !message.IsViewed {
					friend.UnreadCount++
				}
			} else {
				tmpFriends[message.Sender] = &Friend{FriendUserID: message.Sender}
				tmpFriends[message.Sender].RecieveMsgs = append(tmpFriends[message.Sender].RecieveMsgs, message)
				if !message.IsViewed {
					tmpFriends[message.Sender].UnreadCount++
				}
			}
		}
	}
	if len(userids) == 0 {
		friends := make([]Friend, 0)
		for _, friend := range tmpFriends {
			friend.TotalCount = len(friend.RecieveMsgs) + len(friend.SentMsgs)
			friends = append(friends, *friend)
		}
		return friends, nil
	}
	friends := make([]Friend, 0)
	for _, user := range userids {
		if f, ok := tmpFriends[user]; ok {
			f.TotalCount = len(f.RecieveMsgs) + len(f.SentMsgs)
			friends = append(friends, *f)
		}
	}
	return friends, nil
}

// SendMessage 发送新消息
func (u *User) SendMessage(message *Message) error {
	if message.RecieverEmail == "" {
		return fmt.Errorf("no reciever email provided")
	}
	if message.Content == "" {
		return fmt.Errorf("no content provided")
	}
	message.Sender = u.UserID
	friend := User{Email: message.RecieverEmail}
	bExist, err := friend.GetUserByEmail()
	if err != nil {
		return err
	}
	if !bExist {
		return fmt.Errorf("No user existed")
	}

	isFriend, err := u.IsFriend(&friend)
	if err != nil {
		return err
	}
	fmt.Printf("is friend: %t\n", isFriend)
	if !isFriend {
		return fmt.Errorf("Please be friend first")
	}

	isFriend, err = friend.IsFriend(u)
	if err != nil {
		return err
	}
	if !isFriend {
		err = friend.addFriend2(u)
		if err != nil {
			return err
		}
	}

	message.Reciever = friend.UserID
	err = message.New()
	return err
}

func (u *User) addFriend2(to *User) error {
	if u.UserID == 0 || to.UserID == 0 {
		return fmt.Errorf("userid not provided")
	}
	if to.UserID == u.UserID {
		return fmt.Errorf("can not add self as friend")
	}
	rows, err := PrivateMessageBackendPublic.Select(SQL_GET_FRIEND, u.UserID, to.UserID)
	if err != nil {
		return err
	}
	if len(rows) > 0 {
		return nil
	}
	_, err = PrivateMessageBackendPublic.Insert(SQL_ADD_FRIEND, u.UserID, to.UserID, to.Username, time.Now().Unix())
	return err
}

// IsFriend u是否是o的联系人
func (u *User) IsFriend(o *User) (bool, error) {
	if u.UserID == 0 || o.UserID == 0 {
		return false, fmt.Errorf("No UserID provided")
	}
	rows, err := PrivateMessageBackendPublic.Select(SQL_GET_FRIENDSHIP, u.UserID, o.UserID)
	if err != nil {
		return false, err
	}
	if len(rows) == 0 {
		return false, nil
	}
	return true, nil
}

// ReadMessage 阅读message
func (u *User) ReadMessage(message *Message) error {
	err := message.Get()
	if err != nil {
		return err
	}
	if message.Reciever != u.UserID {
		return fmt.Errorf("permission denied")
	}
	if message.IsViewed {
		return fmt.Errorf("Message has already been viewed")
	}
	err = message.Read()
	message.IsViewed = true
	message.UpdateTime = time.Now().Unix()
	return err
}

// DeleteMessage 删除message
func (u *User) DeleteMessage(message *Message) error {
	err := message.Get()
	if err != nil {
		return err
	}
	// 发送和接受者都可以删除消息
	if message.Reciever != u.UserID && message.Sender != u.UserID {
		return fmt.Errorf("permission denied")
	}
	err = message.Delete()
	message.IsDeleted = true
	message.UpdateTime = time.Now().Unix()
	return err
}
