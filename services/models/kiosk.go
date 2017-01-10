package models

import (
    "time"
    "math/rand"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)


type Kiosk struct {  
    ID     int64 `json:"id" form:"-"`
    Name string  `json:"name" form:"name"`
    Status int  `json:"status" form:"status"`
    UserName string `json:"username" form:"username"`
    ConfirmationCode string `json:"confirmationcode" form:"confirmationcode"`
    Configuration     int64 `json:"configuration" form:"configuration"`
    Created time.Time `json:"created" form:"created"`
    Token string `json:"token" form:"token"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

//RemoveKiosk - remove kiosk by providing Kiosk object and username
func RemoveKiosk(k *Kiosk, usr string)(error){
    stmt, err := db.Prepare("DELETE FROM Kiosk where ID=? and UserName=?")
    if err != nil {
        return err
    }

    _, err2 := stmt.Exec(k.ID,usr)
    if err2 != nil {
        return err2
    }

    return nil
}

//UpdateKiosk - Update Kiosk by providing Kiosk object
func UpdateKiosk(k *Kiosk, usr string)(*Kiosk, error){   
    stmt, err := db.Prepare("UPDATE Kiosk SET Status=?,Name=?,Configuration=? where ID=? and UserName=?")
    if err != nil {
        return nil, err
    }

    _, err2 := stmt.Exec(k.Status,k.Name,k.Configuration,k.ID,usr)
    if err2 != nil {
        return nil, err2
    }
    return k, nil
}



//GetKiosks - return list of Kiosks which belong to User
func GetKiosks(usr string)([]*Kiosk, error){
    rows, err := db.Query("SELECT ID,Name,UserName,Status FROM Kiosk where UserName=?",usr)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    bks := make([]*Kiosk, 0)
    for rows.Next() {
        bk := new(Kiosk)
        err := rows.Scan(&bk.ID,&bk.Name,&bk.UserName,&bk.Status)
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

func GetTokenByKioskID(k *Kiosk)(string,error){
    rows, err := db.Query("SELECT Token FROM Kiosk where ID=? AND UserName=?",k.ID, k.UserName)
    if err != nil {
        return "", err
    }
    defer rows.Close()

    for rows.Next() {
        bk := new(Kiosk)
        err := rows.Scan(&bk.Token)
        if err != nil {
            return "", err
        }
        return bk.Token, err
    }
    return  "", err  
}

//GetKioskByID - return Kiosk by ID and USER it belongs to
func GetKioskByID(id string,usr string)(*Kiosk, error){
    rows, err := db.Query("SELECT ID,Name,UserName,Status FROM Kiosk where ID=? AND UserName=?",id,usr)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        bk := new(Kiosk)
        err := rows.Scan(&bk.ID,&bk.Name,&bk.UserName,&bk.Status)
        if err != nil {
            return nil, err
        }
        return bk, err
    }
    return  nil, err
}


//SignUpKiosk Add new record to Kiosk table
func SignUpKiosk(username string)(*Kiosk,string,error){

     u:=new(User)
     u.UserName=username;
     
     u, err:=CheckUserExists(u)
     if err!=nil{
         return nil,"",err
     }

    k,confimrCode,err:=addKiosk(u)
    if err!= nil{
        return nil,"",err
    }

        return k,confimrCode,nil
}

//LinkKioskToAccount Adding linkage betweem Kiosk and User Account. Will return the token for Kiosk
func LinkKioskToAccount(kiosk *Kiosk)(*Kiosk,error){
    k,err:=checkKioskExists(kiosk)
    if err!=nil{
        return nil,err
    }

    if k!=nil {
        k.Token=kiosk.Token
        k.Name=kiosk.Name
        return linkKiosk(k)
    }
    
    return nil,nil
}


func linkKiosk(k *Kiosk)(*Kiosk,error){
    //update
    stmt, err := db.Prepare("UPDATE Kiosk SET Status=?,Name=?,Token=? where ID=?")
    if err != nil {
        return nil, err
    }

    res, err := stmt.Exec(2,k.Name,k.Token,k.ID)
    if err != nil {
        return nil, err
    }
    res.LastInsertId()
    k.Status=2

    return k, nil
}




func checkKioskExists(k *Kiosk)(*Kiosk,error){
    rows, err := db.Query("SELECT ID,ConfirmationCode FROM Kiosk where ID=? AND ConfirmationCode=?",k.ID,k.ConfirmationCode)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        bk := new(Kiosk)
        err := rows.Scan(&bk.ID,&bk.ConfirmationCode)
        if err != nil {
            return nil, err
        }
        return bk, err
    }
    return  nil, err
}




func addKiosk(u *User)(*Kiosk,string,error){
    // insert
    stmt, err := db.Prepare("INSERT Kiosk SET UserName=?,Status=?,ConfirmationCode=?,Created=?")
    if err != nil {
        return nil,"", err
    }
    confirmCode:=randStringBytesMaskImpr(6)
    created:=time.Now()
    res, err := stmt.Exec(u.UserName,1,confirmCode,created)
    if err != nil {
        return nil,"", err
    }

    k:=new(Kiosk)    

    id, err := res.LastInsertId()
    if err != nil {
        return nil,"", err
    }
    
    k.ID=id;
    k.UserName=u.UserName
    k.Status=1
    k.Created=created

    return k,confirmCode,nil 
}



func randStringBytesMaskImpr(n int) string {
    b := make([]byte, n)
    // A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
    for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = rand.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return string(b)
}




