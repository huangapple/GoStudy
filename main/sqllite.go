package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDBName = "adoc"

var code2nameMap map[string]string

type Info struct {
	Name string `json:"name" bson:"name" db:"name"`
	Code string `json:"code" bson:"code" db:"code"`
}
type ADOCInfo struct {
	Info
	StreetCode   string `db:"streetCode" json:"streetCode,omitempty" bson:"streetCode,omitempty"`
	ProvinceCode string `db:"provinceCode" json:"provinceCode,omitempty" bson:"provinceCode,omitempty"`
	CityCode     string `db:"cityCode" json:"cityCode,omitempty" bson:"cityCode,omitempty"`
	AreaCode     string `db:"areaCode" json:"areaCode,omitempty" bson:"areaCode,omitempty"`
}

func (info *ADOCInfo) GetInfo() *Info {
	return &info.Info
}

func (info *ADOCInfo) GetProvince() *Info {
	return &Info{
		Name: code2nameMap[info.ProvinceCode],
		Code: info.ProvinceCode,
	}
}
func (info *ADOCInfo) GetCity() *Info {
	return &Info{
		Name: code2nameMap[info.CityCode],
		Code: info.CityCode,
	}
}
func (info *ADOCInfo) GetArea() *Info {
	return &Info{
		Name: code2nameMap[info.AreaCode],
		Code: info.AreaCode,
	}
}
func (info *ADOCInfo) GetStreet() *Info {
	return &Info{
		Name: code2nameMap[info.StreetCode],
		Code: info.StreetCode,
	}
}

type MongoADOC struct {
	ID           string `json:"_id" bson:"_id"`
	Type         string `json:"type" bson:"type"`
	ProvinceInfo *Info  `json:"provinceInfo,omitempty" bson:"provinceInfo,omitempty"`
	VillageInfo  *Info  `json:"villageInfo,omitempty" bson:"villageInfo,omitempty"`
	StreetInfo   *Info  `json:"streetInfo,omitempty" bson:"streetInfo,omitempty"`
	DistrictInfo *Info  `json:"districtInfo,omitempty" bson:"districtInfo,omitempty"`
	CityInfo     *Info  `json:"cityInfo,omitempty" bson:"cityInfo,omitempty"`
}

