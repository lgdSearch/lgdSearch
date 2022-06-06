package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/dgraph-io/badger/v3"
	"github.com/nfnt/resize"
	"github.com/wangbin/jiebago"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"lgdSearch/pkg/db/badgerStorage"
	"lgdSearch/pkg/db/boltStorage"
	"lgdSearch/pkg/models"
	"lgdSearch/pkg/pagination"
	"lgdSearch/pkg/trie"
	"lgdSearch/pkg/utils"
	"lgdSearch/pkg/utils/colf/doc"
	"lgdSearch/pkg/utils/colf/keyIds"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	//Shard is the default engine shard
	//This means that the engine will have shard AVL tree and leveldb
	Shard = 100

	BadgerShard = 10

	BoltBucketSize = 100
)

type Option struct {
	PictureName    string
	BoltKeyIdsName string
	BadgerDocName  string
	DictionaryName string
	DocIndexName   string
}

type Engine struct {
	IndexPath string

	Option *Option

	//关键字和Id映射
	KeyIdsStorage   *boltStorage.BoltdbStorage
	boltBucketNames [][]byte

	//文档仓
	DocStorages []*badgerStorage.BadgerStorage

	//锁
	sync.Mutex
	//等待
	sync.WaitGroup

	//文件分片
	Shard int

	//添加索引的通道
	AddDocumentWorkerChan chan models.IndexDoc

	// index 是data的id索引，可以理解为 doc 的个数
	index uint32

	//是否调试模式
	IsDebug bool
}

var seg jiebago.Segmenter
var SearchEngine *Engine

func Set(e *Engine) {
	SearchEngine = e
}

func (e *Engine) Init() {
	e.Add(1)
	defer e.Done()
	//线程数=cpu数
	runtime.GOMAXPROCS(runtime.NumCPU())
	//保持和gin一致
	e.IsDebug = os.Getenv("GIN_MODE") != "release"
	if e.Option == nil {
		e.Option = e.GetOptions()
	}
	log.Println("数据存储目录：", e.IndexPath)

	err := seg.LoadDictionary(e.getFilePath(e.Option.DictionaryName))
	if err != nil {
		panic("dictionary not find")
	}

	if e.Shard == 0 {
		e.Shard = Shard
	}

	utils.Read(&e.index, e.getFilePath(e.Option.DocIndexName))
	if e.index == 0 {
		e.index = 50000
	}

	//初始化chan
	e.AddDocumentWorkerChan = make(chan models.IndexDoc, 1000)

	filterMapInit() // 初始化分词过滤

	e.KeyIdsStorage, err = boltStorage.Open(e.getFilePath(e.Option.BoltKeyIdsName))
	for i := 0; i < BadgerShard; i++ {
		option := badger.DefaultOptions(e.getFilePath(fmt.Sprintf("%s_%d.db", e.Option.BadgerDocName, i)))
		option.Logger = nil
		badgerdb := badgerStorage.Open(option)
		e.DocStorages = append(e.DocStorages, badgerdb)
	}

	for i := 0; i < BoltBucketSize; i++ {
		bucketName := utils.Uint32ToBytes(uint32(i))
		e.boltBucketNames = append(e.boltBucketNames, bucketName)
		err := e.KeyIdsStorage.CreateBucketIfNotExist(bucketName)
		if err != nil {
			return
		}
	}

	//初始化文件存储
	go e.DocumentWorkerExec()
	for shard := 0; shard < e.Shard; shard++ {
		// 初始化 pictureMap
		picture := make(map[uint32]interface{})
		utils.Read(&picture, e.getFilePath(fmt.Sprintf("%s_%d.txt", e.Option.PictureName, shard)))
		pictureUrlMap = append(pictureUrlMap, picture)
	}
	log.Println("初始化完成")

	//初始化完成，自动检测索引并持久化到磁盘
	go e.automaticFlush()
}

func (e *Engine) IndexDocument(doc models.IndexDoc) {
	e.AddDocumentWorkerChan <- doc
}

// DocumentWorkerExec 添加文档队列
func (e *Engine) DocumentWorkerExec() {
	for {
		docs := <-e.AddDocumentWorkerChan
		e.AddDocument(&docs)
	}
}

