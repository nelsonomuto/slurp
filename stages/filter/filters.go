package filter

import (
	"github.com/omeid/slurp"
	"github.com/omeid/slurp/tools/glob"
)

//A Filter stage, will either close or pass files to the next
// slurp.Stage based on the output of the `filter` function.
func FilterFunc(c *slurp.C, filter func(slurp.File) bool) slurp.Stage {
	return func(files <-chan slurp.File, out chan<- slurp.File) {
		for f := range files {
			if filter(f) {
				f.Close()
			} else {
				out <- f
			}
		}
	}
}


//Filters out files based on a pattern, if they match,
// they will be closed, otherwise sent to the output channel.
func Filter(c *slurp.C, pattern string) slurp.Stage {
	return FilterFunc(c, func(f slurp.File) bool {
		s, err := f.Stat()
		if err != nil {
			c.Errorf("Can't get File Stat: %s", err.Error())
			return false
		}
		m, err := glob.Match(pattern, s.Name())
		if err != nil {
			c.Error(err)
		}
		return m
	})
}
