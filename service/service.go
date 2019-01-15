package service

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"sync"
)

//定义枚举
type State int
const (
	control_orc_dev ="0"
	control_dev_only="1"
	control_ocr_only="2"
)

//Service service for the service struct
type Service struct {
	Address string
	//MgDb    *mgo.Session
	//Action  map[string]string
	//进行添加python库对应的信息
	Pythonlib *PythonLib
	sync.RWMutex
}
//point
type Point struct {
	//坐标类型
	// 1表示只有绝对坐标有意义；
	//2表示只有相对坐标有意义；
	//3表示绝对坐标和相对坐标都有意义。
	CoordinateType int `json:"CoordinateType"`
	//文字位置由x，y坐标组成，单位像素点
	AbsoluteCoorX  int `json:"AbsoluteCoorX"`
	AbsoluteCoorY  int `json:"AbsoluteCoorY"`
	//相对坐标基准分类标识
	//意义与上面的“识别目标分类标识”、“识别目标唯一标识”相同，是为了上层应用
	RelativeTypeID int `json:"RelativeTypeID"`
	//相对坐标基准唯一标识
	RelativeNumID  int `json:"RelativeNumID"`
	//相对于相对坐标基准的坐标由x，y组成，单位像素点
	RelativeCoorX  int `json:"RelativeCoorX"`
	RelativeCoorY  int `json:"RelativeCoorY"`
}
//文本标注信息
type Text struct {
	//0表示识别文字
	TypeID int 		`json:"TypeID"`
	//标识第几个识别对象
	NumID  int 		`json:"NumID"`
	//文字内容
	Info   string 	`json:"Info"`
	//识别目标对象的图像识别矩形尺寸宽（W）单位像素点
	Wide   int    	`json:"Wide"`
	//识别目标对象的图像识别矩形尺寸高（H），单位像素点
	Hight  int 		`json:"Hight"`
	//points
	Points *Point   `json:"Points"`
}
//AirSwitch
//空开识别结果Json格式
type AirSwitch struct {
	//TypeID
	TypeID int `json:"TypeID"`
	NumID  int `json:"NumID"`
	//空开类型识别（如是1P空开还是2P空开）1表示1P，2表示2P，3表示3P，4表示4P。
	Type   int `json:"Type"`
	Color  int `json:"Color"`
	State  int `json:"State"`
	//空开的健康状态或者工作环境是否符合规范要求等，0表示损坏，1表示正常
	HealthState int `json:"HealthState"`
	Wide	int `json:"Wide"`
	Hight	int `json:"Hight"`
	Points  *Point `json:"Points"`
}
//压板识别结果Json格式
//Ena
type Ena struct {
	TypeID int `json:"TypeID"`
	NumID  int `json:"NumID"`
	Type   int `json:"Type"`
	Color  int `json:"Color"`
	State  int `json:"State"`
	HealthState int `json:"HealthState"`
	Wide   int `json:"Wide"`
	Hight  int `json:"Hight"`
	Points *Point `json:"Points"`
} 
//端子排端子识别结果Json格式
type Terminal struct {
	TypeID int `json:"TypeID"`
	NumID  int `json:"NumID"`
	Type   int `json:"Type"`
	Color  int `json:"Color"`
	State  int `json:"State"`
	HealthState int `json:"HealthState"`
	Wide   int `json:"Wide"`
	Hight  int `json:"Hight"`
	Points *Point `json:"Points"`
} 
//msg
type Data struct {
	//Msg    string      `bson:"message"`
	TextRet  []*Text     `json:"text_ret"`
	AirSwitchRet []*AirSwitch `json:"air_switch_ret"`
	EnaRet   []*Ena    `json:"ena_ret"`
	TerminalRet []*Terminal `json:"terminal_ret"`
	LabelPic string `json:"label_pic"`
}
var gService *Service
//Resp to client
type Resp struct {
	Code   int     `json:"code"`
	Result Data    `json:"result"`
	Msg    string  `json:"msg"`
}
// file download respons
type ErrResp struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
} 
//NewService init the Service struct
func NewService() *Service {
	service := new(Service)
	service.Address = Conf.ServerC.Addr
	//service.Action = make(map[string]string)
	//service.MgDb = Newmgo()
	service.Pythonlib = new(PythonLib)
	return service
}

//Start start the kline service
func (ser *Service) Start() {
	gService = ser
	//action and table
	//go mapped(ser)
	//start the web service
	go initPythonLib(ser)
	go Server(ser)
}


//Close close the kline service
func (ser *Service) Close() {
	//ser.MgDb.Close()
}