func (e *Engine) GetOptions() *Option {
	return &Option{
		PictureName:    "picture/pic",
		DictionaryName: "dictionary.txt",
		DocIndexName:   "dataIndex.txt",
		BoltKeyIdsName: "bolt_keyIds.db",
		BadgerDocName:  "badger_doc",
	}
}

// Get the key corresponding shard
func (e *Engine) getShard(id uint32) int {
	return int(id % uint32(e.Shard))
}

func (e *Engine) InitOption(option *Option) {
	if option == nil {
		//默认值
		option = e.GetOptions()
	}
	e.Option = option

	//初始化其他的
	e.Init()
}

func (e *Engine) getFilePath(fileName string) string {
	return e.IndexPath + string(os.PathSeparator) + fileName
}

// AddDocument
func (e *Engine) AddDocument(index *models.IndexDoc) {
	//等待初始化完成
	e.Wait()

	text := index.Text

	// wordLen 是带重复数据的长度
	keys, wordMap := e.WordCutFilter(text)

	for _, key := range keys {
		data := keyIds.StorageId{
			Id:    index.Id,
			Score: wordMap[key],
		}

		keyBuf := utils.Uint32ToBytes(key)
		bucketName := e.boltBucketNames[key%BoltBucketSize]
		buf, found := e.KeyIdsStorage.Get(keyBuf, bucketName)

		ids := new(keyIds.StorageIds)
		if found {
			ids.UnmarshalBinary(buf)
			ids.StorageIds = append(ids.StorageIds, &data)
		} else {
			ids.StorageIds = append(ids.StorageIds, &data)
		}

		bufs, _ := ids.MarshalBinary()
		e.KeyIdsStorage.Set(keyBuf, bufs, bucketName)
	}
	e.addDoc(index)
}

func (e *Engine) addDoc(index *models.IndexDoc) {
	k := utils.Uint32ToBytes(index.Id)

	value := doc.StorageIndexDoc{Text: index.Text, Url: index.Url}
	buf, err := value.MarshalBinary()
	if err != nil {
		log.Println("addDoc", index.Id, err)
	}
	err = e.DocStorages[index.Id%BadgerShard].Set(k, buf)
	if err != nil {
		log.Println("doc set", index.Id, err)
	}
}

var wordFilterMap map[string]interface{}

func filterMapInit() {
	wordFilterMap = make(map[string]interface{})
	wordFilterMap["了"] = nil
	wordFilterMap["的"] = nil
	wordFilterMap["么"] = nil
	wordFilterMap["呢"] = nil
	wordFilterMap["和"] = nil
	wordFilterMap["与"] = nil
	wordFilterMap["于"] = nil
	wordFilterMap["吗"] = nil
	wordFilterMap["吧"] = nil
	wordFilterMap["呀"] = nil
	wordFilterMap["啊"] = nil
	wordFilterMap["哎"] = nil
	wordFilterMap["是"] = nil
	wordFilterMap["人"] = nil
	wordFilterMap["名"] = nil
	wordFilterMap["在"] = nil
	wordFilterMap["不"] = nil
	wordFilterMap["被"] = nil
	wordFilterMap["有"] = nil
	wordFilterMap["无"] = nil
	wordFilterMap["都"] = nil
	wordFilterMap["也"] = nil
	wordFilterMap["【"] = nil
	wordFilterMap["】"] = nil
	wordFilterMap["《"] = nil
	wordFilterMap["》"] = nil
	wordFilterMap["，"] = nil
	wordFilterMap["。"] = nil
	wordFilterMap["？"] = nil
	wordFilterMap["！"] = nil
	wordFilterMap["、"] = nil
	wordFilterMap["；"] = nil
	wordFilterMap["："] = nil
	wordFilterMap["（"] = nil
	wordFilterMap["）"] = nil
}

