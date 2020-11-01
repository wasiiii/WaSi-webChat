package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/bmizerany/pq"
	"github.com/gorilla/websocket"
)

type user struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type binfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type userall struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Role     int    `json:"role"`
}

type uesrlogin struct {
	Acountnum string `json:"acountnum"`
	Password  string `json:"password"`
	Role      int    `json:"role"`
}

type username struct {
	Name string `json:"name"`
}

type userrole struct {
	Role int `json:"role"`
}

type token struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

var db *sql.DB //连接池对象
func initDB() (err error) {
	//连接数据库
	dsn := "user=postgres password=postgres dbname=wasipg sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	return
}

//ValidToken : 检验token有效性
func ValidToken(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	var t1 token
	json.Unmarshal([]byte(body), &t1)

	var ucheck userall
	ss := db.QueryRow("SELECT email, password FROM userinfo WHERE name=$1", t1.Name)
	err = ss.Scan(&ucheck.Email, &ucheck.Password)
	if err != nil {
		fmt.Println("err = ", err)
	}

	if t1.Name+ucheck.Email == t1.Token && ucheck.Password == t1.Password {
		str := token{Token: t1.Token}
		save, err1 := json.Marshal(&str)
		if err1 != nil {
			fmt.Println("err = ", err)
			return
		}
		fmt.Fprintln(w, string(save))
	}
}

//注册
func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	var u1 user
	json.Unmarshal([]byte(body), &u1)

	//查询
	var ucheck user
	ss := db.QueryRow("SELECT name FROM userinfo WHERE name=$1", u1.Name)
	err = ss.Scan(&ucheck.Name)
	if err == nil {
		str := userall{Name: ucheck.Name}
		save, err1 := json.Marshal(&str)
		if err1 != nil {
			fmt.Println("err = ", err)
		}
		fmt.Fprintln(w, string(save))
	}
	if err != nil {
		ss = db.QueryRow("SELECT email FROM userinfo WHERE email=$1", u1.Email)
		err = ss.Scan(&ucheck.Email)
		if err == nil {
			str := userall{Email: ucheck.Email}
			save, err1 := json.Marshal(&str)
			if err1 != nil {
				fmt.Println("err = ", err)
				return
			}
			fmt.Fprintln(w, string(save))
		}
		if err != nil {
			ss = db.QueryRow("SELECT phone FROM userinfo WHERE phone=$1", u1.Phone)
			err = ss.Scan(&ucheck.Phone)
			if err == nil {
				str := userall{Phone: ucheck.Phone}
				save, err1 := json.Marshal(&str)
				if err1 != nil {
					fmt.Println("err = ", err)
					return
				}
				fmt.Fprintln(w, string(save))
			}
			if err != nil {
				_, err = db.Exec(
					"INSERT INTO userinfo (name,email,phone,password) VALUES ($1,$2,$3,$4)",
					u1.Name,
					u1.Email,
					u1.Phone,
					u1.Password,
				)
				if err != nil {
					fmt.Println("err = ", err)
					return
				}
				u1.Role = 0

				s := []byte(u1.Name)
				if string(s[0:5]) == "admin" {
					u1.Role = 1
				}
				Response1(w, u1)
			}
		}
	}
}

//登录
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	//查询
	var u0 uesrlogin
	json.Unmarshal([]byte(body), &u0)
	var ucheck user
	ss := db.QueryRow("SELECT name, email, phone, password FROM userinfo WHERE name=$1", u0.Acountnum)
	err = ss.Scan(&ucheck.Name, &ucheck.Email, &ucheck.Phone, &ucheck.Password)
	if err == nil && u0.Password == ucheck.Password {
		ucheck.Role = 0
		s := []byte(ucheck.Name)
		if string(s[0:5]) == "admin" {
			ucheck.Role = 1
		}
		Response1(w, ucheck)
	}

	if err != nil || u0.Acountnum != ucheck.Name {
		ss = db.QueryRow("SELECT name, email, phone, password FROM userinfo WHERE email=$1", u0.Acountnum)
		err = ss.Scan(&ucheck.Name, &ucheck.Email, &ucheck.Phone, &ucheck.Password)

		if err == nil && u0.Password == ucheck.Password {
			ucheck.Role = 0
			s := []byte(ucheck.Name)
			if string(s[0:5]) == "admin" {
				ucheck.Role = 1
			}

			Response1(w, ucheck)

		}
		if err != nil || u0.Acountnum != ucheck.Email {
			ss = db.QueryRow("SELECT name, email, phone, password FROM userinfo WHERE phone=$1", u0.Acountnum)
			err = ss.Scan(&ucheck.Name, &ucheck.Email, &ucheck.Phone, &ucheck.Password)

			if err == nil && u0.Password == ucheck.Password {
				ucheck.Role = 0
				s := []byte(ucheck.Name)
				if string(s[0:5]) == "admin" {
					ucheck.Role = 1
				}

				Response1(w, ucheck)
			}
		}
	}
}

