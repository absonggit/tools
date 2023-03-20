package main

import (
    "archive/zip"
    "io"
    "os"
    "time"
    "log"
    "github.com/spf13/viper"
    tgbotapi "gopkg.in/telegram-bot-api.v4"
)
type Env struct {
     Chat_id string `yaml:"CHAT_ID"`
     Bot_token string `yaml:"BOT_TOKEN"`
     Files []string `yaml:"FILES"`
     Prefix string `yaml:"PREFIX"`
     Zippath string `yaml:"ZIP_PATH"`
}
func init() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.SetConfigFile("./config.yaml")
    viper.AddConfigPath("./")
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            log.Println("no such config file")
        } else {
            log.Println("read config error")
        }
        log.Fatal(err)
    }
    log.Println("配置读取成功。")
}

func main() {
    var v interface{}
    var files []string
    a := viper.Get(`FILES`) 
    v = a 
    for _, val := range v.([]interface{}) {
        files = append(files, val.(string))
    }
    currentTime := time.Now()
    now := currentTime.Format("2006-01-02_3-4-5_PM")
    output := viper.GetString(`PREFIX`) + now + ".zip"
    if err := ZipFiles(output, files); err != nil {
        log.Fatal(err)
    }
    log.Println("文件压缩完成。")
    bot, err := tgbotapi.NewBotAPI(viper.GetString(`BOT_TOKEN`))
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = false
    msg := tgbotapi.NewDocumentUpload(int64(viper.GetInt(`CHAT_ID`)), output)
    bot.Send(msg)
    log.Println("发送压缩文件。")
    err = os.Remove(output)
    if err != nil {
        log.Panic(err)
    }
    log.Println("压缩文件删除。")
}
func ZipFiles(filename string, files []string) error {
    log.Println("开始压缩文件。")
    newZipFile, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer newZipFile.Close()
    zipWriter := zip.NewWriter(newZipFile)
    defer zipWriter.Close()
    for _, file := range files {
        if err = AddFileToZip(zipWriter, file); err != nil {
            return err
        }
    }
    return nil
}
func AddFileToZip(zipWriter *zip.Writer, filename string) error {
    fileToZip, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer fileToZip.Close()
    info, err := fileToZip.Stat()
    if err != nil {
        return err
    }
    header, err := zip.FileInfoHeader(info)
    if err != nil {
        return err
    }
    header.Name = filename
    /*
      Store   unit16 = 0  //仅存储文件
      Deflate unit16 = 8  //压缩文件
    */
    header.Method = zip.Deflate
    writer, err := zipWriter.CreateHeader(header)
    if err != nil {
        return err
    }
    _, err = io.Copy(writer, fileToZip)
    return err
}
