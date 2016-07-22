package middlewares

import (
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"

	"github.com/reekoheek/go-bono"
)

type cache struct {
	Hit bool
	Val []byte
	Err error
}

func StaticMiddleware(base string) bono.Middleware {
	cacheBag := map[string]*cache{}

	return func(context *bono.Context, next bono.Next) error {
		// should we make it safe from ../..
		path := string(context.Path())
		context.SetContentType(mime.TypeByExtension(filepath.Ext(path)))

		cacheItem := cacheBag[path]
		if cacheItem == nil {
			filePath := filepath.Join(base, path)
			if stat, err := os.Stat(filePath); os.IsNotExist(err) || stat.IsDir() {
				cacheBag[path] = &cache{}
				return next()
			}

			body, err := ioutil.ReadFile(filePath)
			if err != nil {
				cacheBag[path] = &cache{
					Hit: true,
					Err: err,
				}
				return err
			}

			cacheBag[path] = &cache{
				Hit: true,
				Val: body,
			}
			context.SetBody(body)
		} else {
			if cacheItem.Hit == false {
				return next()
			} else if cacheItem.Err != nil {
				return cacheItem.Err
			} else {
				context.SetBody(cacheItem.Val)
			}
		}
		return nil
	}
}
