package main

import (
	"context"
	"encoding/json"
	_ "encoding/json"
	"flag"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/chromedp/chromedp"

	"path/filepath"
	"regexp"
	"strconv"
)

// https://www.artstation.com/artwork/L29Yaw

// *** INFO: Скачивалка файла по урл в облачном хранилище ***
// ПОСЛЕДНЯЯ ТРЕТЬЯ СТАДИЯ ПАРСИНГА

// Сайтец онлайн построитель регэкспов https://regex101.com/r/YCxSpF/6
// Синтаксис https://github.com/google/re2/wiki/Syntax

// https://cdnb.artstation.com/p/assets/marmosets/images/020/218/089/medium/gioele-minigher-mview-image20190827-5360-ufq3v3.jpg?1566891090
// https://cdnb.artstation.com/p/assets/marmosets/attachments/020/218/089/original/toolbag_upload.mview?1566891090

var Сlient *mongo.Client    // экземпляр соединения с монгой
var DBName string = "test"  // название базы данных
var CollName string = "col" // название кооолекции документов

// TODO: сделать уже наконец получение параметров через коммандную строку
var Iterate int
var RunMode bool

// TODO: необходимо ли использование доступа в этом файле к монге? Или делать по нормальному?
// соединение с Монгой
// TODO: доделать фильтровать выборку можно только с помошью структуры

//var result []struct{ Text string `bson:"text"` }
//err := c.Find(nil).Select(bson.M{"text": 1}).All(&result)
//if err != nil {
//// handle error
//}
//for _, v := range result {
//fmt.Println(v.Text)
//}

func mong() {

	// коллекция документов с которой будет проходить работа
	var collection = Сlient.Database(DBName).Collection(CollName)
	// var result bson.M

	// TODO: Делаем запрос к базе данных и берем любой первый документ, проверяем есть ли у него поле с массивом ссылок
	//  на файлы со сценами если данной информации нет то передаем ссылку на страницу парсеру который соберет ссылки
	//  на внедренные сцены

	// TODO: Разделять функционал или нет?

	// TODO: делаем запрос к коллекции и вытаскиваем URLы к embedded сценам
	// URL к ним передаем в функцию скачивалку

	// 	err := collection.FindOne(context.TODO(), bson.D{{"id", 5946729}}).Decode(&result)
	//

	//*********

	// find all documents in which the "name" field is "Bob"
	// specify the Sort option to sort the returned documents by age in ascending order
	opts := options.Find().SetSort(bson.D{{"user_id", 1}})
	cursor, err := collection.Find(context.TODO(), bson.D{{"id", 5946729}}, opts)
	if err != nil {
		log.Fatal(err)
	}

	// get a list of all returned documents and print them out
	// see the mongo.Cursor documentation for more examples of using cursors
	var results []bson.M

	fmt.Println("		mong -> печать курсора выборки")
	fmt.Println(cursor)

	if err = cursor.All(context.TODO(), &results); err != nil {

		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Println("				печать экземпляра результатов мапа")
		fmt.Println(result)

		fmt.Println("_____")
		fmt.Println("* * * * * * * * * * * * * *")
	}

	// конвертирование bson в структуру

	type Struct1 URLs

	fmt.Println("*** *** *** *** *** ***")

	// var m bson.M
	var s Project

	// нужно проверять на то сколько элементов возвращено в ответе, если их нет то на таком будет Fatal
	bsonBytes, _ := bson.Marshal(results[0])

	bson.Unmarshal(bsonBytes, &s)

	fmt.Println("выводим на экран ")
/*	fmt.Println(s.Icons.Model3d)*/
	fmt.Println(s.URLs)

	//*********

	fmt.Println(err)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection

		if err == mongo.ErrNoDocuments {
			// если нет совпадений и документ не найден
			fmt.Println("документ не найден")
			return
		}
		log.Fatal(err)
	}

	// insertResult, err := collection.InsertOne(context.TODO(), dat.Data[index])

}

