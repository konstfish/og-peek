package chrome

import (
	"context"

	"github.com/chromedp/chromedp"
)

func Initialize(ctx context.Context) (context.CancelFunc, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)

	_, cancel = chromedp.NewContext(allocCtx)

	return cancel, nil
}
