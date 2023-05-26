package main

import (
	"context"
	"errors"
	"fmt"
	jpeg "image/jpeg"
	png "image/png"
	"os"
	"path/filepath"
)

func main() {
	//これは何？
	ctx := context.Background()
	if err := run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, files []string) error {
	//ファイルを取得
	for _, file := range files {
		//変換関数を呼ぶ
		if err := convert(ctx, file); err != nil {
			return err
		}
	}
	return nil
}

func convert(ctx context.Context, file string) (rerr error) {
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
			rerr = errors.Join(rerr, os.Remove(jpgfile))
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
