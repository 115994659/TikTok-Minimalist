package util

//#include <stdlib.h>
//int startCmd(const char* cmd){
//	  return system(cmd);
//}
import "C"
import (
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"log"
	"unsafe"
)

type VideoCoverImage struct {
	InputPath  string
	OutputPath string
	StartTime  string
	KeepTime   string
	Filter     string
	FrameCount int64
	debug      bool
}

//ffmpeg所需的参数
const (
	inputVideoPathOption = "-i"
	startTimeOption      = "-ss"
	keepTimeOption       = "-t"
	videoFilterOption    = "-vf"
	formatToImageOption  = "-f"
	autoReWriteOption    = "-y"
	framesOption         = "-frames:v"
)

//文件默认后缀
var (
	defaultVideoSuffix = ".mp4"
	defaultImageSuffix = ".jpg"
)

var videoTemp VideoCoverImage

func NewVideoCoverImage() *VideoCoverImage {
	return &videoTemp
}
func (v *VideoCoverImage) Debug() {
	v.debug = true
}
func paramJoin(s1, s2 string) string {
	return fmt.Sprintf(" %s %s ", s1, s2)
}

func GetDefaultImageSuffix() string {
	return defaultImageSuffix
}
func (v *VideoCoverImage) GetQueryString() (ret string, err error) {
	if v.InputPath == "" || v.OutputPath == "" {
		err = errors.New("输入输出路径未指定")
		return
	}
	ret = config.Info.FfmpegPath
	ret += paramJoin(inputVideoPathOption, v.InputPath)
	ret += paramJoin(formatToImageOption, "image2")
	if v.Filter != "" {
		ret += paramJoin(videoFilterOption, v.Filter)
	}
	if v.StartTime != "" {
		ret += paramJoin(startTimeOption, v.StartTime)
	}
	if v.KeepTime != "" {
		ret += paramJoin(keepTimeOption, v.KeepTime)
	}
	if v.FrameCount != 0 {
		ret += paramJoin(framesOption, fmt.Sprintf("%d", v.FrameCount))
	}
	ret += paramJoin(autoReWriteOption, v.OutputPath)
	return
}
func (v *VideoCoverImage) ExecCommand(cmd string) error {
	if v.debug {
		log.Println(cmd)
	}
	cCmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(cCmd))
	status := C.startCmd(cCmd)
	if status != 0 {
		return errors.New("视频切截图失败")
	}
	return nil
}