//注销
func logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	var u1 username
	json.Unmarshal([]byte(body), &u1)

	for i, a := range userList {
		if a == u1.Name {
			userList = append(userList[:i], userList[i+1:]...)
			break
		}
	}

	for client := range cManager.clients {
		if u1.Name == client.id {
			cManager.unregister <- client
		}
	}

}

//Response1 : 制作token并返回
func Response1(w http.ResponseWriter, user user) {

	tokenString := user.Name + user.Email

	strcurr := userall{Name: user.Name, Email: user.Email, Phone: user.Phone, Password: user.Password, Token: tokenString, Role: user.Role}
	save, err := json.Marshal(&strcurr)
	if err != nil {
		fmt.Println("err = ", err)
	}
	fmt.Fprintln(w, string(save))
}

//用户管理刷新获取所有成员信息
func updateuser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	var u1 userrole
	json.Unmarshal([]byte(body), &u1)

	if u1.Role == 1 {
		rows, err := db.Query("SELECT name, email, phone FROM userinfo")

		usercc := []binfo{}
		for rows.Next() {
			var userobj binfo
			err = rows.Scan(&userobj.Name, &userobj.Email, &userobj.Phone)
			if err != nil {
				fmt.Println("err = ", err)
				return
			}
			usercc = append(usercc, userobj)
		}
		buf, _ := json.Marshal(usercc)
		fmt.Fprintln(w, string(buf))
	}

	if u1.Role == 0 {
		rows, err := db.Query("SELECT name FROM userinfo")

		usercc := []username{}
		for rows.Next() {
			var userobj username
			err = rows.Scan(&userobj.Name)
			if err != nil {
				fmt.Println("err = ", err)
				return
			}
			usercc = append(usercc, userobj)
		}
		buf, _ := json.Marshal(usercc)
		fmt.Fprintln(w, string(buf))
	}

}

//用户管理删除
func userDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "content-type")

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	//查询
	var u0 username
	json.Unmarshal([]byte(body), &u0)

	_, err = db.Exec("Delete FROM userinfo WHERE name=$1", u0.Name)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
}

//用户管理添加
func add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	var u1 user
	json.Unmarshal([]byte(body), &u1)

	//查询
	var ucheck user
	ss := db.QueryRow("SELECT name FROM userinfo WHERE name=$1", u1.Name)
	err = ss.Scan(&ucheck.Name)
	if err == nil {
		str := binfo{Name: ucheck.Name}
		save, err1 := json.Marshal(&str)
		if err1 != nil {
			fmt.Println("err = ", err)
		}
		fmt.Fprintln(w, string(save))
	}
	if err != nil {
		ss = db.QueryRow("SELECT email FROM userinfo WHERE email=$1", u1.Email)
		err = ss.Scan(&ucheck.Email)
		if err == nil {
			str := binfo{Email: ucheck.Email}
			save, err1 := json.Marshal(&str)
			if err1 != nil {
				fmt.Println("err = ", err)
				return
			}
			fmt.Fprintln(w, string(save))
		}
		if err != nil {
			ss = db.QueryRow("SELECT phone FROM userinfo WHERE phone=$1", u1.Phone)
			err = ss.Scan(&ucheck.Phone)
			if err == nil {
				str := binfo{Phone: ucheck.Phone}
				save, err1 := json.Marshal(&str)
				if err1 != nil {
					fmt.Println("err = ", err)
					return
				}
				fmt.Fprintln(w, string(save))
			}
			if err != nil {
				_, err = db.Exec(
					"INSERT INTO userinfo (name,email,phone,password) VALUES ($1,$2,$3,$4)",
					u1.Name,
					u1.Email,
					u1.Phone,
					u1.Password,
				)
				if err != nil {
					fmt.Println("err = ", err)
					return
				}
			}
		}
	}
}

