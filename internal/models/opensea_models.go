package models

type DisplayType string

const (
	NumberDisplayType          DisplayType = "number"
	BoostNumberDisplayType     DisplayType = "boost_number"
	BoostPercentageDisplayType DisplayType = "boost_percentage"
	DateDisplayType            DisplayType = "date"
)

// https://docs.opensea.io/docs/metadata-standards#attributes
type OpenseaNFTAttribute struct {
	// might be either an float64 or a string
	Value interface{} `json:"value"`
	// only allowed if value is float64
	DisplayType DisplayType `json:"display_type,omitempty"`
	TraitType   string      `json:"trait_type"`
}

// struct to fullfil the opensea json interface: https://docs.opensea.io/docs/metadata-standards#metadata-structure
type OpenseaNFT struct {
	// This is the URL to the image of the item.
	// Can be just about any type of image (including SVGs, which will be cached into PNGs by OpenSea), and can be IPFS URLs or paths.
	// We recommend using a 350 x 350 image.
	Image string `json:"image"`
	// This is the URL that will appear below the asset's image on OpenSea and
	// will allow users to leave OpenSea and view the item on your site.
	ExternalUrl string `json:"external_url"`
	// A human readable description of the item. Markdown is supported.
	Description string `json:"description"`
	// Name of the item.
	Name string `json:"name"`
	// These are the attributes for the item, which will show up on the OpenSea page for the item.
	Attributes []OpenseaNFTAttribute `json:"attributes"`
	// Background color of the item on OpenSea. Must be a six-character hexadecimal without a pre-pended #.
	BackgroundColor string `json:"background_color"`
	// A URL to a multi-media attachment for the item. The file extensions GLTF, GLB, WEBM, MP4, M4V, OGV, and OGG are supported,
	// along with the audio-only extensions MP3, WAV, and OGA.
	//
	// Animation_url also supports HTML pages, allowing you to build rich experiences and interactive NFTs using JavaScript canvas, WebGL, and more.
	// Scripts and relative paths within the HTML page are now supported. However, access to browser extensions is not supported.
	AnimationUrl string `json:"animation_url,omitempty"`
	// A URL to a YouTube video.
	YoutubeUrl string `json:"youtube_url,omitempty"`
}
