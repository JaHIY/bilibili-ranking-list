package main

import (
    "encoding/json"
    "fmt"
    "github.com/JaHIY/bilibili-ranking-list/models"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "time"
)

type BilibiliAnimeList struct {
    Code int `json:"code"`
    Message string `json:"message"`
    Data BilibiliAnimeListData `json:"data"`
}

type BilibiliAnimeListData struct {
    HasNext int `json:"has_next"`
    Num int `json:"num"`
    Size int `json:"size"`
    Total int `json:"total"`
    List []BilibiliAnime `json:"list"`
}

type BilibiliAnime struct {
    Badge string `json:"badge"`
    BadgeType int `json:"badge_type"`
    Cover string `json:"cover"`
    IndexShow string `json:"index_show"`
    IsFinish int `json:"is_finish"`
    Link string `json:"link"`
    MediaId int `json:"media_id"`
    Order string `json:"order"`
    OrderType string `json:"order_type"`
    SeasonId int `json:"season_id"`
    Title string `json:"title"`
    TitleIcon string `json:"title_icon"`
}

type BilibiliMedia struct {
    Code int `json:"code"`
    Message string `json:"message"`
    Result BilibiliMediaResult `json:"result"`
}

type BilibiliMediaResult struct {
    Actors string `json:"actors"`
    Alias string `json:"alias"`
    Areas []BilibiliMediaResultArea `json:"areas"`
    Copyright BilibiliMediaResultCopyright `json:"copyright"`
    Cover string `json:"cover"`
    EpisodeIndex BilibiliMediaResultEpisodeIndex `json:"episode_index"`
    Evaluate string `json:"evaluate"`
    MediaId int `json:"media_id"`
    Mode int `json:"mode"`
    OriginName string `json:"origin_name"`
    Payment BilibiliMediaResultPayment `json:"payment"`
    Publish BilibiliMediaResultPublish `json:"publish"`
    Rating BilibiliMediaResultRating `json:"rating"`
    Rights BilibiliMediaResultRights `json:"rights"`
    SeasonId int `json:"season_id"`
    SeasonStatus int `json:"season_status"`
    Seasons []BilibiliMediaResultSeason `json:"seasons"`
    Staff string `json:"staff"`
    Stat BilibiliMediaResultStat `json:"stat"`
    Styles []BilibiliMediaResultStyle `json:"styles"`
    TimeLength int `json:"time_length"`
    Title string `json:"title"`
    Type int `json:"type"`
    TypeName string `json:"type_name"`
}

type BilibiliMediaResultArea struct {
    Id int `json:"id"`
    Name string `json:"name"`
}

type BilibiliMediaResultCopyright struct {
    IsFinish int `json:"is_finish"`
    IsStarted int `json:"is_started"`
}

type BilibiliMediaResultEpisodeIndex struct {
    Id int `json:"id"`
    Index string `json:"index"`
    IndexShow string `json:"index_show"`
    IsNew int `json:"is_new"`
    PlayIndexShow string `json:"play_index_show"`
}

type BilibiliMediaResultPayment struct {
    Discount int `json:"discount"`
    PayType BilibiliMediaResultPaymentPayType `json:"pay_type"`
    Price string `json:"price"`
    Promotion string `json:"promotion"`
    Tip string `json:"tip"`
    VipDiscount int `json:"vip_discount"`
    VipFirstPromotion string `json:"vip_first_promotion"`
    VipPromotion string `json:"vip_promotion"`
}

type BilibiliMediaResultPaymentPayType struct {
    AllowDiscount int `json:"allow_discount"`
    AllowPack int `json:"allow_pack"`
    AllowTicket int `json:"allow_ticket"`
    AllowTimeLimit int `json:"allow_time_limit"`
    AllowVipDiscount int `json:"allow_vip_discount"`
}

type BilibiliMediaResultPublish struct {
    IsFinish int `json:"is_finish"`
    IsStarted int `json:"is_started"`
    PubDate string `json:"pub_date"`
    PubDateShow string `json:"pub_date_show"`
    ReleaseDateShow string `json:"release_date_show"`
    TimeLengthShow string `json:"time_length_show"`
}

type BilibiliMediaResultRating struct {
    Count int `json:"count"`
    Score float32 `json:"score"`
}

type BilibiliMediaResultRights struct {
    AllowBp int `json:"allow_bp"`
    AllowBpRank int `json:"allow_bp_rank"`
    AllowReview int `json:"allow_review"`
    CanWatch int `json:"can_watch"`
    Copyright string `json:"copyright"`
}

type BilibiliMediaResultSeason struct {
    IsNew int `json:"is_new"`
    SeasonId int `json:"season_id"`
    SeasonTitle string `json:"season_title"`
}

type BilibiliMediaResultStat struct {
    Danmakus int `json:"danmakus"`
    Favorites int `json:"favorites"`
    Views int `json:"views"`
}

type BilibiliMediaResultStyle struct {
    Id int `json:"id"`
    Name string `json:"name"`
}

func makeBilibiliAnimeListUrl(query url.Values) (string, error) {
    u, err := url.Parse("https://api.bilibili.com/pgc/season/index/result")
    if err != nil {
        return "", err
    }

    q := url.Values{}
    q.Set("season_version", "-1")
    q.Set("area", "-1")
    q.Set("is_finish", "-1")
    q.Set("copyright", "-1")
    q.Set("season_status", "-1")
    q.Set("season_month", "-1")
    q.Set("year", "-1")
    q.Set("style_id", "-1")
    q.Set("order", "5")
    q.Set("st", "1")
    q.Set("sort", "1")
    q.Set("season_type", "1")
    q.Set("type", "1")
    q.Set("pagesize", "50")

    for key, value := range query {
        q.Set(key, value[0])
        for _, item := range value[1:] {
            q.Add(key, item)
        }
    }

    u.RawQuery = q.Encode()

    return u.String(), nil
}