func (e *Engine) WordCutFilter(text string) ([]uint32, map[uint32]float32) {
	//不区分大小写
	text = strings.ToLower(text)

	// wordMap is to save word, keyMap is to save hash(word)
	var keyMap = make(map[uint32]float32)
	var wordMap = make(map[uint32]interface{})
	words := make([]uint32, 0)

	resultChan := seg.CutForSearch(text, true)
	pre := true
	for {
		w, ok := <-resultChan
		if !ok {
			break
		}

		switch len(w) { //  过滤分词
		case 1:
			{
				if (w[0] > 47 && w[0] < 58) || (w[0] > 64 && w[0] < 91) || (w[0] > 96 && w[0] < 123) {
					// 保留数字和字母
				} else if w[0] <= 32 || w[0] == 127 || !pre { // 不要空格和分隔符, 不常用字符和符号只保留首位
					pre = true
					continue
				}
			}
		case 3:
			{
				_, ok := wordFilterMap[w]
				if ok && !pre {
					continue
				}
			}
		}

		value := utils.StringToInt(w)
		words = append(words, value)

		wordMap[value] = nil
	}
	lenWords := float32(len(words))
	for index, val := range words { // 越前面的比重越大
		// math.log10(10 + index) 是为了平衡 index=10,len=10 和 index=5000,len=10000的情况
		// 这种情况下 f(10, 10) > f(5000, 10000), etc...
		// f(10, 10000) > f(10, 10) > f(5000, 10000)
		keyMap[val] += float32(2*len(words)-index+1) / (float32(math.Log10(float64(10+index))) * lenWords)
	}

	var keysSlice []uint32 = make([]uint32, len(keyMap))

	index := 0
	for word, _ := range wordMap {
		keysSlice[index] = word
		index += 1
	}

	return keysSlice, keyMap
}

// WordCut 分词，只取长度大于2的词 | 是假的，会取到长度为1的词, 需要取到长度为1的词
// filter 是关键词过滤
// 这里输入的请求长度有限制,最多30,那么分词不会很多,同时考虑对于关键词过滤是否也做限制 10 个
func (e *Engine) WordCut(text string, filter []string) []string {
	//不区分大小写
	text = strings.ToLower(text)

	// wordMap is to save word, keyMap is to save hash(word)
	var wordMap = make(map[string]int)

	resultChan := seg.CutForSearch(text, true)
	pre := true
	for {
		w, ok := <-resultChan
		if !ok {
			break
		}
		if filter != nil { // 先在分词中过滤一遍
			for _, val := range filter { // 最大 O(20 * 10)
				if w == val {
					continue
				}
			}
		}

		switch len(w) { // 过滤分词
		case 1:
			{
				if (w[0] > 47 && w[0] < 58) || (w[0] > 64 && w[0] < 91) || (w[0] > 96 && w[0] < 123) {
					// 保留数字和字母
				} else if w[0] <= 32 || w[0] == 127 || !pre { // 不要空格和分隔符, 不常用字符和符号只保留首位
					pre = true // 对于这些分隔符我认为是一句新的话的开头
					continue
				}
			}
		case 3:
			{
				_, ok := wordFilterMap[w]
				if ok && !pre {
					continue
				}
			}
		}
		wordMap[w]++
		pre = false
	}

	var wordsSlice = make([]string, len(wordMap))

	index := 0
	for word, _ := range wordMap {
		wordsSlice[index] = word
		index += 1
	}

	return wordsSlice
}

