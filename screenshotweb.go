package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/dchest/uniuri"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

type ScreenshotJob struct {
	MagazineRequest
	QualitiesSlice []string
	C              chan string
}

type ScreenshotDispatcher struct {
	maxWorkers int

	WorkerPool chan chan ScreenshotJob
}

func NewScreenshotDispatcher(maxWorkers int) *ScreenshotDispatcher {
	pool := make(chan chan ScreenshotJob, maxWorkers)
	return &ScreenshotDispatcher{WorkerPool: pool, maxWorkers: maxWorkers}
}

func (d *ScreenshotDispatcher) Run(ctx context.Context) {
	chromeCtx, cancel := chromedp.NewContext(ctx)

	if err := chromedp.Run(chromeCtx); err != nil {
		cancel()
		panic(err)
	}

	for i := 0; i < d.maxWorkers; i++ {
		worker := NewScreenshotWorker(d.WorkerPool)
		worker.Start(chromeCtx)
	}

	go d.dispatch(chromeCtx)
}

func (d *ScreenshotDispatcher) dispatch(ctx context.Context) {
	for {
		select {
		case job := <-ScreenshotJobQueue:
			go func(job ScreenshotJob) {
				jobChannel := <-d.WorkerPool

				jobChannel <- job
			}(job)
		case <-ctx.Done():
			return
		}
	}
}

type ScreenshotWorker struct {
	WorkerPool chan chan ScreenshotJob
	JobChannel chan ScreenshotJob
	quit       chan bool
}

func NewScreenshotWorker(workerPool chan chan ScreenshotJob) ScreenshotWorker {
	return ScreenshotWorker{
		WorkerPool: workerPool,
		JobChannel: make(chan ScreenshotJob),
		quit:       make(chan bool),
	}
}

func (w ScreenshotWorker) Start(ctx context.Context) {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				urlStr, err := makeScreenshot(ctx, job)
				if err != nil {
					log.Println(err)
				}

				job.C <- urlStr
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (w ScreenshotWorker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func makeScreenshot(ctx context.Context, job ScreenshotJob) (imageUrl string, err error) {
	link, err := getLinkToPhoto(job.MagazineRequest)
	if err != nil {
		return
	}

	urlObj, err := url.Parse(fmt.Sprintf("%s/static/magazine.html?img_src=%s", PUBLIC_URL, link))
	if err != nil {
		return
	}
	values := urlObj.Query()
	for k, v := range job.QualitiesSlice {
		values.Set("qualities_"+strconv.Itoa(k), v)
	}

	urlObj.RawQuery = values.Encode()
	buf, err := processScreenshotMessage(ctx, urlObj.String())
	if err != nil {
		return
	}

	s := uniuri.NewLen(32)

	file, err := os.Create("./static/images/covers/" + s + ".png")
	if err != nil {
		panic(err)
		return
	}

	defer file.Close()

	_, err = file.Write(buf)
	if err != nil {
		return
	}

	imageUrl = PUBLIC_URL + "/static/images/covers/" + s + ".png"
	return
}

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(VIEWPORT_WIDTH, VIEWPORT_HEIGHT),
		chromedp.Navigate(urlstr),
		chromedp.Sleep(150 * time.Millisecond),
		chromedp.FullScreenshot(res, quality),
	}
}

func processScreenshotMessage(ctx context.Context, url string) (buf []byte, err error) {
	tabctx, cancel := chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()
	if err := chromedp.Run(tabctx, fullScreenshot(url, 100, &buf), page.Close()); err != nil {
		return nil, err
	}
	return
}
