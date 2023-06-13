package weave

import "git.sr.ht/~ionous/tapestry/lang"



// use the domain rules ( and hierarchy ) to strip determiners off of the passed word
func (d *Domain) StripDeterminer(word string) (retDet, retWord string) {
	// fix: determiners should be specified by the author ( and libraries )
	return lang.SliceArticle(word)
}

// use the domain rules ( and hierarchy ) to strip determiners off of the passed word
func (d *Domain) UniformDeterminer(word string) (retDet, retWord string) {
	// fix: determiners should be specified by the author ( and libraries )
	det, name := lang.SliceArticle(word)
	if name, ok := UniformString(name); ok {
		retDet, retWord = det, name
	}
	return
}
