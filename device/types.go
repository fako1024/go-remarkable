package device

type File struct {
	Name string `json:"displayName"`
	ID   string `json:"id"`
}

type Folder struct {
	Name    string             `json:"displayName"`
	ID      string             `json:"id"`
	Files   []File             `json:"files"`
	Folders map[string]*Folder `json:"folders"`
}

type Tree = Folder

type Element struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Children []Element `json:"children"`
}

func (f *Folder) Flatten() []Element {

	var res []Element
	for _, v := range f.Files {
		res = append(res, Element{
			ID:   v.ID,
			Name: v.Name,
			Type: "file",
		})
	}
	for _, v := range f.Folders {
		res = append(res, Element{
			ID:       v.ID,
			Name:     v.Name,
			Type:     "folder",
			Children: v.Flatten(),
		})
	}

	return res
}

// Document denotes a document / file (name + content)
type Document struct {
	Name    string
	Content []byte
}

// Documents denotes several documents
type Documents []Document

type FileMetaData struct {
	Deleted          bool   `json:"deleted"`
	LastModified     string `json:"lastModified"`
	Metadatamodified bool   `json:"metadatamodified"`
	Modified         bool   `json:"modified"`
	Parent           string `json:"parent"`
	Pinned           bool   `json:"pinned"`
	Synced           bool   `json:"synced"`
	Type             string `json:"type"`
	Version          int    `json:"version"`
	VisibleName      string `json:"visibleName"`
}

type FileContentData struct {
	ExtraMetadata  ExtraMetadata `json:"extraMetadata"`
	FileType       string        `json:"fileType"`
	FontName       string        `json:"fontName"`
	LastOpenedPage int           `json:"lastOpenedPage"`
	LineHeight     int           `json:"lineHeight"`
	Margins        int           `json:"margins"`
	PageCount      int           `json:"pageCount"`
	TextScale      int           `json:"textScale"`
	Transform      Transform     `json:"transform"`
}

type Transform struct {
	M11 int `json:"m11"`
	M12 int `json:"m12"`
	M13 int `json:"m13"`
	M21 int `json:"m21"`
	M22 int `json:"m22"`
	M23 int `json:"m23"`
	M31 int `json:"m31"`
	M32 int `json:"m32"`
	M33 int `json:"m33"`
}

type ExtraMetadata struct {
	LastColor      string `json:"LastColor"`
	LastTool       string `json:"LastTool"`
	ThicknessScale string `json:"ThicknessScale"`
}
