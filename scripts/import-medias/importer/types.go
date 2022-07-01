package importer

type Name struct {
	NConst            string
	PrimaryName       string
	BirthYear         string // year or \N
	DeathYear         string // year or \N
	PrimaryProfession string // comma separated job names
	KnownForTitles    string // comma separated Title.TConst
}

type TitleName struct {
	TitleId         string
	Ordering        int32
	Title           string
	Region          string
	Language        string
	Type            string
	Attributes      string // description of the title or \N
	IsOriginalTitle int8   // 0 for false and 1 for true
}

type Title struct {
	TConst         string
	TitleType      string
	PrimaryTitle   string // string separated by `; or `
	OriginalTitle  string
	IsAdult        int8 // 0 for false and 1 for true
	StartYear      int32
	EndYear        int32  // \N for movies and last episode for series
	RuntimeMinutes int32  // number of \N
	Genres         string // comma separated strings
}

type TitleCrew struct {
	TConst    string
	Directors string // comma separated Name.NConst or \N
	Writers   string // comma separated Name.NConst or \N
}

type TitleEpisode struct {
	TConst        string // Title.TConst of the episode
	ParentTConst  string // Title.TConst of the serie
	SeasonNumber  string // int or \N
	EpisodeNumber string // int or \N
}

type TitlePrincipals struct {
	TConst     string // Title.TConst
	Ordering   int32
	NConst     string // Name.NConst
	category   string // Role of the person
	job        string // Precision of the role
	characters string // Array of strings or \N
}

type TitleRating struct {
	TConst        string  // Title.TConst
	QverageRating float32 // rating between 0 and 10
	NumVotes      int32
}
