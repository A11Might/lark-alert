package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ConcatMP3Files 使用ffmpeg将多个MP3文件合并为一个
// inputFiles: 输入文件路径切片
// outputFile: 输出文件路径
// 返回执行日志和可能发生的错误
func ConcatMP3Files(inputFiles []string, outputFile string) (string, error) {
	// 检查输入参数
	if len(inputFiles) == 0 {
		return "", fmt.Errorf("没有提供输入文件")
	}

	// 使用|连接所有输入文件
	concatStr := "concat:" + strings.Join(inputFiles, "|")

	// 创建ffmpeg命令
	cmd := exec.Command("ffmpeg", "-i", concatStr, "-acodec", "copy", outputFile)

	// 创建缓冲区捕获标准输出和错误输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("执行FFmpeg命令失败: %v\n错误输出: %s", err, stderr.String())
	}

	// 返回执行日志，ffmpeg通常将信息输出到stderr，即使成功执行
	return stderr.String(), nil
}

// ConvertToOpus 将MP3文件转换为Opus格式
// 参数:
//   - inputFile: 输入文件路径
//   - outputFile: 输出文件路径
//
// 返回:
//   - 执行结果输出和可能的错误
func ConvertToOpus(inputFile, outputFile string) (string, error) {
	cmd := exec.Command("ffmpeg", "-y", "-i", inputFile, "-acodec", "libopus", "-ac", "1", "-ar", "16000", outputFile)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("执行FFmpeg命令失败: %v\n错误输出: %s", err, stderr.String())
	}

	// FFmpeg通常输出到stderr，即使成功执行
	return stderr.String(), nil
}

// GetAudioDuration 获取音频文件的时长（秒）
// 参数:
//   - filePath: 音频文件路径
//
// 返回:
//   - 时长(秒)和可能的错误
func GetAudioDuration(filePath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("执行FFprobe命令失败: %v\n错误输出: %s", err, stderr.String())
	}

	durationStr := strings.TrimSpace(stdout.String())
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("解析时长失败: %v", err)
	}

	return duration, nil
}
