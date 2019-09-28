package models

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/mdg-iitr/Codephile/models/db"
	"github.com/mdg-iitr/Codephile/models/submission"
	"github.com/mdg-iitr/Codephile/scripts"
	"log"
	"github.com/mdg-iitr/Codephile/models/profile"
)

type User struct {
	ID          bson.ObjectId          `bson:"_id" json:"id" schema:"-"`
	Username    string                 `bson:"username" json:"username" schema:"username"`
	Password    string                 `bson:"password" json:"-" schema:"password"`
	Handle      Handle                 `bson:"handle" json:"handle" schema:"handle"`
	Submissions submission.Submissions `bson:"submission" json:"-" schema:"-"`
}
type Handle struct {
	Codeforces  string `bson:"codeforces" json:"codeforces" schema:"codeforces"`
	Codechef    string `bson:"codechef" json:"codechef" schema:"codechef"`
	Spoj        string `bson:"spoj" json:"spoj" schema:"spoj"`
	Hackerrank  string `bson:"hackerrank" json:"hackerrank" schema:"hackerrank"`
	Hackerearth string `bson:"hackerearth" json:"hackerearth" schema:"hackerearth"`
}

func AddUser(u User) (string, error) {
	u.ID = bson.NewObjectId()
	collection := db.NewCollectionSession("coduser")
	defer collection.Close()
	err := collection.Session.Insert(u)
	if err != nil {
		panic(err)
	}
	return u.ID.Hex(), nil
}

func GetUser(uid bson.ObjectId) (*User, error) {
	var user User
	collection := db.NewCollectionSession("coduser")
	defer collection.Close()
	err := collection.Session.FindId(uid).Select(bson.M{"_id": 1, "username": 1, "handle": 1}).One(&user)
	//fmt.Println(err.Error())
	if err != nil {
		return nil, errors.New("user not exists")
	}
	return &user, nil
}

func GetAllUsers() []User {
	var users []User
	collection := db.NewCollectionSession("coduser")
	defer collection.Close()
	err := collection.Session.Find(nil).Select(bson.M{"_id": 1, "username": 1, "handle": 1}).All(&users)
	if err != nil {
		panic(err)
	}
	return users
}

func UpdateUser(uid bson.ObjectId, uu *User) (a *User, err error) {
	if u, err := GetUser(uid); err == nil {
		if uu.Username != "" {
			u.Username = uu.Username
		}
		if uu.Handle.Codechef != "" {
			u.Handle.Codechef = uu.Handle.Codechef
		}
		if uu.Handle.Codeforces != "" {
			u.Handle.Codeforces = uu.Handle.Codeforces
		}
		if uu.Handle.Hackerearth != "" {
			u.Handle.Hackerearth = uu.Handle.Hackerearth
		}
		if uu.Handle.Hackerrank != "" {
			u.Handle.Hackerrank = uu.Handle.Hackerrank
		}
		if uu.Handle.Hackerearth != "" {
			u.Handle.Hackerearth = uu.Handle.Hackerearth
		}
		collection := db.NewCollectionSession("coduser")
		_, err := collection.Session.UpsertId(uid, &u)
		return u, err
	}
	return nil, errors.New("User Not Exist")
}
func AutheticateUser(username string, password string) (*User, bool) {
	var user User
	collection := db.NewCollectionSession("coduser")
	defer collection.Close()
	err := collection.Session.Find(bson.M{"username": username}).One(&user)
	//fmt.Println(err.Error())
	if err != nil {
		return nil, false
	}
	if user.Password == password {
		return &user, true
	}
	return nil, false
}

func AddSubmissions(user *User, site string) error {
	var sub submission.Submissions
	var handle string
	coll := db.NewCollectionSession("coduser")
	switch site {
	case "codechef":
		handle = user.Handle.Codechef
		if handle == "" {
			return errors.New("handle not available")
		}
		sub.Codechef = scripts.GetCodechefSubmissions(handle)
		err := coll.Session.UpdateId(user.ID, bson.M{"$set": bson.M{"submission.codechef": sub.Codechef}})
		if err != nil {
			log.Fatal(err.Error())
		}
		return nil
	case "codeforces":
		handle = user.Handle.Codeforces
		if handle == "" {
			return errors.New("handle not available")
		}
		sub.Codeforces = scripts.GetCodeforcesSubmissions(handle).Data
		err := coll.Session.UpdateId(user.ID, bson.M{"$set": bson.M{"submission.codeforces": sub.Codeforces}})
		if err != nil {
			log.Fatal(err.Error())
		}
		return nil
	case "spoj":
		handle = user.Handle.Spoj
		if handle == "" {
			return errors.New("handle not available")
		}
		sub.Spoj = scripts.GetSpojSubmissions(handle)
		err := coll.Session.UpdateId(user.ID, bson.M{"$set": bson.M{"submission.spoj": sub.Spoj}})
		if err != nil {
			log.Fatal(err.Error())
		}
		return nil
	case "hackerrank":
		handle = user.Handle.Hackerrank
		if handle == "" {
			return errors.New("handle not available")
		}
		sub.Hackerrank = scripts.GetHackerrankSubmissions(handle).Data
		err := coll.Session.UpdateId(user.ID, bson.M{"$set": bson.M{"submission.hackerrank": sub.Hackerrank}})
		if err != nil {
			log.Fatal(err.Error())
		}
		return nil
	}
	return nil
}

func GetSubmissions(ID bson.ObjectId) (*submission.Submissions, error) {
	coll := db.NewCollectionSession("coduser")
	var user User
	err := coll.Session.FindId(ID).Select(bson.M{"submission": 1}).One(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user.Submissions, nil
}

func AddorUpdateProfile(uid bson.ObjectId, site string, handle string) ( *User, error) {

	var UserProfile profile.ProfileInfo 
	//runs code to fetch the particular script's getProfile function
	switch site {
    case "codechef":
		UserProfile = scripts.GetCodechefProfileInfo(handle)
		break;
	case "codeforces":
		UserProfile = scripts.GetCodeforcesProfileInfo(handle)
		break;
	case "spoj":
		UserProfile = scripts.GetSpojProfileInfo(handle)
		break;
	case "hackerrank":
		UserProfile = scripts.GetHackerrankProfileInfo(handle)
		break;
	}  // add a default case for non-existent website
	//Profile fetched. Store in database 
	 user,err := GetUser(uid)
	 if err == nil {
		var ProfileTobeInserted profile.Profile
		ProfileTobeInserted.Website = site
		ProfileTobeInserted.Profileinfo = UserProfile
		// ProfileTobeInserted is all set to be put in the database
		collection := db.NewCollectionSession("coduser")
		defer collection.Close()
		// err2 := collection.Session.Update(bson.D{{"_id" , user.ID}},bson.D{{"$set" , ProfileTobeInserted}})
		NewNode := site + "Profile"
		SelectedUser := bson.D{{"_id", user.ID}}
		Update := bson.D{{"$set" , bson.D{ {NewNode , ProfileTobeInserted }}}}
		_, err2 := collection.Session.Upsert(SelectedUser,Update)
		//inserted into the document
		if err2 == nil{
			return user , nil
		} else {
			return nil, err2
		}
	 } else {
		 //handle the error (Invalid user)
		  return nil, err 
	 } 
}

// func Login(username, password string) bool {
// 	for _, u := range UserList {
// 		if u.Username == username && u.Password == password {
// 			return true
// 		}
// 	}
// 	return false
// }

// func DeleteUser(uid string) {
// 	delete(UserList, uid)
// }
