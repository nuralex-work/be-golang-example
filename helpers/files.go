package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/gorm/utils"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo/v4"
	excel "github.com/xuri/excelize/v2"
)

func CheckCSVFile(s string) bool {
	extention := s[len(s)-4:]
	return extention == ".csv"

}

//func GetFileToUpload(c echo.Context, param string, path string) (filename string, err error) {
//	if _, err := os.Stat(path); os.IsNotExist(err) {
//		os.Mkdir(path, 0755)
//	}
//
//	file, err := c.FormFile(param)
//	if err != nil {
//		return "", err
//	}
//
//	src, err := file.Open()
//	if err != nil {
//		return "", err
//	}
//	defer src.Close()
//
//	dst, err := os.Create(path + file.Filename)
//	if err != nil {
//		return "", err
//	}
//	defer dst.Close()
//
//	if _, err = io.Copy(dst, src); err != nil {
//		return "", err
//	}
//
//	return file.Filename, err
//
//}

//	func GetJournalEntry(file interface{}, path string) ([]*structs.JournalEntry, error) {
//		defer os.RemoveAll("./" + path + file.(string))
//		filePath := "./" + path + file.(string)
//		csvFile, err := os.Open(filePath)
//		if err != nil {
//			return nil, err
//		}
//		defer csvFile.Close()
//		var journal []*structs.JournalEntry
//
//		gocsv.SetCSVReader(setDelimiter)
//		if err = gocsv.UnmarshalFile(csvFile, &journal); err != nil {
//			return nil, err
//		}
//		return journal, err
//
// }
//
//	func setDelimiter(in io.Reader) gocsv.CSVReader {
//		r := csv.NewReader(in)
//		r.Comma = ';'
//		return r
//
// }
func RemoveFiles(filePath string) (err error) {
	err = os.RemoveAll(filePath)
	if err != nil {
		fmt.Println("Error removing file:", err)
		return
	}

	fmt.Println("File removed successfully!")
	return
}
func FilesToBase64(files *multipart.FileHeader) (string, error) {
	file, err := files.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	// Read the file contents
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	encodedString := base64.StdEncoding.EncodeToString(fileContent)
	return encodedString, nil
}
func ToCharStr(i int) string {
	return string(rune('A' - 1 + i))
}
func UploadFile(c echo.Context, param string, path string) (data string, err error) {
	// Read form fields
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0777)

		if err != nil {
			fmt.Println("tidak bisa create")
		}
	}
	file, err := c.FormFile(param)
	if err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	// Destination
	dst, err := os.Create(path + file.Filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return file.Filename, err
}
func ReadExcelWithHeader(file interface{}, path string, sheetName string, CustomKeys []string) (res []map[string]string, header []string, jmlRows int) {
	//fmt.Println("sheetName", sheetName)
	countCustomkeys := len(CustomKeys)
	// Membaca file excel
	// Deklarasi 2 kali karena menggunakan library yang berbeda
	excel, _ := excel.OpenFile("./" + path + file.(string))
	f, _ := excelize.OpenFile("./" + path + file.(string))

	// Get Tatal Column
	var datacol []string
	cols, err := excel.GetCols(sheetName)
	if err != nil {
		fmt.Println("sheetName :", sheetName)
		return
	} else {
		defer os.RemoveAll("./" + path + file.(string))
	}

	for _, col := range cols {
		for k, rowCell := range col {
			// Deklarasi Kondisi Column yang mau dibaca dengan kondisi cell tidak kososng dan index k merupakan index 0
			if rowCell != "" && k == 0 {
				datacol = append(datacol, rowCell)
			}
		}
	}

	// Deklarasi Total Rows and Cell
	CountRows := len(f.GetRows(sheetName))
	totalCols := len(datacol)
	var totalRows int

	for i := 1; i <= CountRows; i++ {
		if f.GetCellValue(sheetName, fmt.Sprintf("A%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("B%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("G%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("H%d", i)) != "" {
			totalRows++
		}
	}
	index := 2

	//var r = make([]map[string]interface{}, totalRows - 1)
	var r = make([]map[string]string, totalRows-1)
	for i := 2; i <= totalRows; i++ {
		//r[i-2] = make(map[string]interface{})
		r[i-2] = make(map[string]string)
		if f.GetCellValue(sheetName, fmt.Sprintf("A%d", 1)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("A%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("B%d", i)) != "" {
			for j := 1; j <= totalCols; j++ {
				key := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(f.GetCellValue(sheetName, fmt.Sprintf(ToCharStr(j)+"%d", 1)))), " ", "_")
				if CustomKeys != nil && countCustomkeys > 0 && j <= countCustomkeys {
					key = CustomKeys[j-1]
				}
				if i == 2 {
					header = append(header, key)
				}

				value := f.GetCellValue(sheetName, fmt.Sprintf(ToCharStr(j)+"%d", index))
				r[i-2][key] = value
			}
		}

		index++
		jmlRows++
	}

	return r, header, jmlRows
}
func ReadExcelWithHeaderAndDate(file interface{}, path string, sheetName string, CustomKeys []string, dateskey []string) (res []map[string]string, header []string, jmlRows int, errs []string) {
	//fmt.Println("sheetName", sheetName)
	countCustomkeys := len(CustomKeys)
	// Membaca file excel
	// Deklarasi 2 kali karena menggunakan library yang berbeda
	excel, _ := excel.OpenFile("./" + path + file.(string))
	f, _ := excelize.OpenFile("./" + path + file.(string))

	// Get Tatal Column
	var datacol []string
	cols, err := excel.GetCols(sheetName)
	if err != nil {
		fmt.Println("sheetName :", sheetName)
		return
	} else {
		defer os.RemoveAll("./" + path + file.(string))
	}

	for _, col := range cols {
		for k, rowCell := range col {
			// Deklarasi Kondisi Column yang mau dibaca dengan kondisi cell tidak kososng dan index k merupakan index 0
			if rowCell != "" && k == 0 {
				datacol = append(datacol, rowCell)
			}
		}
	}

	// Deklarasi Total Rows and Cell
	CountRows := len(f.GetRows(sheetName))
	totalCols := len(datacol)
	var totalRows int

	for i := 1; i <= CountRows; i++ {
		if f.GetCellValue(sheetName, fmt.Sprintf("A%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("B%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("G%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("H%d", i)) != "" {
			totalRows++
		}
	}
	index := 2

	//var r = make([]map[string]interface{}, totalRows - 1)
	var r = make([]map[string]string, totalRows-1)
	for i := 2; i <= totalRows; i++ {
		//r[i-2] = make(map[string]interface{})
		r[i-2] = make(map[string]string)
		if f.GetCellValue(sheetName, fmt.Sprintf("A%d", 1)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("A%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("B%d", i)) != "" {
			for j := 1; j <= totalCols; j++ {
				key := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(f.GetCellValue(sheetName, fmt.Sprintf(ToCharStr(j)+"%d", 1)))), " ", "_")
				if CustomKeys != nil && countCustomkeys > 0 && j <= countCustomkeys {
					key = CustomKeys[j-1]
				}
				if i == 2 {
					header = append(header, key)
				}

				value := f.GetCellValue(sheetName, fmt.Sprintf(ToCharStr(j)+"%d", index))

				excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
				// Coba parse value ke float64 untuk mengetahui apakah ini serial number
				var date time.Time
				var serial float64

				if _, err := fmt.Sscanf(value, "%f", &serial); err == nil && utils.Contains(dateskey, key) && value != "" {
					d, er := time.Parse("02-01-06", value)
					if er == nil {
						value = d.Format("02-01-2006")
					} else {
						date = excelEpoch.Add(time.Duration(serial*24) * time.Hour)
						log.Println("date.Year(): ", date.Year())
						if date.Year() <= 1900 || date.Year() > time.Now().Year()+100 {
							errs = append(errs, fmt.Sprintf("Invalid format date on cell %s", fmt.Sprintf(ToCharStr(j)+"%d", index)))
						}
						value = date.Format("02-01-2006")
					}
				}

				r[i-2][key] = value
			}
		}

		index++
		jmlRows++
	}

	return r, header, jmlRows, errs
}
func ReadExcelWithHeaderAndDateInterface(file interface{}, path string, sheetName string, CustomKeys []string, dateskey []string) (res []map[string]interface{}, header []string, jmlRows int, errs []string) {
	//fmt.Println("sheetName", sheetName)
	countCustomkeys := len(CustomKeys)
	// Membaca file excel
	// Deklarasi 2 kali karena menggunakan library yang berbeda
	excel, _ := excel.OpenFile("./" + path + file.(string))
	f, _ := excelize.OpenFile("./" + path + file.(string))

	// Get Tatal Column
	var datacol []string
	cols, err := excel.GetCols(sheetName)
	if err != nil {
		fmt.Println("sheetName :", sheetName)
		return
	} else {
		defer os.RemoveAll("./" + path + file.(string))
	}

	for _, col := range cols {
		for k, rowCell := range col {
			// Deklarasi Kondisi Column yang mau dibaca dengan kondisi cell tidak kososng dan index k merupakan index 0
			if rowCell != "" && k == 0 {
				datacol = append(datacol, rowCell)
			}
		}
	}

	// Deklarasi Total Rows and Cell
	CountRows := len(f.GetRows(sheetName))
	totalCols := len(datacol)
	var totalRows int

	for i := 1; i <= CountRows; i++ {
		if f.GetCellValue(sheetName, fmt.Sprintf("A%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("B%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("G%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("H%d", i)) != "" {
			totalRows++
		}
	}
	index := 2

	var r = make([]map[string]interface{}, totalRows-1)
	//var r = make([]map[string]string, totalRows-1)
	for i := 2; i <= totalRows; i++ {
		r[i-2] = make(map[string]interface{})
		//r[i-2] = make(map[string]string)
		if f.GetCellValue(sheetName, fmt.Sprintf("A%d", 1)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("A%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("B%d", i)) != "" {
			for j := 1; j <= totalCols; j++ {
				key := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(f.GetCellValue(sheetName, fmt.Sprintf(ToCharStr(j)+"%d", 1)))), " ", "_")
				if CustomKeys != nil && countCustomkeys > 0 && j <= countCustomkeys {
					key = CustomKeys[j-1]
				}
				if i == 2 {
					header = append(header, key)
				}

				value := f.GetCellValue(sheetName, fmt.Sprintf(ToCharStr(j)+"%d", index))

				excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
				// Coba parse value ke float64 untuk mengetahui apakah ini serial number
				var date time.Time
				var serial float64
				if _, err := fmt.Sscanf(value, "%f", &serial); err == nil && utils.Contains(dateskey, key) && value != "" {
					date = excelEpoch.Add(time.Duration(serial*24) * time.Hour)
					if date.Year() < 1900 || date.Year() > time.Now().Year()+100 {
						errs = append(errs, fmt.Sprintf("Invalid format date on cell %s", fmt.Sprintf(ToCharStr(j)+"%d", index)))
					}
					value = date.Format("02-01-2006")
				}

				r[i-2][key] = value
			}
		}

		index++
		jmlRows++
	}

	return r, header, jmlRows, errs
}

func ReadExcelWithHeaderV2(file interface{}, path string, sheetName string, CustomKeys []string) (res []map[string]interface{}, header []string, jmlRows int) {
	defer os.RemoveAll("./" + path + file.(string))
	//fmt.Println("sheetName", sheetName)
	countCustomkeys := len(CustomKeys)
	// Membaca file excel
	// Deklarasi 2 kali karena menggunakan library yang berbeda
	excel, _ := excel.OpenFile("./" + path + file.(string))
	f, _ := excelize.OpenFile("./" + path + file.(string))

	// Get Tatal Column
	var datacol []string
	cols, err := excel.GetCols(sheetName)
	if err != nil {
		fmt.Println("sheetName :", sheetName)
		return
	}

	for _, col := range cols {
		for k, rowCell := range col {
			// Deklarasi Kondisi Column yang mau dibaca dengan kondisi cell tidak kososng dan index k merupakan index 0
			if rowCell != "" && k == 0 {
				datacol = append(datacol, rowCell)
			}
		}
	}

	// Deklarasi Total Rows and Cell
	CountRows := len(f.GetRows(sheetName))
	totalCols := len(datacol)
	var totalRows int

	for i := 1; i <= CountRows; i++ {
		if f.GetCellValue(sheetName, fmt.Sprintf("A%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("B%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("G%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("H%d", i)) != "" {
			totalRows++
		}
	}
	index := 2

	var r = make([]map[string]interface{}, totalRows-1)
	//var r = make([]map[string]string, totalRows-1)
	for i := 2; i <= totalRows; i++ {
		r[i-2] = make(map[string]interface{})
		//r[i-2] = make(map[string]string)
		if f.GetCellValue(sheetName, fmt.Sprintf("A%d", 1)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("A%d", i)) != "" || f.GetCellValue(sheetName, fmt.Sprintf("B%d", i)) != "" {
			for j := 1; j <= totalCols; j++ {
				key := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(f.GetCellValue(sheetName, fmt.Sprintf(ToCharStr(j)+"%d", 1)))), " ", "_")
				if CustomKeys != nil && countCustomkeys > 0 && j <= countCustomkeys {
					key = CustomKeys[j-1]
				}
				if i == 2 {
					header = append(header, key)
				}

				value := f.GetCellValue(sheetName, fmt.Sprintf(ToCharStr(j)+"%d", index))
				r[i-2][key] = value
			}
		}

		index++
		jmlRows++
	}

	return r, header, jmlRows
}
func UploadHandler(c echo.Context, param string, path string, maxSize int64) (string, error) {
	// Read form file
	maxFileSize := maxSize * 1024 * 1024

	file, err := c.FormFile(param)
	if err != nil {
		return "", errors.New("Failed to get image: " + err.Error())
	}
	if file.Size > maxFileSize {
		return "", errors.New(fmt.Sprintf("File size exceeds %vMB limit", maxSize))
	}

	src, err := file.Open()
	if err != nil {
		return "", errors.New("Failed to open image: " + err.Error())
	}
	defer src.Close()

	// Decode image
	img, format, err := image.Decode(src)
	if err != nil {
		return "", errors.New("Failed to decode image: " + err.Error())
	}
	// Generate random filename with original extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		// fallback extension based on detected format
		switch format {
		case "jpeg":
			ext = ".jpg"
		case "png":
			ext = ".png"
		default:
			ext = ".img"
		}
	}

	// Generate random filename with original extension
	randomName, err := GenerateRandomString(16)
	if err != nil {
		return "", errors.New("Failed to generate filename: " + err.Error())
	}
	filename := randomName + ext
	// Destination
	dstPath := filepath.Join(path, filename)
	// Create uploads directory if not exists
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", errors.New("Failed to create directory: " + err.Error())
	}
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", errors.New("Failed to create file: " + err.Error())
	}
	defer dst.Close()
	// Compress and encode image
	switch format {
	case "jpeg":
		// Compress JPEG with quality 75 (adjust as needed)
		opts := jpeg.Options{Quality: 75}
		if err := jpeg.Encode(dst, img, &opts); err != nil {
			return "", errors.New("Failed to encode JPEG: " + err.Error())
		}
	case "png":
		// For PNG, encode with default compression
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		if err := encoder.Encode(dst, img); err != nil {
			return "", errors.New("Failed to encode PNG: " + err.Error())
		}
	default:
		// For other formats, just save original bytes (no compression)
		// Re-open the file to copy raw bytes
		src2, err := file.Open()
		if err != nil {
			return "", errors.New("Failed to re-open image: " + err.Error())
		}

		defer src2.Close()
		if _, err := io.Copy(dst, src2); err != nil {
			return "", errors.New("Failed to save file: " + err.Error())
		}
	}

	return filename, nil
}
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
func IsImageFile(typefile string) (valid bool) {
	return strings.Contains(typefile, "image")
}
