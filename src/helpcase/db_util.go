package helpcase
import (
	"os"
	"database/sql"
	"log"
	"unicode/utf8"
	"strconv"
)

const isHelpcaseExistQuery = "select id, status from HELPCASE where id=?"
const getHelpcase = "select id, title, date, status, amount, detail_url, donation_url, publish_date from HELPCASE where id=?"
const getAllHelpcase = "select id, title, date, status, amount, detail_url, donation_url, publish_date from HELPCASE order by date desc"
const createCaseQuery = "insert into HELPCASE values (?,?,?,?,?,?,?,'')"
const updateCaseQuery = "update HELPCASE set title=?, date=?, status=?, amount=?, detail_url=?, donation_url=? where id=?"
const CaseClose = "已結案"
const getCaseDetailQuery ="select id from HELPCASE_DETAIL where id=?"
const getAllCaseDetailQuery ="select id, title, date, detail_url, likecount, paracount, wordcount, imgcount, content from HELPCASE_DETAIL"
const createCaseDetail="insert into HELPCASE_DETAIL values (?,?,?,?,?,?,?,?,?)"
const updateFBLike="update HELPCASE_DETAIL set LIKECOUNT=? where id=?"
const updateDonationPublishDate="update HELPCASE set publish_date=? where id=?"
const getDonatorQuery="select id from DONATION_DETAIL where id=? and ROW_NUM=?"
const createDonationDetail="insert into DONATION_DETAIL values (?,?,?,?,?,?,'')"
const getDonatorObject="select id, row_num, name, amount, date, longer_four from DONATION_DETAIL where id=? order by ROW_NUM"

type CaseDetailRecord struct {
	SerialNo		string
	Title			string
	Date			string
	DetailUrl 		string
	LikeCount		int
	WordCount		int
	ParagraphCount	int
	ImgCount		int
	Content			string
}

func IsHelpcaseExist(helpcase *Helpcase) bool {

	var isUpdate = false

	pwd, _ := os.Getwd()
	db, _ := sql.Open("sqlite3", pwd + "/test.db")

	defer db.Close()

	tx, _ := db.Begin()

	stmt1, _ := db.Prepare(isHelpcaseExistQuery)
	defer stmt1.Close()

	var id string
	var status string

	err := stmt1.QueryRow(helpcase.SerialNo).Scan(&id, &status)

	if err != nil {
		isUpdate=true
		log.Printf("Case %s not in database. Inserting ...\n", helpcase.SerialNo)
		stmt2, err := tx.Prepare(createCaseQuery)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt2.Close()

		_, err = stmt2.Exec(helpcase.SerialNo, helpcase.Title, helpcase.Date, helpcase.Status, helpcase.Amount, helpcase.DetailUrl, helpcase.DonationUrl)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		if status != CaseClose {
			isUpdate=true
			log.Printf("Case %s in database and status is open. Updating number ...\n", helpcase.SerialNo)
			stmt3, _ := tx.Prepare(updateCaseQuery)
			defer stmt3.Close()

			_, err = stmt3.Exec(helpcase.Title, helpcase.Date, helpcase.Status, helpcase.Amount, helpcase.DetailUrl, helpcase.DonationUrl, helpcase.SerialNo)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("Case %s in database and status is close. Done\n", helpcase.SerialNo)
		}
	}
	log.Printf("Handle case %s for main table done\n", helpcase.SerialNo)
	tx.Commit()

	return isUpdate
}

func GetHelpcase(serialNo string) *Helpcase {
	pwd, _ := os.Getwd()
	db, _ := sql.Open("sqlite3", pwd + "/test.db")

	defer db.Close()

	stmt1, _ := db.Prepare(getHelpcase)
	defer stmt1.Close()

	var id string
	var title string
	var date string
	var status string
	var amount int
	var detailUrl string
	var donationUrl string
	var publishDate string

	err := stmt1.QueryRow(serialNo).Scan(&id, &title, &date, &status, &amount, &detailUrl, &donationUrl, &publishDate)

	if err != nil {
		log.Println(err)
		return nil
	}

	return &Helpcase{
		SerialNo: id,
		Title: title,
		Date: date,
		Status: status,
		Amount: amount,
		DetailUrl: detailUrl,
		DonationUrl: donationUrl,
		PublishDate: publishDate}
}

func GetAllHelpcase() []*Helpcase {
	pwd, _ := os.Getwd()
	db, _ := sql.Open("sqlite3", pwd + "/test.db")

	defer db.Close()

	rows, err := db.Query(getAllHelpcase)

	if err != nil {
		log.Println(err)
		return nil
	}

	defer rows.Close()

	results := make([]*Helpcase, 0)
	for rows.Next() {
		var id string
		var title string
		var date string
		var status string
		var amount int
		var detailUrl string
		var donationUrl string
		var publishDate string
		rows.Scan(&id, &title, &date, &status, &amount, &detailUrl, &donationUrl, &publishDate)

		results = append(results, &Helpcase{
			SerialNo: id,
			Title: title,
			Date: date,
			Status: status,
			Amount: amount,
			DetailUrl: detailUrl,
			DonationUrl: donationUrl,
			PublishDate: publishDate})
	}

	return results
}

