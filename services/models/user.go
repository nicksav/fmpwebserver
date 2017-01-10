package models

import (
    "time"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "strconv"
    "strings"
    "io"
    "net/url"
    "golang.org/x/crypto/bcrypt"
)

const (
	WebPortalMethod           = 1 
    DeviceMethod              = 2
    AllMethod                 = 3
)

type User struct {  
    ID     int64 `json:"id" form:"-"`
    FirstName string  `json:"firstName" form:"firstName"`
    LastName string  `json:"lastName" form:"lastName"`
    Status int  `json:"status" form:"status"`
    UserName string `json:"username" form:"username"`
    Password string `json:"password" form:"password"`
    Customer int `json:"customer" form:"customer"`
    Created time.Time `json:"created" form:"created"`
}


func AllUsers() ([]*User, error) {
    rows, err := db.Query("SELECT * FROM User")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    bks := make([]*User, 0)
    for rows.Next() {
        bk := new(User)
        err := rows.Scan(&bk.ID,&bk.FirstName,&bk.LastName,&bk.Status, &bk.UserName, &bk.Password, &bk.Customer, &bk.Created)
        if err != nil {
            return nil, err
        }
        bks = append(bks, bk)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return bks, nil
}


func CheckUserExists(usr *User)(*User,error){
    rows, err := db.Query("SELECT UserName,Status FROM User where UserName=?",usr.UserName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        bk := new(User)
        err := rows.Scan(&bk.UserName,&bk.Status)
        if err != nil {
            return nil, err
        }
        return bk, err
    }
    return  nil, err
}

func CheckUserExistsWithPassword(usr *User)(*User,error){
    rows, err := db.Query("SELECT UserName,Status,Password FROM User where UserName=?",usr.UserName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        bk := new(User)
        err := rows.Scan(&bk.UserName,&bk.Status,&bk.Password)
        if err != nil {
            return nil, err
        }
        return bk, err
    }
    return  nil, err
}


func RequestConfirmationCode(usr *User)(string,error){
    code,err:=requestConfirmCode(usr)
    if err != nil {
        return "",err
    }

    return code,nil
}


func SignUpUser(usr *User)(*User,error){

     // insert
    stmt, err := db.Prepare("INSERT User SET FirstName=?,LastName=?,Status=?,UserName=?,Password=?,Customer=?,Created=?")
    if err != nil {
        return nil, err
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(usr.Password), 10)
    res, err := stmt.Exec(usr.FirstName,usr.LastName,1,usr.UserName,string(hashedPassword),usr.Customer,time.Now())
    if err != nil {
        return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }
    
    usr.Password="";
    usr.ID=id;


    return usr,nil
}


func requestConfirmCode(usr *User)(string,error){
   timestamp:= time.Now().Unix()
   confirmcode:=[]byte(strconv.FormatInt(timestamp,10)+";"+usr.UserName)

   key := []byte("1111111111111111")
   resultCode, err:=encrypt(key,confirmcode)
   if err != nil {
        return "", err
    }

    return url.QueryEscape(base64.URLEncoding.EncodeToString(resultCode)),nil

}



func VlidateConfirmCode(confirmCode string) (*User,error){
    data, err := base64.URLEncoding.DecodeString(confirmCode)

    key := []byte("1111111111111111")
    result, err := decrypt(key, data)
    if err != nil {
        return nil, err
    }

    s := strings.Split(string(result), ";")
    username := s[1]

    newUser:=  new(User)
    newUser.UserName=username;
    user, err:=CheckUserExists(newUser)
    if err != nil {
        return nil, err
    }

    if(user!=nil && user.Status==1){
        stmt, err := db.Prepare("UPDATE User SET Status=? where UserName=?")
        if err != nil {
            return nil, err
        }

        res, err := stmt.Exec(2,user.UserName)
        
        if err != nil {
            return nil, err
        }

        id, err := res.LastInsertId()
        if err != nil {
            return nil, err
        }


        user.Status=2;
        user.ID=id;
        return user,nil
    }else{
        return user,nil
    }   

}



func encrypt(key, text []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    b := base64.StdEncoding.EncodeToString(text)
    ciphertext := make([]byte, aes.BlockSize+len(b))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
    return ciphertext, nil
}

func decrypt(key, text []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    if len(text) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }
    iv := text[:aes.BlockSize]
    text = text[aes.BlockSize:]
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(text, text)
    data, err := base64.StdEncoding.DecodeString(string(text))
    if err != nil {
        return nil, err
    }
    return data, nil
}







//GetUserByUsername - Get User by UserName. Should return User or nill
func GetUserByUsername(UserName string)(*User,error){
    rows, err := db.Query("SELECT * FROM User where UserName=?",UserName)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    bk := new(User)
    for rows.Next() {
        err := rows.Scan(&bk.ID,&bk.FirstName,&bk.LastName,&bk.Status, &bk.UserName, &bk.Password, &bk.Customer, &bk.Created)
        if err != nil {
            return nil, err
        }

    }
    return bk,nil;
}