func getBilibiliAnimeList(toChannel chan<- int) {
    defer close(toChannel)

    client := &http.Client{}
    for page := 1; ; page += 1 {
        q := url.Values{
            "page": []string{strconv.Itoa(page)},
        }

        bilibiliAnimeListUrl, err := makeBilibiliAnimeListUrl(q)
        if err != nil {
            log.Fatal(err)
        }

        result := func () (BilibiliAnimeList) {
            req, err := http.NewRequest("GET", bilibiliAnimeListUrl, nil)
            if err != nil {
                log.Fatal(err)
            }

            req.Header.Set("Connection", "close")
            resp, err := client.Do(req)
            if err != nil {
                log.Fatal(err)
            }
            defer resp.Body.Close()

            result := BilibiliAnimeList{}
            err = json.NewDecoder(resp.Body).Decode(&result)
            if err != nil {
                log.Fatal(err)
            }

            return result
        } ()

        if result.Data.HasNext != 1 {
            break
        }

        for _, item := range result.Data.List {
            toChannel <- item.MediaId
        }

    }
}

func makeBilibiliMediaUrl(query url.Values) (string, error) {
    u, err := url.Parse("https://api.bilibili.com/pgc/view/web/media")
    if err != nil {
        return "", err
    }

    q := url.Values{}
    for key, value := range query {
        q.Set(key, value[0])
        for _, item := range value[1:] {
            q.Add(key, item)
        }
    }

    u.RawQuery = q.Encode()

    return u.String(), nil
}

func getBilibiliMedia(fromChannel <-chan int, toChannel chan<- BilibiliMedia) {
    defer close(toChannel)

    client := &http.Client{}
    for mediaId := range fromChannel {
        q := url.Values{
            "media_id": []string{strconv.Itoa(mediaId)},
        }
        bilibiliMediaUrl, err := makeBilibiliMediaUrl(q)
        if err != nil {
            log.Fatal(err)
        }

        result := func () (BilibiliMedia) {
            req, err := http.NewRequest("GET", bilibiliMediaUrl, nil)
            if err != nil {
                log.Fatal(err)
            }

            resp, err := client.Do(req)
            if err != nil {
                log.Fatal(err)
            }

            defer resp.Body.Close()

            result := BilibiliMedia{}
            err = json.NewDecoder(resp.Body).Decode(&result)
            if err != nil {
                log.Fatal(err)
            }

            return result
        } ()

        toChannel <- result

    }
}

func writeDatabase(fromChannel <-chan BilibiliMedia) {
    db, err := gorm.Open("sqlite3", "./anime.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if !db.HasTable(&models.Media{}) {
        err := db.CreateTable(&models.Media{}).Error
        if err != nil {
            log.Fatal(err)
        }
    }

    if !db.HasTable(&models.Area{}) {
        err := db.CreateTable(&models.Area{}).Error
        if err != nil {
            log.Fatal(err)
        }
    }

    if !db.HasTable(&models.Style{}) {
        err := db.CreateTable(&models.Style{}).Error
        if err != nil {
            log.Fatal(err)
        }
    }

    for media := range fromChannel {
        result := media.Result

        fmt.Println(result.Publish.PubDate)
        pubDate, err := time.Parse("2006-01-02", result.Publish.PubDate)
        if err != nil {
            log.Fatal(err)
        }

        tx := db.Begin()
        defer func() {
            r := recover()
            if r != nil {
                tx.Rollback()
            }
        }()

        err = tx.Error
        if err != nil {
            log.Fatal(err)
        }

        var areasModel []*models.Area
        for _, area := range result.Areas {
            var areaModel models.Area
            err := tx.Where(models.Area{Name: area.Name}).FirstOrCreate(&areaModel).Error
            if err != nil {
                log.Fatal(err)
            }

            areasModel = append(areasModel, &areaModel)
        }

        var stylesModel []*models.Style
        for _, style := range result.Styles {
            var styleModel models.Style
            err := tx.Where(models.Style{Name: style.Name}).FirstOrCreate(&styleModel).Error
            if err != nil {
                log.Fatal(err)
            }

            stylesModel = append(stylesModel, &styleModel)
        }

        var mediaModel models.Media
        err = tx.Where(models.Media{BilibiliMediaId: result.MediaId}).Assign(models.Media{
            Actors: result.Actors,
            Alias: result.Alias,
            Areas: areasModel,
            BilibiliMediaId: result.MediaId,
            Copyright: result.Rights.Copyright,
            Cover: result.Cover,
            Evaluate: result.Evaluate,
            EpisodeIndex: result.EpisodeIndex.Index,
            OriginName: result.OriginName,
            IsFinish: result.Publish.IsFinish,
            IsStarted: result.Publish.IsStarted,
            PubDate: pubDate,
            RatingCount: result.Rating.Count,
            RatingScore: int(result.Rating.Score * 10),
            Staff: result.Staff,
            StatDanmakus: result.Stat.Danmakus,
            StatFavorites: result.Stat.Favorites,
            StatViews: result.Stat.Views,
            Styles: stylesModel,
            Title: result.Title,
        }).FirstOrCreate(&mediaModel).Error
        if err != nil {
            tx.Rollback()
            log.Fatal(err)
        }

        err = tx.Commit().Error
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println(mediaModel.Title)
    }

}

func main() {
    mediaIdChannel := make(chan int, 10)
    bilibiliMediaChannel := make(chan BilibiliMedia, 10)
    go getBilibiliAnimeList(mediaIdChannel)
    go getBilibiliMedia(mediaIdChannel, bilibiliMediaChannel)

    writeDatabase(bilibiliMediaChannel)
}
