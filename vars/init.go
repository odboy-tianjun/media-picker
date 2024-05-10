package vars

var GlobalFilePathList []string                          // 所有文件的路径
var GlobalFilePath2MimeInfoMap = make(map[string]string) // 文件路径和mime信息的映射关系
var GlobalFilePath2FileNameMap = make(map[string]string) // 文件路径和文件名称的映射关系
var GlobalFilePath2FileExtMap = make(map[string]string)  // 文件路径和文件扩展名的映射关系
var GlobalVideoPathList []string                         // 所有视频文件的路径
var GlobalImagePathList []string                         // 所有图片文件的路径
