package main

import (
	"context"
	"fmt"
	jpeg "image/jpeg"
	png "image/png"
	"os"
	"path/filepath"
	"runtime/trace"
	"sync"
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	// ワイルドカードを展開してファイルリストを取得
	files, err := filepath.Glob(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := run(ctx, files); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, files []string) error {
	if err := convertAll(ctx, files); err != nil {
		return err
	}
	return nil
}

func convertAll(ctx context.Context, files []string) error {
	var wg sync.WaitGroup
	ctx, task := trace.NewTask(ctx, "convert All")
	defer task.End()

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			//convertを呼ぶ
			if err := convert(ctx, file); err != nil {
				fmt.Println(err) //エラーを表示す
			}
		}(file)
	}
	wg.Wait()
	//成功したらnilを返す
	return nil
}

func convert(ctx context.Context, file string) (rerr error) {
	region := trace.StartRegion(ctx, "convert")
	defer region.End()

	//以下は元のコードと同じ...

	//ファイルを開く
	src, err := os.Open(file)
	if err != nil {
		return err
	}
	//閉じるのを予約
	defer src.Close()
	//ピング画像をデコード
	pngimg, err := png.Decode(src)
	//エラーがあればその都度返す
	if err != nil {
		return err
	}

	//ファイル名を変更
	ext := filepath.Ext(file)
	//拡張しを変更".ping"->".jpg"に
	jpgfile := file[:len(file)-len(ext)] + ".jpg"

	//jpgファイルを作成
	dst, err := os.Create(jpgfile)
	if err != nil {
		return err
	}
	defer func() {
		dst.Close()
		if rerr != nil {
			//失敗したらファイルを削除
			os.Remove(jpgfile)
		}
	}()
	//jpgファイルにエンコード
	if err := jpeg.Encode(dst, pngimg, nil); err != nil {
		return err
	}

	//dstを同期
	if err := dst.Sync(); err != nil {
		return err
	}
	//成功したらnilを返す
	return nil
}