// функция вытаcкивания имени файла из урл
// TODO: доделать
//  1. определиться с путем к файлу, можно наверное по айди к ембеддед вытащить регуляркой из URL сохраненого в MongoDB
// 																		https://www.artstation.com/embed/20218089
//
func extractor(fileURL string) (name string) {

	fmt.Println("			EXTRACTOR ->")

	re := regexp.MustCompile(`(?P<url>.*/)(?P<a>.*\.\w{5})`)
	match := re.FindStringSubmatch(fileURL)
	// match := re.FindStringSubmatch("https://cdnb.artstation.com/p/assets/marmosets/attachments/020/218/089/original/toolbag_upload.mview?1566891090")

	fmt.Println(match[2])
	return match[2]
}

// убирает числовой айди за именем файла, кторорый возможно служит для учета доступа к определенным файлам
func urlCrop(fileURL string) (name string) {

	fmt.Println("			EXTRACTOR ->")

	// старая регуляррка
	// re := regexp.MustCompile(`(?P<url>https://)(?P<path>\w.+/)(?P<name>\w+.\w{5})`)

	// новая регулярка имя файла во второй группе
	// (?P<url>.*/)(?P<a>.*\.\w{5})

	re := regexp.MustCompile(`(?P<url>.*/)(?P<a>.*\.\w{5})`)

	match := re.FindStringSubmatch(fileURL)
	// match := re.FindStringSubmatch("https://cdnb.artstation.com/p/assets/marmosets/attachments/020/218/089/original/toolbag_upload.mview?1566891090")

	fmt.Println(match[2])
	return match[1] + match[2]
}

func extrsctEmbeddedId(embeddedPath string) (embdId string) {
	re := regexp.MustCompile(`(?P<name>\d+$)`)
	match := re.FindStringSubmatch(embeddedPath)
	// src="https://www.artstation.com/embed/16064327
	return match[0]
}

// функция скачивания файла по урл
//
func downloadFile(filep, filename, url string) (err error) {

	// create the path folders
	_ = os.MkdirAll(filep, os.ModeDir)

	// Create the file
	out, err := os.Create(filepath.Join(".", filep, filename))
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	fmt.Println("записываем файл на диск")
	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("функция скачивания файла успешно завершена")
	return nil
}

//***
// Часть получения данных через Chromedev
// получает URL страницы с работой и возвращает байтовый массив с объектами в которых ссылки
//***
// TODO: Сделать уже наконец все в одном контексте
func urlsFetch(urlstr string, Ctx context.Context) (Outer []byte /*map[string]interface{}*/, err error) {

	// create context
	// ctx, cancel := chromedp.NewContext(context.Background())
	// defer cancel()

	var tsk = chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(map[string]interface{}{
			"X-Header":                  "my request header",
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
			"Accept-Encoding":           "gzip, deflate, br",
			"Accept-Language":           "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
			"Connection":                "keep-alive",
			"Dnt":                       "1",
			"Sec-Fetch-Mode":            "navigate",
			"Sec-Fetch-Site":            "none",
			"Upgrade-Insecure-Requests": "1",
			"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36",
		})),

		chromedp.ResetViewport(),
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second),

		// выполняем код яваскрипт так, как будто мы его выполняем в консоли браузера, выведеное значение присваиваем переменной
		chromedp.EvaluateAsDevTools(`var result = {"data":[]};
console.log(result);
qw = document.querySelectorAll(
    "body > div.wrapper > div.wrapper-main > div > div.gallery-single.form-prolongate-session > main.artwork-container > project-assets.project-assets > div.artwork.marmoset > div.artwork-data > div.artwork-embedded > iframe"
);

qw.forEach(
  function (item, index, array) {
    result.data.push(
      {"pageURL": item.src,"fileURL": item.contentWindow.currentViewer.sceneURL}
      )
  }
  ); 
    result;`, &Outer),
		// TODO: виснет хром, может так поможет?
		//  может headful режим что нибудь даст?
		chromedp.Reload(),
		chromedp.Stop(),
	}

	fmt.Println("выполняем задания chromedp")
	fmt.Println("")
	// capture text element
	if err := chromedp.Run(Ctx, tsk); err != nil {

		// TODO: виснет хром, может так поможет?
		//  может headful режим что нибудь даст?
		return nil, err
	}

	chromedp.Stop()
	// TODO: виснет хром, может так поможет?
	//  может headful режим что нибудь даст?
	return Outer, nil

}