func GetJsonStr(i interface{}) string {
	var bytes []byte
	var err error

	bytes, err = json.Marshal(i)

	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func dump(i interface{}) {
	jsonBytes, _ := json.MarshalIndent(i, "", "    ")
	fmt.Println(string(jsonBytes))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func exportFile(fileName string, data string) {
	err := ioutil.WriteFile(fileName, []byte(data), os.ModePerm)
	checkErr(err)
}

func main() {

	code2nameMap = make(map[string]string)
	db, err := sqlx.Connect("sqlite3", "./data.sqlite")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb:27017"))

	checkErr(err)

	err = client.Connect(nil)
	checkErr(err)
	mgo := client.Database("information").Collection(mongoDBName)
	checkErr(err)
	defer db.Close()

	//导出省
	exportProvince(db, mgo)
	//导出市
	exportCity(db, mgo)
	//导出区
	exportDistrict(db, mgo)
	//导出街道
	exportStreet(db, mgo)
	//导出村庄
	exportVillage(db, mgo)

	fmt.Println("全部完成")
}

//导出省份
func exportProvince(db *sqlx.DB, mgo *mongo.Collection) {

	list := make([]*ADOCInfo, 0)

	err := db.Unsafe().Select(&list, "select * from province")
	checkErr(err)

	//adocList := make([]*MongoADOC, 0, len(list))
	for _, info := range list {

		mongoADOC := &MongoADOC{
			ID:           info.Code,
			Type:         "province",
			ProvinceInfo: info.GetInfo(),
		}

		code2nameMap[info.Code] = info.Name

		result, err := mgo.InsertOne(nil, mongoADOC)
		checkErr(err)
		fmt.Println(GetJsonStr(result))
	}

	//txt := fmt.Sprintf("db.%s.insertMany(%s)\n", mongoDBName, GetJsonStr(adocList))
	//exportFile("province.mongo", txt)
}

//导出市
func exportCity(db *sqlx.DB, mgo *mongo.Collection) {

	list := make([]*ADOCInfo, 0)

	err := db.Unsafe().Select(&list, "select * from city")
	checkErr(err)

	//adocList := make([]*MongoADOC, 0, len(list))
	for _, info := range list {

		mongoADOC := &MongoADOC{
			ID:           info.Code,
			Type:         "city",
			ProvinceInfo: info.GetProvince(),
			CityInfo:     info.GetInfo(),
		}

		code2nameMap[info.Code] = info.Name
		result, err := mgo.InsertOne(nil, mongoADOC)
		checkErr(err)
		fmt.Println(GetJsonStr(result))
	}

	//txt := fmt.Sprintf("db.%s.insertMany(%s)\n", mongoDBName, GetJsonStr(adocList))
	//exportFile("city.mongo", txt)
}

//导出区
func exportDistrict(db *sqlx.DB, mgo *mongo.Collection) {

	list := make([]*ADOCInfo, 0)

	err := db.Unsafe().Select(&list, "select * from area")
	checkErr(err)

	//adocList := make([]*MongoADOC, 0, len(list))
	for _, info := range list {

		mongoADOC := &MongoADOC{
			ID:           info.Code,
			Type:         "district",
			ProvinceInfo: info.GetProvince(),
			CityInfo:     info.GetCity(),
			DistrictInfo: info.GetInfo(),
		}

		code2nameMap[info.Code] = info.Name
		result, err := mgo.InsertOne(nil, mongoADOC)
		checkErr(err)
		fmt.Println(GetJsonStr(result))
	}
	//txt := fmt.Sprintf("db.%s.insertMany(%s)\n", mongoDBName, GetJsonStr(adocList))
	//exportFile("district.mongo", txt)
}

//导出街道
func exportStreet(db *sqlx.DB, mgo *mongo.Collection) {

	list := make([]*ADOCInfo, 0)

	err := db.Unsafe().Select(&list, "select * from street")
	checkErr(err)

	//adocList := make([]*MongoADOC, 0, len(list))
	for _, info := range list {

		mongoADOC := &MongoADOC{
			ID:           info.Code,
			Type:         "street",
			ProvinceInfo: info.GetProvince(),
			CityInfo:     info.GetCity(),
			DistrictInfo: info.GetArea(),
			StreetInfo:   info.GetInfo(),
		}

		code2nameMap[info.Code] = info.Name
		result, err := mgo.InsertOne(nil, mongoADOC)
		checkErr(err)
		fmt.Println(GetJsonStr(result))
	}
	//txt := fmt.Sprintf("db.%s.insertMany(%s)\n", mongoDBName, GetJsonStr(adocList))
	//exportFile("street.mongo", txt)
}

//导出村庄
func exportVillage(db *sqlx.DB, mgo *mongo.Collection) {

	list := make([]*ADOCInfo, 0)

	err := db.Unsafe().Select(&list, "select * from village")
	checkErr(err)

	//adocList := make([]*MongoADOC, 0, len(list))
	for _, info := range list {

		mongoADOC := &MongoADOC{
			ID:           info.Code,
			Type:         "village",
			ProvinceInfo: info.GetProvince(),
			CityInfo:     info.GetCity(),
			DistrictInfo: info.GetArea(),
			StreetInfo:   info.GetStreet(),
			VillageInfo:  info.GetInfo(),
		}

		code2nameMap[info.Code] = info.Name
		result, err := mgo.InsertOne(nil, mongoADOC)
		checkErr(err)
		fmt.Println(GetJsonStr(result))
	}
	//txt := fmt.Sprintf("db.%s.insertMany(%s)\n", mongoDBName, GetJsonStr(adocList))
	//exportFile("village.mongo", txt)
}
