package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dsoprea/go-exif"
	log "github.com/dsoprea/go-logging"
)

type IfdEntry struct {
	IfdPath     string      `json:"ifd_path"`
	FqIfdPath   string      `json:"fq_ifd_path"`
	IfdIndex    int         `json:"ifd_index"`
	TagId       uint16      `json:"tag_id"`
	TagName     string      `json:"tag_name"`
	TagTypeId   uint16      `json:"tag_type_id"`
	TagTypeName string      `json:"tag_type_name"`
	UnitCount   uint32      `json:"unit_count"`
	Value       interface{} `json:"value"`
	ValueString string      `json:"value_string"`
}

type IfdEntryMap map[string]IfdEntry

func Arrange(folderPath, timeformat string) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Program error.")
		}
	}()

	//flag.StringVar(&folderPath, "folderpath", "/Users/cage/Documents/末理整照片(重要)勿刪/aaaaaa", "Folder-path of image")
	//flag.StringVar(&timeformat, "timeformat", "200601", "timeformat for folder (ex: 200601 month, 20060102 month & day)")
	//
	//flag.Parse()
	//
	//if folderPath == "" {
	//	fmt.Printf("Please provide a folder-path for images.\n")
	//	os.Exit(1)
	//}

	mask := syscall.Umask(0) // 改为 0000 八进制
	defer syscall.Umask(mask)

	//
	err := filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				IfdEntries, err := GetPhotoExif(path)
				if err == nil {
					if m, ok := IfdEntries["DateTimeOriginal"]; ok {
						taken, _ := time.Parse("2006:01:02 15:04:05", m.ValueString)

						newFolder := fmt.Sprintf("%s/%s", folderPath, taken.Format(timeformat))
						_ = os.MkdirAll(newFolder, 0766)
						newFile := fmt.Sprintf("%s/%s", newFolder, info.Name())
						os.Rename(path, newFile)

						fmt.Println(fmt.Sprintf("%s → %s", path, newFile))
					}
				}
			}
			return nil
		})

	if err != nil {
		fmt.Print(err)
	}

	RemoveEmptyFolder(folderPath)
}

func GetPhotoExif(fname string) (IfdEntryMap, error) {
	f, err := os.Open(fname)
	if err != nil {
		return IfdEntryMap{}, nil
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return IfdEntryMap{}, nil
	}

	rawExif, err := exif.SearchAndExtractExif(data)
	if err != nil {
		return IfdEntryMap{}, nil
	}

	// Run the parse.
	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	entries := make(IfdEntryMap, 0)
	visitor := func(fqIfdPath string, ifdIndex int, tagId uint16, tagType exif.TagType, valueContext exif.ValueContext) (err error) {
		defer func() {
			if state := recover(); state != nil {
				err = log.Wrap(state.(error))
				log.Panic(err)
			}
		}()

		ifdPath, err := im.StripPathPhraseIndices(fqIfdPath)
		log.PanicIf(err)

		it, err := ti.Get(ifdPath, tagId)
		if err != nil {
			if log.Is(err, exif.ErrTagNotFound) {
				// fmt.Printf("WARNING: Unknown tag: [%s] (%04x)\n", ifdPath, tagId)
				return nil
			} else {
				log.Panic(err)
			}
		}

		valueString := ""
		var value interface{}
		if tagType.Type() == exif.TypeUndefined {
			var err error
			value, err = exif.UndefinedValue(ifdPath, tagId, valueContext, tagType.ByteOrder())
			if log.Is(err, exif.ErrUnhandledUnknownTypedTag) {
				value = nil
			} else if err != nil {
				log.Panic(err)
			} else {
				valueString = fmt.Sprintf("%v", value)
			}
		} else {
			valueString, err = tagType.ResolveAsString(valueContext, true)
			log.PanicIf(err)

			value = valueString
		}

		entry := IfdEntry{
			IfdPath:     ifdPath,
			FqIfdPath:   fqIfdPath,
			IfdIndex:    ifdIndex,
			TagId:       tagId,
			TagName:     it.Name,
			TagTypeId:   tagType.Type(),
			TagTypeName: tagType.Name(),
			UnitCount:   valueContext.UnitCount,
			Value:       value,
			ValueString: valueString,
		}
		entries[it.Name] = entry

		return nil
	}

	_, err = exif.Visit(exif.IfdStandard, im, ti, rawExif, visitor)
	if err != nil {
		return IfdEntryMap{}, nil
	}

	return entries, nil
}
