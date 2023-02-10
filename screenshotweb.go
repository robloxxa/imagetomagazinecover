package main

import (
	"context"
	"github.com/chromedp/chromedp"
)

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}

func processScreenshotMessage(url string) (buf []byte, err error) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()
	if err := chromedp.Run(ctx, fullScreenshot(url, 100, &buf)); err != nil {
		return nil, err
	}
	return
}
