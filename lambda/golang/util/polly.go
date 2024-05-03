package util

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
)

type PollyInstance struct {
	PollyClient *polly.Client
}

func NewPollyInstance(config aws.Config) *PollyInstance {
	return &PollyInstance{
		PollyClient: polly.NewFromConfig(config),
	}
}

func (p *PollyInstance) SynthesizeSpeech(voice string, text string) (*polly.SynthesizeSpeechOutput, error) {
	return p.PollyClient.SynthesizeSpeech(context.TODO(), &polly.SynthesizeSpeechInput{
		OutputFormat: "mp3",
		Engine:       "standard",
		TextType:     "ssml",
		Text:         aws.String(text),
		VoiceId:      types.VoiceId(voice),
	})
}
