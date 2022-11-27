package clip

import clip "golang.design/x/clipboard"

func Write(data []byte, f clip.Format) error {
	err := clip.Init()
	if err != nil {
		return err
	}

	clip.Write(f, data)

	return nil
}
