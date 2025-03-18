package util

import (
	"fmt"
	"os"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"resty.dev/v3"
)

// https://elevenlabs.io/docs/api-reference/text-to-speech/convert
func TextToSpeech(fileName, text string) error {
	c := resty.New()
	defer c.Close()

	_, err := c.R().
		SetSaveResponse(true).
		SetOutputFileName(fileName).
		SetHeader("Content-Type", "application/json").
		SetHeader("xi-api-key", os.Getenv("XI_API_KEY")).
		SetBody(map[string]string{
			"text":     text,
			"model_id": "eleven_multilingual_v2",
		}).
		Post("https://api.elevenlabs.io/v1/text-to-speech/nPczCjzI2devNBz1zQrb?output_format=mp3_44100_128")
	if err != nil {
		return err
	}
	return nil
}

func synthesizeStartedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Println("Synthesis started.")
}

func synthesizingHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Printf("Synthesizing, audio chunk size %d.\n", len(event.Result.AudioData))
}

func synthesizedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Printf("Synthesized, audio length %d.\n", len(event.Result.AudioData))
}

func cancelledHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Println("Received a cancellation.")
}

// https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/get-started-text-to-speech?tabs=windows%2Cterminal&pivots=programming-language-go
func TextToSpeechByAzure(filename, content string) error {
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	audioConfig, err := audio.NewAudioConfigFromDefaultSpeakerOutput()
	if err != nil {
		return fmt.Errorf("Got an error: ", err)
	}
	defer audioConfig.Close()
	speechConfig, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		return fmt.Errorf("Got an error: ", err)
	}
	defer speechConfig.Close()

	speechConfig.SetSpeechSynthesisVoiceName("zh-CN-Xiaochen:DragonHDLatestNeural")
	speechConfig.SetSpeechSynthesisOutputFormat(common.Audio16Khz32KBitRateMonoMp3)

	speechSynthesizer, err := speech.NewSpeechSynthesizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		return fmt.Errorf("Got an error: ", err)
	}
	defer speechSynthesizer.Close()

	speechSynthesizer.SynthesisStarted(synthesizeStartedHandler)
	speechSynthesizer.Synthesizing(synthesizingHandler)
	speechSynthesizer.SynthesisCompleted(synthesizedHandler)
	speechSynthesizer.SynthesisCanceled(cancelledHandler)

	task := speechSynthesizer.SpeakTextAsync(content)
	var outcome speech.SpeechSynthesisOutcome
	select {
	case outcome = <-task:
	case <-time.After(60 * time.Second):
		return fmt.Errorf("Timed out")
	}
	defer outcome.Close()
	if outcome.Error != nil {
		return fmt.Errorf("Got an error: ", outcome.Error)
	}

	if outcome.Result.Reason == common.SynthesizingAudioCompleted {
		fmt.Printf("Speech synthesized to speaker for text [%s].\n", content)
		if err := os.WriteFile(filename, outcome.Result.AudioData, 0644); err != nil {
			return fmt.Errorf("文件保存失败:", err)
		}
		fmt.Printf("语音已保存至 %s\n", filename)
	} else {
		cancellation, _ := speech.NewCancellationDetailsFromSpeechSynthesisResult(outcome.Result)
		fmt.Printf("CANCELED: Reason=%d.\n", cancellation.Reason)

		if cancellation.Reason == common.Error {
			return fmt.Errorf("CANCELED: ErrorCode=%d\nCANCELED: ErrorDetails=[%s]\nCANCELED: Did you set the speech resource key and region values?\n",
				cancellation.ErrorCode,
				cancellation.ErrorDetails)
		}
	}
	return nil
}