//MultiSearch 多线程搜索
func (e *Engine) MultiSearch(request *models.SearchRequest) *models.SearchResult {
	//等待搜索初始化完成
	if e.IsDebug {
		log.Println("Search start")
	}

	e.Wait()
	//分词搜索
	words := e.WordCut(request.Query, request.FilterWord)
	if e.IsDebug {
		log.Println("分词数: ", len(words))
	}

	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(words))

	var fastSort pagination.FastSort
	_time := utils.ExecTime(func() {
		var allValues = make([]*models.SliceItem, 0)

		for _, word := range words {
			go e.SimpleSearch(word, utils.StringToInt(word), func(values []*models.SliceItem) {
				lock.Lock()
				allValues = append(allValues, values...)
				lock.Unlock()
				wg.Done()
			})
		}

		wg.Wait()
		fastSort = pagination.FastSort{ScoreMap: make(map[uint32]float32, int(float32(len(allValues))*1.5))} // 预分配空间，优化效率
		fastSort.Add(allValues)
	})
	if e.IsDebug {
		log.Println("搜索时间:", _time, "ms")
	}
	// 处理分页
	request = request.GetAndSetDefault()

	//读取文档
	var result = &models.SearchResult{
		Total:     fastSort.Count(),
		Time:      float32(_time),
		Page:      request.Page,
		Limit:     request.Limit,
		Words:     words,
		Documents: nil,
		Related:   trie.Tree.Search([]rune(request.Query)),
	}

	tim := time.Now()
	trie.BuildTree(request.FilterWord)

	// 设置fail指针
	trie.SetNodeFailPoint()

	if e.IsDebug {
		log.Println("ACTrie init Success! 耗时: ", time.Since(tim))
	}

	_time = utils.ExecTime(func() {

		pager := new(pagination.Pagination)
		var resultIds []models.SliceItem
		_tt := utils.ExecTime(func() {
			resultIds = fastSort.GetAll()
		})

		if e.IsDebug {
			log.Println("处理排序耗时", _tt, "ms")
			log.Println("结果集大小", len(resultIds))
		}

		pager.Init(request.Limit, len(resultIds))
		//设置总页数
		result.PageCount = pager.PageCount

		//读取单页的id
		if pager.PageCount != 0 {

			start, end := pager.GetPage(request.Page)
			log.Println("start: ", start, " --- end: ", end)
			if start == -1 {
				return
			}
			items := resultIds[start : end+1]
			if e.IsDebug {
				log.Println("Page: ", "start ", start, "end ", end)
			}

			result.Documents = make([]models.ResponseDoc, len(items))

			var wg sync.WaitGroup
			wg.Add(len(items))

			//只读取前面 limit 个
			_tt := time.Now()
			for index, item := range items {
				go func(index int, item models.SliceItem) {
					defer wg.Done()

					_cost := time.Now()
					buf := e.GetDocById(item.Id)

					if e.IsDebug {
						log.Println("Id: ", item.Id, "--- GetDocById: ", time.Since(_cost), "--- 数据长度: ", len(buf))
					}

					storageDoc := new(doc.StorageIndexDoc)
					if buf != nil {
						storageDoc.UnmarshalBinary(buf)

						// 查找 ACTrie, 如果有过滤词, 过滤掉
						if request.FilterWord != nil &&
							len(request.FilterWord) > 0 &&
							trie.AcAutoMatch(storageDoc.Text) {
							return
						}
					}
					result.Documents[index].Score = item.Score

					if buf != nil {
						result.Documents[index].Id = item.Id
						result.Documents[index].Url = storageDoc.Url
						text := storageDoc.Text
						//处理关键词高亮
						highlight := request.Highlight
						if highlight != nil {
							//全部小写
							text = strings.ToLower(text)
							for _, word := range words {
								text = strings.ReplaceAll(
									text, word,
									fmt.Sprintf("%s%s%s", highlight.PreTag, word, highlight.PostTag))
							}
						}
						result.Documents[index].Text = text
					}
				}(index, item)
			}

			wg.Wait()

			if e.IsDebug {
				log.Println("分页耗时: ", time.Since(_tt))
			}

		}
	})
	if e.IsDebug {
		log.Println("处理数据耗时：", _time, "ms")
	}

	return result
}

//SimpleSearch key is one of keys, keys -> 当前查询语句的所有分词
func (e *Engine) SimpleSearch(word string, key uint32, call func(ranks []*models.SliceItem)) {
	_time := time.Now()

	s := e.KeyIdsStorage

	kv := utils.Uint32ToBytes(key)

	data, find := s.Get(kv, e.boltBucketNames[key%BoltBucketSize]) // key.ids

	if find {
		array := new(keyIds.StorageIds)
		array.UnmarshalBinary(data)

		results := make([]*models.SliceItem, len(array.StorageIds))
		if e.IsDebug {
			log.Println("读数据时间: ", time.Since(_time), "--- word: ", word, "--- key: ", key, "--- Ids长度:", len(array.StorageIds))
		}

		for index, id := range array.StorageIds { // 遍历 ids
			rank := &models.SliceItem{}
			rank.Id = id.Id
			rank.Score = float32(math.Log10(float64(e.index)/float64(len(array.StorageIds)+1))) * id.Score
			results[index] = rank
		}
		call(results)
	} else {
		call(nil)
	}
}

