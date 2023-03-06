package lib

import (
	"bytes"
	"io"

	"image"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
	"github.com/taubyte/go-sdk/event"

	"github.com/taubyte/go-sdk/http/client"

	_ "image/jpeg"
	"image/png"
)

func get(path string) (io.ReadCloser, error) {
	c, err := client.New()
	if err != nil {
		return nil, err
	}

	r, err := c.Request(path)
	if err != nil {
		return nil, err
	}

	res, err := r.Do()
	if err != nil {
		return nil, err
	}

	return res.Body(), nil
}

//export process
func process(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	imageRC, err := get("https://upload.wikimedia.org/wikipedia/commons/5/55/Mona_Lisa_headcrop.jpg")
	if err != nil {
		h.Write([]byte("Error fetching the image. Failed with " + err.Error()))
		return 1
	}
	defer imageRC.Close()

	img, _, err := image.Decode(imageRC)
	if err != nil {
		h.Write([]byte("Error Deconing the image. Failed with " + err.Error()))
		return 1
	}

	inverted := effect.Invert(img)
	resized := transform.Resize(inverted, 800, 800, transform.Linear)
	rot := transform.Rotate(resized, 60, nil)

	var b bytes.Buffer
	err = png.Encode(&b, rot)
	if err != nil {
		h.Write([]byte("PNG Encoding failed with " + err.Error()))
		return 1
	}

	h.Headers().Set("Content-Type", "image/png")
	h.Write(b.Bytes())

	return 0
}
