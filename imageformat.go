package imageformat

import (
	"math"
	"os"
	"path/filepath"
)

type Photo struct {
	Filepath    string
	Orientation string
	Ratio       float64
	Width       int
}

type Operation struct {
	Rotate     string "none"
	Centercrop bool
	Fillsides  bool
	Scaledown  bool
}

func RoundToInt(x float64) int {
	t := math.Trunc(x)
	if math.Abs(x-t) >= 0.5 {
		return int(t + math.Copysign(1, x))
	}
	return int(t)
}

func Visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			//log.Fatal(err)
			return err
		} else {
			*files = append(*files, path)
			return nil
		}
	}
}

func Filter(this Photo, base_ratio float64, base_width float64, keep_h_ratio bool) (steps Operation) {
	// Determine operations to be applied on each image

	// FIRST STEP Determine rotation required
	if this.Orientation == "upper-right" {
		steps.Rotate = "cw" // will rotate clock-wise
	} else {
		if this.Orientation == "lower-left" {
			steps.Rotate = "acw" // will rotate anti-clock-wise
		} else {
			steps.Rotate = "none" // will not rotate
		}
	}

	if steps.Rotate == "acw" || steps.Rotate == "cw" {
		// All rotated images will be vertical
		steps.Fillsides = true
		// LAST STEP Check scaling post rotation
		if (float64(this.Width) * base_ratio) >= base_width {
			steps.Scaledown = true
		}
	} else {
		// SECOND STEP Check for wide images
		if this.Ratio > base_ratio {
			if keep_h_ratio == true {
				steps.Fillsides = true // respect original frame and ratio, fillsides will also fill top/bottom
			} else {
				steps.Centercrop = true // crop wide images to fit base ratio
			}
			// LAST STEP Check for big wide images
			if this.Width > int(base_width) {
				steps.Scaledown = true
			}
		} else {
			// THIRD STEP Check for vertical images
			if this.Ratio <= 1.0 {
				steps.Fillsides = true
				// LAST STEP Check for big vertical images
				if (float64(this.Width) * base_ratio) > base_width {
					steps.Scaledown = true
				}
			}
		}

	}
	return steps
}
