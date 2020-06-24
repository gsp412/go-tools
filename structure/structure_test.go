package structure

import "testing"

type ClassUserTeacher struct {
	Id      int64  `gsp:"Id, pk"`
	Name    string `gsp:"Name"`
	UId     int64  `gsp:"Users__Id, pk"`
	UName   string `gsp:"Users__Name"`
	UPNick  string `gsp:"Users__Pkgs__NickName, pk"`
	UPColor string `gsp:"Users__Pkgs__Color"`
	TType   string `gsp:"Teacher__Type, pk"`
	Code    int64  `gsp:"Teacher__Code"`
}

type Class struct {
	Id      int64
	Name    string
	Users   []*User
	Teacher []*Teacher
}

type Package struct {
	NickName string
	Color    string
}

type User struct {
	Id   int64
	Name string
	Pkgs []*Package
}

type Teacher struct {
	Type string
	Code int64
}

func TestStructure(t *testing.T) {
	in := []*ClassUserTeacher{
		{Id: 1, Name: "班级1", UId: 1, UName: "张三", UPNick: "小白", UPColor: "白色", TType: "英语老师", Code: 13},
		{Id: 1, Name: "班级1", UId: 2, UName: "李四", UPNick: "小黑", UPColor: "黑色", TType: "英语老师", Code: 13},
		{Id: 1, Name: "班级1", UId: 3, UName: "王五", UPNick: "小绿", UPColor: "绿色", TType: "英语老师", Code: 13},
		{Id: 1, Name: "班级1", UId: 1, UName: "张三", UPNick: "小蓝", UPColor: "蓝色", TType: "数学老师", Code: 15},
		{Id: 1, Name: "班级1", UId: 2, UName: "李四", UPNick: "小白", UPColor: "白色", TType: "数学老师", Code: 15},
		{Id: 1, Name: "班级1", UId: 3, UName: "王五", UPNick: "小白", UPColor: "白色", TType: "数学老师", Code: 15},
		{Id: 2, Name: "班级2", UId: 4, UName: "赵六", TType: "英语老师", Code: 13},
		{Id: 2, Name: "班级2", UId: 5, UName: "陈七", TType: "英语老师", Code: 13},
		{Id: 2, Name: "班级2", UId: 6, UName: "冯八", TType: "英语老师", Code: 13},
		{Id: 2, Name: "班级2", UId: 4, UName: "赵六", TType: "数学老师", Code: 15},
		{Id: 2, Name: "班级2", UId: 5, UName: "陈七", TType: "数学老师", Code: 15},
		{Id: 2, Name: "班级2", UId: 6, UName: "冯八", TType: "数学老师", Code: 15},
	}

	for i := 0; i < 100000; i++ {
		var out []*Class

		err := Structure(&in, &out)
		if nil != err {
			t.Error(err)
			continue
		}
	}

	t.Logf("SUCCESS \n")
}
