package screenshot

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func Capture(ctx context.Context, url string) ([]byte, error) {
	taskCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	taskCtx, cancel = context.WithTimeout(taskCtx, 15*time.Second)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(taskCtx, captureScreenshot(url, &buf)); err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %v", err)
	}

	// filename := fmt.Sprintf("%s.png", formatting.UrlToSlug(url))
	// writeScreenshot(filename, buf)
	// fmt.Printf("screenshot saved: %s\n", filename)

	return buf, nil
}

func writeScreenshot(filename string, buf []byte) error {
	if err := os.WriteFile(filename, buf, 0644); err != nil {
		return fmt.Errorf("failed to save screenshot: %v", err)
	}
	return nil
}

func captureScreenshot(url string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.EmulateViewport(1200, 630),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(1200 * time.Millisecond)
			return nil
		}),
		chromedp.CaptureScreenshot(res),
	}
}