// ------------------------------------- wordToSearchPicture

func (e *Engine) MultiSearchPicture(request *models.SearchRequest) *models.SearchPictureResult {
	//等待搜索初始化完成

	if e.IsDebug {
		log.Println("Search start")
	}

	e.Wait()
	//分词搜索
	words := e.WordCut(request.Query, request.FilterWord)
	if e.IsDebug {
		log.Println("分词数: ", len(words))
	}

	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(words))

	var fastSort pagination.FastSort
	_time := utils.ExecTime(func() {
		var allValues = make([]*models.SliceItem, 0)

		for _, word := range words {
			go e.SimpleSearch(word, utils.StringToInt(word), func(values []*models.SliceItem) {
				lock.Lock()
				allValues = append(allValues, values...)
				lock.Unlock()
				wg.Done()
			})
		}

		wg.Wait()
		fastSort = pagination.FastSort{ScoreMap: make(map[uint32]float32, int(float32(len(allValues))*1.5))} // 预分配空间，优化效率
		fastSort.Add(allValues)
	})
	if e.IsDebug {
		log.Println("搜索时间:", _time, "ms")
	}
	// 处理分页
	request = request.GetAndSetDefault()

	//读取文档
	var result = &models.SearchPictureResult{
		Total:     fastSort.Count(),
		Time:      float32(_time),
		Page:      request.Page,
		Limit:     request.Limit,
		Words:     words,
		Documents: nil,
	}

	_time = utils.ExecTime(func() {

		pager := new(pagination.Pagination)
		var resultIds []models.SliceItem
		_tt := utils.ExecTime(func() {
			resultIds = fastSort.GetAll()
		})

		if e.IsDebug {
			log.Println("处理排序耗时", _tt, "ms")
			log.Println("结果集大小", len(resultIds))
		}

		pager.Init(request.Limit, len(resultIds))
		//设置总页数
		result.PageCount = pager.PageCount

		//读取单页的id
		if pager.PageCount != 0 {

			start, end := pager.GetPage(request.Page)
			if start == -1 {
				return
			}

			items := resultIds[start : end+1]
			if e.IsDebug {
				log.Println("Page: ", "start ", start, "end ", end)
			}

			result.Documents = make([]models.ResponseUrl, len(items))
			var wg sync.WaitGroup
			wg.Add(len(items)) // 并发上传图片

			//只读取前面 limit 个
			_tt := time.Now()
			for index, item := range items {
				go func(index int, item models.SliceItem) {
					defer wg.Done()
					_cost := time.Now()
					buf := e.GetDocById(item.Id)

					if e.IsDebug {
						log.Println("Id: ", item.Id, "--- GetDocById: ", time.Since(_cost), "--- 数据长度: ", len(buf))
					}

					result.Documents[index].Score = item.Score

					if buf != nil {
						storageDoc := new(doc.StorageIndexDoc)
						storageDoc.UnmarshalBinary(buf)

						result.Documents[index].Id = item.Id
						result.Documents[index].Url = storageDoc.Url
						e.GetPictureUrl(storageDoc.Url, item.Id, &result.Documents[index].ThumbnailUrl, func() {
						}) // get url
						text := storageDoc.Text
						//处理关键词高亮
						highlight := request.Highlight
						if highlight != nil {
							//全部小写
							text = strings.ToLower(text)
							for _, word := range words {
								text = strings.ReplaceAll(text, word, fmt.Sprintf("%s%s%s", highlight.PreTag, word, highlight.PostTag))
							}
						}
						result.Documents[index].Text = text
					}
				}(index, item)
			}
			wg.Wait() // 等待并发结束

			if e.IsDebug {
				log.Println("分页耗时: ", time.Since(_tt))
			}
		}
	})
	if e.IsDebug {
		log.Println("处理数据耗时：", _time, "ms")
	}

	return result
}

func pictureSearchHandleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