//
func initPythonLib(ser *Service)  {
	ser.Pythonlib = NewLib()
	//importModule
	ser.Pythonlib.LoadLib(Conf.PythonC.LibDir,Conf.PythonC.LibName)
	ser.Pythonlib.Init()
}
//Server the server is to do web work
func Server(service *Service) error {

	irs := echo.New()
	irs.POST("/OcrProcess", handler)
	irs.GET("/download",staticServer)
	irs.GET("/test",testPython)
	//add 接口
	//irs.GET("/akbusdt", getAbkHandler)
	irs.Start(service.Address)
	return nil
}
func testPython(ctx echo.Context) error {
	ret :=gService.Pythonlib.GetResult()
	fmt.Println(ret)
	var m2 map[string]Text
	json.Unmarshal([]byte(ret),&m2)
	fmt.Println(m2)
	result :=[]Text{}
	for key,item := range m2 {
		fmt.Println("kkkk===>>>>",key)
		fmt.Println("value==>>>",item)
		result = append(result,item)
	}

	ctx.JSON(http.StatusOK, &result)
	return nil
}
func getTextResult(textRet *[]*Text ) error  {
	ocrResult :=gService.Pythonlib.GetResult()
	var tempResult map[string]Text
	json.Unmarshal([]byte(ocrResult),&tempResult)
	for _,item := range tempResult {
		(*textRet) = append((*textRet),&item)
	}
	return nil
}
func handler(ctx echo.Context) error {
	err := save_pciture(ctx)
	if err != nil {
		Logger.Error("[getData]", zap.String("@getData", "file is null"))
		ctx.JSON(http.StatusOK, &ErrResp{
			Code:   http.StatusBadRequest,
			Msg:    "param is error"})
		return nil
	}
	//进行获取参数
	controlValue := ctx.FormValue("ControlValue")
	if controlValue == "" {
		Logger.Error("[getData]", zap.String("@getData", "ControlValue is error"))
		ctx.JSON(http.StatusOK, &ErrResp{
			Code:   http.StatusBadRequest,
			Msg:    "param is error"})
		return nil
	}
	//
	nullMap := Data{}
	//Text
	//设置标注pic
	nullMap.LabelPic ="LabelPic.png"
	switch controlValue {
	case control_ocr_only:
		getTextResult(&nullMap.TextRet)
	case control_dev_only:
		//AirSwitch
		AirSwitchDatas := &AirSwitch{
			TypeID:1,
			NumID:4,
			Type:1,
			Color:0xfffff,
			State:0,
			HealthState:1,
			Wide:20,
			Hight:30,
			Points:&Point{

			},
		}
		nullMap.AirSwitchRet = append(nullMap.AirSwitchRet,AirSwitchDatas)
		EnaDatas := &Ena{
			TypeID:1,
			NumID:4,
			Type:1,
			Color:0xfffff,
			State:0,
			HealthState:1,
			Wide:20,
			Hight:30,
			Points:&Point{

			},
		}
		nullMap.EnaRet = append(nullMap.EnaRet,EnaDatas)
		TerminalDatas := &Terminal{
			TypeID:1,
			NumID:4,
			Type:1,
			Color:0xfffff,
			State:0,
			HealthState:1,
			Wide:20,
			Hight:30,
			Points:&Point{

			},
		}
		nullMap.TerminalRet = append(nullMap.TerminalRet,TerminalDatas)
	case control_orc_dev:
		getTextResult(&nullMap.TextRet)
		//AirSwitch
		AirSwitchDatas := &AirSwitch{
			TypeID:1,
			NumID:4,
			Type:1,
			Color:0xfffff,
			State:0,
			HealthState:1,
			Wide:20,
			Hight:30,
			Points:&Point{

			},
		}
		nullMap.AirSwitchRet = append(nullMap.AirSwitchRet,AirSwitchDatas)
		EnaDatas := &Ena{
			TypeID:1,
			NumID:4,
			Type:1,
			Color:0xfffff,
			State:0,
			HealthState:1,
			Wide:20,
			Hight:30,
			Points:&Point{

			},
		}
		nullMap.EnaRet = append(nullMap.EnaRet,EnaDatas)
		TerminalDatas := &Terminal{
			TypeID:1,
			NumID:4,
			Type:1,
			Color:0xfffff,
			State:0,
			HealthState:1,
			Wide:20,
			Hight:30,
			Points:&Point{

			},
		}
		nullMap.TerminalRet = append(nullMap.TerminalRet,TerminalDatas)
	default:
		Logger.Error("[getData]", zap.String("@getData", "ControlValue is error"))
		ctx.JSON(http.StatusOK, &ErrResp{
			Code:   http.StatusBadRequest,
			Msg:    "param is error"})
		return nil
	}
	Logger.Info(">>>>>>>>>>>>>>>>>>>>",zap.Any("@nullMap",nullMap))
	ctx.JSON(http.StatusOK, &Resp{
		Code:   http.StatusOK,
		Result: nullMap,
		Msg:    "request ok"})
	return nil
}
func save_pciture(ctx echo.Context) error {
	file ,err:= ctx.FormFile("file")
	if err !=nil {
		return err;
	}
	src,err :=file.Open()
	if err !=nil {
		return err;
	}
	defer src.Close()
	dst, err := os.OpenFile(Conf.ServerC.Dir + file.Filename,os.O_WRONLY | os.O_CREATE,0666)
	if err !=nil {
		return err;
	}
	defer dst.Close()
	if _,err :=io.Copy(dst,src) ;err !=nil {
		return err;
	}
	return nil;
}

//文件下载
func staticServer(ctx echo.Context) error{
	filename := ctx.FormValue("file")
	if filename == "" {
		Logger.Error("[getData]", zap.String("@getData", "filename is null"))
		ctx.JSON(http.StatusOK, &ErrResp{
			Code:   http.StatusBadRequest,
			Msg:    "param is error"})
		return nil
	}
	ctx.Response().Writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	path := Conf.ServerC.Dir + filename
	err := checkExists(path)
	if err != true {
		ctx.JSON(http.StatusOK, &ErrResp{
			Code:   http.StatusBadRequest,
			Msg:    "file "+filename+" is not exit"})
		return nil
	}
	http.ServeFile(ctx.Response().Writer, ctx.Request(), path)
	return nil
}

//查看文件是否存在的接口
func checkExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}