//用户管理修改
func modify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	var u1 userall
	json.Unmarshal([]byte(body), &u1)

	var ucheck user
	ss := db.QueryRow("SELECT name FROM userinfo WHERE email=$1", u1.Email)
	err = ss.Scan(&ucheck.Name)
	if err == nil && ucheck.Name != u1.Name {
		str := userall{Name: ucheck.Name}
		save, err1 := json.Marshal(&str)
		if err1 != nil {
			fmt.Println("err = ", err)
		}
		fmt.Fprintln(w, string(save))
		return
	}

	_, err = db.Exec("UPDATE userinfo SET email=$1, password=$2 WHERE name=$3", u1.Email, u1.Password, u1.Name)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
}

//用户管理刷新页面查看是否被管理员删除并确定其权限
func checkF5(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("content-type", "application/json")             //返回数据格式是json

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	var u1 username
	var u2 username
	json.Unmarshal([]byte(body), &u1)

	ss := db.QueryRow("SELECT name FROM userinfo WHERE name=$1", u1.Name)
	err = ss.Scan(&u2.Name)
	if err != nil {
		return
	}
	s := []byte(u1.Name)
	if string(s[0:5]) == "admin" {
		str := userrole{Role: 1}
		save, err1 := json.Marshal(&str)
		if err1 != nil {
			fmt.Println("err = ", err1)
		}
		fmt.Fprintln(w, string(save))
		return
	}
	str := userrole{Role: 0}
	save, err1 := json.Marshal(&str)
	if err1 != nil {
		fmt.Println("err = ", err1)
	}
	fmt.Fprintln(w, string(save))

}

//聊天页面检查当前聊天用户是否被管理员删除
func checkroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	var u1 username
	json.Unmarshal([]byte(body), &u1)

	var ucheck username
	ss := db.QueryRow("SELECT name FROM userinfo WHERE name=$1", u1.Name)
	err = ss.Scan(&ucheck.Name)
	if err != nil {
		str := username{Name: ucheck.Name}
		save, err1 := json.Marshal(&str)
		if err1 != nil {
			fmt.Println("err = ", err)
			return
		}
		fmt.Fprintln(w, string(save))
	}
}

func main() {
	//连接数据库
	err := initDB()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	fmt.Println("连接数据库成功")

	//ws控制器不断处理管道数据，并进行同步处理
	go cManager.run()

	http.HandleFunc("/ValidToken", ValidToken) //验证TOKEN
	http.HandleFunc("/register", register)     //注册
	http.HandleFunc("/login", login)           //登录
	http.HandleFunc("/logout", logout)         //注销
	http.HandleFunc("/updateuser", updateuser) //用户管理页面刷新获取所有成员信息
	http.HandleFunc("/userDelete", userDelete) //用户管理页面删除
	http.HandleFunc("/add", add)               //用户管理页面添加
	http.HandleFunc("/modify", modify)         //用户管理页面修改
	http.HandleFunc("/checkF5", checkF5)       //用户管理页面F5刷新页面查看是否被管理员删除，并确定其权限
	http.HandleFunc("/checkroom", checkroom)   //聊天页面检查当前聊天用户是否被管理员删除
	http.HandleFunc("/ws", ws)

	err = http.ListenAndServe(":9000", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return
	}

}