func pictureSearchGetRemote(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		// 如果有错误返回错误内容
		return nil, err
	}
	// 使用完成后要关闭，不然会占用内存
	defer res.Body.Close()
	// 读取字节流
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bytes, err
}
func pictureCompress(buf []byte) ([]byte, error) {
	//文件压缩
	decodeBuf, layout, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	// 修改图片的大小
	set := resize.Resize(0, 200, decodeBuf, resize.Lanczos3)
	NewBuf := bytes.Buffer{}
	switch layout {
	case "png":
		err = png.Encode(&NewBuf, set)
	case "jpeg", "jpg":
		err = jpeg.Encode(&NewBuf, set, &jpeg.Options{Quality: 80})
	default:
		return nil, errors.New("该图片格式不支持压缩")
	}
	if err != nil {
		return nil, err
	}
	if NewBuf.Len() < len(buf) {
		buf = NewBuf.Bytes()
	}
	return buf, nil
}

// 判断服务器上有没有这张图片的缩略图
var pictureUrlMap []map[uint32]interface{}

func (e *Engine) putInOSS(url string, id uint32) bool {
	resByte, err := pictureSearchGetRemote(url)
	if err != nil {
		fmt.Println(err)
	}
	resBytes, err := pictureCompress(resByte)
	if err != nil {
		//fmt.Println(err)
		return false
	}
	//get()
	//Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := "http://oss-cn-hangzhou.aliyuncs.com"
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	accessKeyId := "LTAI5tRRkejmZXzM5k9QvpfW"
	accessKeySecret := "TkOOCo6dTGfs6l9j4iZ8UfQJa67Qg3"
	bucketName := "lgdsearch"
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		pictureSearchHandleError(err)
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		pictureSearchHandleError(err)
	}

	reader := bytes.NewReader(resBytes)
	pictureUrl := fmt.Sprintf("example/%d.jpg", id)
	err = bucket.PutObject(pictureUrl, reader)
	if err != nil {
		pictureSearchHandleError(err)
	}
	pictureUrlMap[e.getShard(id)][id] = nil
	return true
}

func (e *Engine) GetPictureUrl(url string, id uint32, thumbnailUrl *string, call func()) {
	_, ok := pictureUrlMap[e.getShard(id)][id]
	if !ok {
		// 压缩图片到服务器
		ok = e.putInOSS(url, id)
	}
	if ok {
		*thumbnailUrl = fmt.Sprintf("https://lgdsearch.oss-cn-hangzhou.aliyuncs.com/example/%d.jpg", id)
	} else {
		*thumbnailUrl = url
	}
	call()
}

// ------------------------------------------------------ pictureSearch

func (e *Engine) flushDataIndex() {
	utils.Write(&e.index, "./pkg/data/dataIndex.txt")
}

func (e *Engine) getPictureSize() int {
	size := 0
	for _, urlMap := range pictureUrlMap {
		size += len(urlMap)
	}
	return size
}

func (e *Engine) flushPictureMap() {
	for index, urlMap := range pictureUrlMap {
		utils.Write(&urlMap, e.getFilePath(fmt.Sprintf("%s_%d.txt", e.Option.PictureName, index)))
	}
}

// 自动保存索引，120秒钟检测一次
func (e *Engine) automaticFlush() {
	ticker := time.NewTicker(time.Second * 120)
	index := e.index
	pictureSize := e.getPictureSize()

	for {
		<-ticker.C
		if index != e.index {
			index = e.index
			e.flushDataIndex()
		}
		if pictureSize != e.getPictureSize() {
			pictureSize = e.getPictureSize()
			e.flushPictureMap()
		}
		//定时GC
		runtime.GC()
	}

}

// GetDocById 通过id获取文档
func (e *Engine) GetDocById(id uint32) []byte {
	key := utils.Uint32ToBytes(id)
	buf, found := e.DocStorages[id%BadgerShard].Get(key)
	if found {
		return buf
	}

	return nil
}

// Close will save AVL data
func (e *Engine) Close() {
	e.Lock()
	defer e.Unlock()

	//保存文件
	e.flushDataIndex()

	for _, db := range e.DocStorages {
		db.Close()
	}
	e.KeyIdsStorage.Close()
}
