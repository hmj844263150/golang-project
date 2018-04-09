package api

import (
	"sort"
)

type Seg struct {
	Parent        *Seg
	Name          string
	IsPlaceholder bool
	Ptype         uint8
	Children      []*Seg
	FixedIndex    int
	Api           *Api
}

func (s *Seg) find(seg *Seg) *Seg {
	for _, child := range s.Children {
		if !child.IsPlaceholder && !seg.IsPlaceholder {
			if child.Name == seg.Name {
				return child
			}
		}
		if child.IsPlaceholder && seg.IsPlaceholder {
			if child.Name == seg.Name && child.Ptype == seg.Ptype {
				return child
			}
		}
	}
	return nil
}

func (s *Seg) Lookup(part string) []*Seg {
	i := sort.Search(s.FixedIndex, func(i int) bool {
		return part <= s.Children[i].Name
	})
	if i < s.FixedIndex && part == s.Children[i].Name {
		return []*Seg{s.Children[i]}
	}
	return s.Children[s.FixedIndex:]
}

func (s *Seg) Build() {
	sort.Sort(s)
	done := false
	s.FixedIndex = len(s.Children)
	for i, seg := range s.Children {
		if !done && seg.IsPlaceholder {
			s.FixedIndex = i
			done = true
		}
		seg.Build()
	}
}

func (s *Seg) Len() int {
	return len(s.Children)
}

func (s *Seg) Less(i, j int) bool {
	si, sj := s.Children[i], s.Children[j]
	if !si.IsPlaceholder && sj.IsPlaceholder {
		return true
	}
	if si.IsPlaceholder && !sj.IsPlaceholder {
		return false
	}
	return si.Name < sj.Name
}

func (s *Seg) Swap(i, j int) {
	s.Children[i], s.Children[j] = s.Children[j], s.Children[i]
}