// функция которая ходит по ссылкам и запускает скачивание
func IteratorURL(result Project, collection *mongo.Collection) (err error) {

	fmt.Println("")
	fmt.Println("запущен итератор")

	for _, url := range result.URLs {
		fmt.Println(url.FileURL)
		fmt.Println(url.PageURL)

		var t = url.FileURL

		fmt.Println(t)

		filep := filepath.Join(".", strconv.Itoa(result.User_id), strconv.Itoa(result.Id))

		// TODO: заменить имя файлы на айди ембеддед из ссылки на страницу пока сделано айди плюс оригинальное имя
		var filename = extrsctEmbeddedId(url.PageURL) + "-" + extractor(url.FileURL)

		fmt.Println("")
		fmt.Println("			-> Скачаем файл " + t)
		fmt.Println("")

		err := downloadFile(filep, filename, t)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	fmt.Println("подготавливаем запрос в базу данных для записи информации об успешно скачаном файле")

	filter := bson.D{{"id", result.Id}}
	update := bson.D{
		{
			"$set", bson.D{{"checked", true},},
		},
	}

	fmt.Println(" записываем данные об успешно скачаном файле")
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Найдено %v документов и обновлено %v документов \n", updateResult.MatchedCount, updateResult.ModifiedCount)

	fmt.Println("итератор выполнен")
	return nil
}

func get(Ctx context.Context) {

	Iterate -= 1

	var result Project

	// коллекция документов с которой будет проходить работа
	var collection = Сlient.Database(DBName).Collection(CollName)
	err := collection.FindOne(context.TODO(), bson.D{{"checked", false}}).Decode(&result)

	println("декодирование")


	if err != nil {
		if err == mongo.ErrNoDocuments {
			// если нет совпадений и документ не найден
			// TODO: делать
			//  получается хуйня если на странице нет EMBEDDED

			fmt.Println("")
			fmt.Println("		--> пришел ответ что документ с checked=false не найден")
			fmt.Println("		--> берем из базы запись с отсутствующим checked")
			err = collection.FindOne(context.TODO(), bson.D{{"checked", nil}}).Decode(&result)

			if err != nil {

				fmt.Println("		--> ответ содержал ошибку:")
				fmt.Println(err)
				fmt.Println("		--> выходим из условия:")
				return
			}

			fmt.Println("		--> откроем страницу " + result.Permalink + " в браузере")
			var fetchRes, err = urlsFetch(result.Permalink, Ctx)
			if err != nil {
				fmt.Println("	 		--> пришла ошибка")
				fmt.Println(err)
			}

			fmt.Println("			--> обработаем данные")
				dat := urlarr{}
			if err := json.Unmarshal(fetchRes, &dat); err != nil {
				panic(err)
			}

			//fmt.Println(dat.Data[0].FileURL)

			//fmt.Println(fetchRes)

			//
			//if err := json.Unmarshal(fetchRes.Bytes(), &bb); err != nil {
			//	panic(err)
			//}

			fmt.Println("			--> подготовим запрос к MongoDB")
			filter := bson.D{{"id", result.Id}}
			update := bson.D{
				{
					"$set", bson.D{{"urls", dat.Data},},
				},
				{
					"$set", bson.D{{"checked", false}},
				},
			}

			fmt.Println("			--> выполним запрос к MongoDB и обновим данные")
			updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Найдено %v документов и обновлено %v документов \n записаны ссылки на страницы и файлы", updateResult.MatchedCount, updateResult.ModifiedCount)
			fmt.Println("		<--")
			// Iterate =true
		} else {
			println("декодирование пошло по пизде")
			fmt.Println("		--> ААА!!! что-то непонятное")
			// выпадаем в фатал
			log.Fatal(err)

		}
		return
	}

	fmt.Printf("найден один документ: %+v\n обработаем и будем скачивать файлы", result)
	// передаем результат поиска в функцию обработки перед закачкой
	err = IteratorURL(result, collection)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("		<--")
	fmt.Println("")
	// Iterate = true
}

func init() {
	// инициализирован ОРМ для монги
	Сlient, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	// Создаем соединение
	err := Сlient.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	// Проверка соединения
	err = Сlient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

func main() {

	// create context
	var Ctx, Cancel = chromedp.NewContext(context.Background())

	defer Cancel()

	// добавил получение колитества итераций через аттрибуты коммандной строки
	flag.IntVar(&Iterate, "int", 5, "an int")
	flag.BoolVar(&RunMode, "runmode", false, "an bool")

	var iterate = Iterate
	var runmode = RunMode

	flag.Parse()

	//var prj Project

	fmt.Println("Downloading microservice started-> ")

	//get()

	// -> берем из базы любой документ у которого { Check: false } и URLS.lenght > 0
	// 		записываем данные в переменную в формате структуры Project, берем:
	// 		prj.Id
	// 		prj.User_id
	// 		prj.URLs[]
	// 			начинаем перебирать циклом все prj.URLs[x] и брать fileURL pageURL и передавать функции
	//
	// 			скачивания файлов (prj.User_id, prj.Id, prj.URLs[], fileURL, pageURL) путь до файла
	//			./ + prj.User_id + / + prj.Id + / + prj.URLs[index].pageURL + / + prj.URLs[index].fileURL
	//
	// 			если функция отработает и не вернет ошибку будем считать что файл записан и ссылка отработана
	// 			если во время цикла не происходило ошибок, записываем в документ { Check: true }
	// 			так же увеличиваем переменную счетчик полученных файлов
	// 	если условие не выполняется и документ у которого Check != true и URLS.lenght > 0 не найден
	// то берем любой первый попавшийся документ у которого поле Check отсутствует (db.inventory.find( { item: null } ))
	// и URLS.lenght отсутствует и запускаем функцию краулинга
	// дополнительных данных fileURL pageURL с записью в базу данных , и проставляем Check: false
	// если и таких данных не найдено то пишем сообщение и завершаем программу

	// ->

	for Iterate > 0 {

		fmt.Println(strconv.Itoa(Iterate))
		fmt.Println("****************************")
		fmt.Println("")
		fmt.Println("запрашиваем информацию в базе данных")

		if RunMode {

			fmt.Println("RunMOde")

			Iterate -= 1
			var result Project
			// коллекция документов с которой будет проходить работа
			var collection = Сlient.Database(DBName).Collection(CollName)
			err := collection.FindOne(context.TODO(), bson.D{{"checked", nil}}).Decode(&result)

			if err == nil {

				// если есть совпадения
				fmt.Println("")
				fmt.Println("		--> пришел ответ что документ с checked=nil найден")
				fmt.Println("		--> берем из базы эту запись")
				fmt.Println("		--> откроем страницу " + result.Permalink + " в браузере")

				var fetchRes, err = urlsFetch(result.Permalink, Ctx)
				if err != nil {
					fmt.Println("	 		--> пришла ошибка " + result.Permalink)
					fmt.Println(err)
				}

				fmt.Println("			--> обработаем данные")
				dat := urlarr{}
				if err := json.Unmarshal(fetchRes, &dat); err != nil {
					panic(err)
				}

				fmt.Println("			--> подготовим запрос к MongoDB")
				filter := bson.D{{"id", result.Id}}
				update := bson.D{
					{
						"$set", bson.D{{"urls", dat.Data},},
					},
					{
						"$set", bson.D{{"checked", false}},
					},
				}

				fmt.Println("			--> выполним запрос к MongoDB и обновим данные")
				updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Найдено %v документов и обновлено %v документов \n записаны ссылки на страницы и файлы", updateResult.MatchedCount, updateResult.ModifiedCount)
				fmt.Println("		<--")
			} else {
				fmt.Println("ошибка")
				return
			}

		} else {
			get(Ctx)
		}

	}

	fmt.Println(`/c start f:/PARSE/main.exe -int=` + strconv.Itoa(iterate) + ` -runmode=` + strconv.FormatBool(!runmode))
	cmnd := exec.Command("C:\\WINDOWS\\system32\\cmd.exe", "/c start f:/PARSE/main.exe -int="+strconv.Itoa(iterate)+" "+"-runmode="+strconv.FormatBool(!runmode))
	cmnd.Start()
	fmt.Println("НИЧЕГО НЕ НАШЕЛ И ЗАВЕРШИЛСЯ И ЗАПУСТИЛ ЕЩЕ РАЗ")

}