type msg struct {
	From string `json:"from"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
	Date string `json:"date"`
	Img  string `json:"img"`
}

type client struct {
	//作为clientmanager记录名字map的key值
	id   string
	ws   *websocket.Conn
	send chan []byte
	data *msg
}

type clientmanager struct {
	//client注册连接器
	clients map[*client]bool
	//记录名字
	ftname map[string]string
	//从连接器发送消息
	chat chan []byte
	//从连接器注册请求
	register chan *client
	//销毁请求
	unregister chan *client
}

var cManager = clientmanager{
	//client注册连接器
	clients: make(map[*client]bool),
	//记录名字
	ftname: make(map[string]string),
	//从连接器发送消息
	chat: make(chan []byte),
	//从连接器注册请求
	register: make(chan *client),
	//销毁请求
	unregister: make(chan *client),
}

var userList = []string{}

//先实现ws的读和写
//ws连接中写数据
func (c *client) write() {
	defer func() {
		cManager.unregister <- c
		c.ws.Close()
	}()

	//从管道遍历数据
	for message := range c.send {
		c.ws.WriteMessage(websocket.TextMessage, message)
	}
}

//ws连接中读数据
func (c *client) read() {
	defer func() {
		cManager.unregister <- c
		c.ws.Close()
	}()

	for {
		//不断读websocket数据
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			fmt.Println("err = ", err)
			//读不进数据，将用户移除
			cManager.unregister <- c
			c.ws.Close()
			break
		}

		//读取数据
		json.Unmarshal(message, &c.data)
		cManager.ftname[c.data.From] = c.data.From
		if contains(userList, c.data.From) {
			userList = append(userList, c.data.From)
			fmt.Println(c.data.From + "  login successfully!")
			c.id = c.data.From
			fmt.Println(userList)
			//登录成功 查库 离线消息
			rows, err := db.Query("SELECT * FROM chat")
			if err != nil {
				fmt.Println("err = ", err)
			}
			if err == nil {
				for rows.Next() {
					var userobj msg
					rows.Scan(&userobj.From, &userobj.To, &userobj.Date, &userobj.Msg, &userobj.Img)
					if userobj.To == c.data.From {
						_, err = db.Exec(
							"delete from chat where receiver=$1;",
							userobj.To,
						)
						if err != nil {
							fmt.Println("err = ", err)
						}
						buf, _ := json.Marshal(&userobj)
						cManager.chat <- buf
					}
				}
			}
		}
		c.data.Date = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
		dataB, _ := json.Marshal(c.data)
		cManager.chat <- dataB
	}
}

//查询用户列表中是否存在该用户
//避免重复添加
//或查询
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return false
		}
	}
	return true
}

//升级为ws请求
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//ws的回调函数
func ws(w http.ResponseWriter, r *http.Request) {
	//获取ws对象
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer ws.Close()

	//创建连接对象
	//初始化连接对象
	c := &client{ws: ws, send: make(chan []byte, 128), data: &msg{}}

	//在ws中注册
	cManager.register <- c

	go c.write()
	c.read()

	defer func() {
		cManager.unregister <- c
	}()

}

//处理ws的逻辑实现
func (manager *clientmanager) run() {
	//监听管道数据，在后端不断处理管道数据
	for {
		//根据不同的数据管道，处理不同逻辑
		select {
		//注册
		case client := <-manager.register:
			//标志注册了
			manager.clients[client] = true

			//注销
		case client := <-manager.unregister:
			if _, ok := manager.clients[client]; ok {

				for i, a := range userList {
					if a == client.id {
						userList = append(userList[:i], userList[i+1:]...)
						break
					}
				}
				close(client.send)
				delete(manager.clients, client)
			}
		case message := <-manager.chat:
			var u1 msg
			err := json.Unmarshal(message, &u1)
			if err != nil {
				fmt.Println("err = ", err)
			}
			if u1.Msg != "" || u1.Img != "" {
				flag := 0
				for _, r := range userList {
					if r == u1.To {
						for client := range manager.clients {
							if r == client.id {
								select {
								//如果成功向client.send写入message，则进行该case处理语句
								case client.send <- message:
									flag = 1
								default:
									//如果上面没有成功，执行这里
									//防止死循环
									close(client.send)
									delete(manager.clients, client)
								}
							}
						}
					}
				}
				if flag == 0 {
					_, err = db.Exec(
						"insert into chat (sender,receiver,date,msg,img) values ($1,$2,$3,$4,$5);",
						u1.From,
						u1.To,
						u1.Date,
						u1.Msg,
						u1.Img,
					)
					if err != nil {
						fmt.Println("err = ", err)
					}
				}

			}
		}
	}
}