func GetAllHelpcaseDetail() []*CaseDetailRecord {
	pwd, _ := os.Getwd()
	db, _ := sql.Open("sqlite3", pwd + "/test.db")

	defer db.Close()

	rows, err := db.Query(getAllCaseDetailQuery)

	if err != nil {
		log.Println(err)
		return nil
	}

	defer rows.Close()

	results := make([]*CaseDetailRecord, 0)
	for rows.Next() {
		var id string
		var title string
		var date string
		var detailUrl string
		var likecount int
		var paracount int
		var wordcount int
		var imgcount int
		var content string
		rows.Scan(&id, &title, &date, &detailUrl, &likecount, &paracount, &wordcount, &imgcount, &content)

		results = append(results, &CaseDetailRecord{
			SerialNo: id,
			Title: title,
			Date: date,
			DetailUrl: detailUrl,
			LikeCount: likecount,
			ParagraphCount: paracount,
			WordCount: wordcount,
			ImgCount: imgcount,
			Content:content})
	}

	return results
}

func GetAllDonationDetail(serialNo string) []*Donator {
	pwd, _ := os.Getwd()
	db, _ := sql.Open("sqlite3", pwd + "/test.db")

	defer db.Close()

	rows, err := db.Query(getDonatorObject, serialNo)

	if err != nil {
		log.Println(err)
		return nil
	}

	defer rows.Close()

	results := make([]*Donator, 0)
	for rows.Next() {
		var id string
		var row_num int
		var date string
		var name string
		var amount int
		var longer_four int
		rows.Scan(&id, &row_num, &name, &amount, &date, &longer_four)
		results = append(results, &Donator{
			SerialNo: strconv.Itoa(row_num),
			Name: name,
			Amount: amount,
			Date: date,
			LongFour:longer_four})
	}

	return results
}

func CheckOrCreateCaseDetail(helpcase *Helpcase, casedetail *CaseDetail) {
	pwd, _ := os.Getwd()
	db, _ := sql.Open("sqlite3", pwd + "/test.db")

	stmt1, _ := db.Prepare(getCaseDetailQuery)
	defer stmt1.Close()

	defer db.Close()

	tx, _ := db.Begin()

	var id string

	err := stmt1.QueryRow(casedetail.SerialNo).Scan(&id)

	if err != nil {
		stmt2, err := tx.Prepare(createCaseDetail)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt2.Close()

		_, err = stmt2.Exec(
			helpcase.SerialNo,
			helpcase.Title,
			helpcase.Date,
			helpcase.DetailUrl,
			casedetail.LikeCount,
			casedetail.ParagraphCount,
			casedetail.WordCount,
			casedetail.ImgCount,
			casedetail.Content)

		if err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit()
	log.Printf("Updating case detail %s done\n", helpcase.SerialNo)
}

func CheckOrCreateDonationDetail(helpcase *Helpcase, donationDetail *DonationDetail) {
	pwd, _ := os.Getwd()
	db, _ := sql.Open("sqlite3", pwd + "/test.db")

	defer db.Close()
	tx, _ := db.Begin()

	stmt0, _ := tx.Prepare(updateDonationPublishDate)
	defer stmt0.Close()

	_, err2 := stmt0.Exec(donationDetail.PublishDate, helpcase.SerialNo)

	if err2 != nil {
		log.Fatal(err2)
	}

	tx.Commit()

	for i := range donationDetail.Donators {
		donator := donationDetail.Donators[i]

		tx_inner, _ := db.Begin()
		stmt1, _ := tx_inner.Prepare(getDonatorQuery)

		var id string
		err := stmt1.QueryRow(donationDetail.SerialNo, donator.SerialNo).Scan(&id)

		if err != nil {
			stmt2, _ := tx_inner.Prepare(createDonationDetail)

			var islong	int
			if utf8.RuneCountInString(donator.Name)>3 {
				islong = 1
			} else {
				islong = 0
			}

			_, err = stmt2.Exec(helpcase.SerialNo, donator.SerialNo, donator.Name, donator.Amount, donator.Date, islong)
			if err != nil {
				log.Fatal(err)
			}
			stmt2.Close()
		}
		stmt1.Close()
		tx_inner.Commit()
	}

	log.Printf("Updating donation detail %s done\n", helpcase.SerialNo)
}

func CleanDB() {
	pwd, _ := os.Getwd()
	db, _ := sql.Open("sqlite3", pwd + "/test.db")

	defer db.Close()
	tx, _ := db.Begin()

	stmt0, _ := tx.Prepare("delete from HELPCASE")
	defer stmt0.Close()
	stmt0.Exec()

	stmt1, _ := tx.Prepare("delete from HELPCASE_DETAIL")
	defer stmt1.Close()
	stmt1.Exec()

	stmt2, _ := tx.Prepare("delete from DONATION_DETAIL")
	defer stmt2.Close()
	stmt2.Exec()

	tx.Commit()
}