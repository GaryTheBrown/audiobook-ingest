package ffmpeg

type StreamType int

const (
	STNIL StreamType = iota
	Audio
	Video
	Subtitle
)

var streamTypeVals = []string{
	"NULL",
	"a",
	"v",
	"s",
}

func (st *StreamType) ToVal() string {
	return streamTypeVals[*st]
}

type AudioType int

const (
	ATNIL AudioType = iota
	Aaac
)

var audioTypeStrings = []string{
	"NULL",
	"aac",
}

func (at *AudioType) ToString() string {
	return audioTypeStrings[*at]
}

type VideoType int

const (
	VTNULL VideoType = iota
	Vpng
)

var videoTypeStrings = []string{
	"NULL",
	"png",
}

func (vt *VideoType) ToString() string {
	return videoTypeStrings[*vt]
}

type DispositionType int

const (
	DTNULL DispositionType = iota
	Ddefault
	Ddub
	Doriginal
	Dcomment
	Dlyrics
	Dkaraoke
	Dforced
	Dhearing_impaired
	Dvisual_impaired
	Dclean_effects
	Dattached_pic
	Dtimed_thumbnails
	Dnon_diegetic
	Dcaptions
	Ddescriptions
	Dmetadata
	Ddependent
	Dstill_image
	Dmultilayer
)

var dispositionTypeStrings = []string{
	"NULL",
	"default",
	"dub",
	"original",
	"comment",
	"lyrics",
	"karaoke",
	"forced",
	"hearing_impaired",
	"visual_impaired",
	"clean_effects",
	"attached_pic",
	"timed_thumbnails",
	"non_diegetic",
	"captions",
	"descriptions",
	"metadata",
	"dependent",
	"still_image",
	"multilayer",
}

func (dt *DispositionType) ToString() string {
	return dispositionTypeStrings[*dt]
}
