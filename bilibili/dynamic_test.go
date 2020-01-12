package biliAPI

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"testing"
	"time"
)

func TestGetDynamic(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ret := &DynamicRetStruct{}
		var err error
		curs := int64(0)
		for {
			ret, err = GetDynamicSpaceHistory(3379951, curs, true)
			if err == nil {
				if ret.Data.Attentions != nil {
					fmt.Println("关注列表: ", ret.Data.Attentions.UIds)
				}
				fmt.Println("/////////////////////////////////////////////////////////////////////////////////////")

				for _, c := range ret.Data.Cards {
					fmt.Println("动态ID: ", c.Desc.DynamicId)
					fmt.Println("动态发布时间: ", time.Unix(c.Desc.Timestamp, 0))
					fmt.Println("动态类型: ", c.Desc.Type, DynamicTypeId2String(c.Desc.Type))
					switch c.Desc.Type {
					case 1: // 引用
						fmt.Println("引用类型: ", c.CardContent.T1.OriginType)
						fmt.Println("原文: ", c.CardContent.T1.Item.Content) // 原文
						fmt.Printf("> 引用 ")
						switch c.CardContent.T1.OriginType {
						case 2:
							fmt.Println("动态")
							fmt.Println("> 引用账户 : ", c.CardContent.T1.OriginUser.Info.UName)
							fmt.Println("> 引用动态内容: ", c.CardContent.T1.OriginContent.T2.Item.Description)
							fmt.Println("> 引用动态发布时间: ", time.Unix(c.CardContent.T1.OriginContent.T2.Item.UploadTime, 0))
							fmt.Println("> 引用动态附带图片数量: ", c.CardContent.T1.OriginContent.T2.Item.PicturesCount)
							for i, p := range c.CardContent.T1.OriginContent.T2.Item.Pictures {
								fmt.Println("第", i+1, "张: ", p.ImgSrc)
							}
						case 4:
							fmt.Println("活动")
							fmt.Println("> 引用活动内容: ", c.CardContent.T1.OriginContent.T4.Item.Content)
							fmt.Println("活动发布时间: ", time.Unix(c.CardContent.T1.OriginContent.T4.Item.Timestamp, 0))
						case 8:
							fmt.Println("投稿")
							fmt.Println("> 引用账户 : ", c.CardContent.T1.OriginUser.Info.UName)
							fmt.Println("> 引用投稿ID: ", c.CardContent.T1.OriginContent.T8.Aid)
							fmt.Println("> 引用投稿标题: ", c.CardContent.T1.OriginContent.T8.Title)
							fmt.Println("> 引用投稿描述: ", c.CardContent.T1.OriginContent.T8.Desc)
							fmt.Println("> 引用投稿封面: ", c.CardContent.T1.OriginContent.T8.Pic)
						case 16:
							fmt.Println("小视频")
							fmt.Println("> 引用账户 : ", c.CardContent.T1.OriginUser.Info.UName)
							fmt.Println("> 引用小视频ID: ", c.CardContent.T1.OriginContent.T16.Item.Id)
							fmt.Println("> 引用小视频图文: ", c.CardContent.T1.OriginContent.T16.Item.Description)
							fmt.Println("> 引用小视频链接: ", c.CardContent.T1.OriginContent.T16.Item.VideoPlayUrl)
						case 32:
							panic("捕获引用类型32!")
						case 64:
							fmt.Println("专栏")
							fmt.Println("> 引用账户 : ", c.CardContent.T1.OriginUser.Info.UName)
							fmt.Println("> 引用专栏ID: ", c.CardContent.T1.OriginContent.T64.Id)
							fmt.Println("> 引用专栏动态: ", c.CardContent.T1.OriginContent.T64.Dynamic)
							fmt.Println("> 引用专栏标题: ", c.CardContent.T1.OriginContent.T64.Title)
							fmt.Println("> 引用专栏摘要: ", c.CardContent.T1.OriginContent.T64.Summary)
						case 256:
							fmt.Println("创作")
							fmt.Println("> 引用作品标题: ", c.CardContent.T1.OriginContent.T256.Title)
							fmt.Println("> 引用作品介绍: ", c.CardContent.T1.OriginContent.T256.Intro)
							fmt.Println("> 引用作品类型: ", c.CardContent.T1.OriginContent.T256.TypeInfo)
							fmt.Println("> 引用作者: ", c.CardContent.T1.OriginContent.T256.Author)
						case 512:
							fmt.Println("番剧")
							fmt.Println("> 引用番剧ID: ", c.CardContent.T1.OriginContent.T512.Aid)
							fmt.Println("> 引用番剧名: ", c.CardContent.T1.OriginContent.T512.ApiSeasonInfo.Title)
							fmt.Println("> 引用番剧封面: ", c.CardContent.T1.OriginContent.T512.Cover)
							fmt.Println("> 引用番剧第", c.CardContent.T1.OriginContent.T512.Index, "集")
							fmt.Println("> 引用番剧最新一集描述: ", c.CardContent.T1.OriginContent.T512.NewDesc)
							fmt.Println("> 引用番剧播放量: ", c.CardContent.T1.OriginContent.T512.PlayCount)
						case 1024:
							fmt.Println("已失效资源")
							fmt.Println("> 引用提示: ", c.CardContent.T1.Item.Tips)
						case 2048:
							fmt.Println("宣传")
							fmt.Println("> 引用宣传动态内容", c.CardContent.T1.OriginContent.T2048.Vest.Content)
							fmt.Println("> 引用宣传页标题", c.CardContent.T1.OriginContent.T2048.Sketch.Title)
							fmt.Println("> 引用宣传页封面", c.CardContent.T1.OriginContent.T2048.Sketch.CoverUrl)
							fmt.Println("> 引用宣传页描述", c.CardContent.T1.OriginContent.T2048.Sketch.DescText)
							fmt.Println("> 引用宣传页链接", c.CardContent.T1.OriginContent.T2048.Sketch.TargetUrl)
						case 4099:
							fmt.Println("番剧")
							fmt.Println("> 引用番剧ID: ", c.CardContent.T1.OriginContent.T4099.Aid)
							fmt.Println("> 引用番剧名: ", c.CardContent.T1.OriginContent.T4099.ApiSeasonInfo.Title)
							fmt.Println("> 引用番剧封面: ", c.CardContent.T1.OriginContent.T4099.Cover)
							fmt.Println("> 引用番剧第", c.CardContent.T1.OriginContent.T4099.Index, "集")
							fmt.Println("> 引用番剧最新一集描述: ", c.CardContent.T1.OriginContent.T4099.NewDesc)
							fmt.Println("> 引用番剧播放量: ", c.CardContent.T1.OriginContent.T4099.PlayCount)
						case 4200:
							fmt.Println("房间")
							fmt.Println("> 引用账户 : ", c.CardContent.T1.OriginUser.Info.UName)
							fmt.Println("> 引用房间ID: ", c.CardContent.T1.OriginContent.T4200.RoomId)
							fmt.Println("> 引用房间标题: ", c.CardContent.T1.OriginContent.T4200.Title)
							fmt.Println("> 引用房间标签: ", c.CardContent.T1.OriginContent.T4200.Tags)
						case 4300:
							fmt.Println("收藏夹")
							fmt.Println("> 引用收藏夹ID : ", c.CardContent.T1.OriginContent.T4300.Fid)
							fmt.Println("> 引用收藏夹标题 : ", c.CardContent.T1.OriginContent.T4300.Title)
							fmt.Println("> 引用收藏夹所属 : ", c.CardContent.T1.OriginContent.T4300.Upper.Name)
						default:
							panic(fmt.Sprint("捕获未知引用类型!", c.CardContent.T1.OriginType, c.CardContent.T1.OriginString))
						}
					case 2: // 图文
						fmt.Println("图文内容: ", c.CardContent.T2.Item.Description)
						fmt.Println("发布时间: ", time.Unix(c.CardContent.T2.Item.UploadTime, 0))
						fmt.Println("附带图片数量: ", c.CardContent.T2.Item.PicturesCount)
						for i, p := range c.CardContent.T2.Item.Pictures {
							fmt.Println("第", i+1, "张: ", p.ImgSrc)
						}
					case 4:
						fmt.Println("活动内容: ", c.CardContent.T4.Item.Content)
						fmt.Println("活动发布时间: ", time.Unix(c.CardContent.T4.Item.Timestamp, 0))
					case 8: //投稿
						fmt.Println("投稿ID: ", c.CardContent.T8.Aid)
						fmt.Println("投稿标题: ", c.CardContent.T8.Title)
						fmt.Println("投稿描述: ", c.CardContent.T8.Desc)
						fmt.Println("投稿时间: ", time.Unix(c.CardContent.T8.CTime, 0))
						//fmt.Println(c.CardContent.JumpUrl)
						fmt.Println("投稿封面: ", c.CardContent.T8.Pic)
					case 16: // 小视频
						fmt.Println("> 小视频ID: ", c.CardContent.T16.Item.Id)
						fmt.Println("> 小视频图文: ", c.CardContent.T16.Item.Description)
						fmt.Println("> 小视频链接: ", c.CardContent.T16.Item.VideoPlayUrl)
					case 32:
						panic("捕获未知动态类型32!")
					case 64: // 专栏
						fmt.Println("专栏ID: ", c.CardContent.T64.Id)
						fmt.Println("专栏标题: ", c.CardContent.T64.Title)
						fmt.Println("专栏摘要: ", c.CardContent.T64.Summary)
						fmt.Println("专栏横幅: ", c.CardContent.T64.BannerUrl)
						fmt.Println("专栏发布日期: ", time.Unix(c.CardContent.T64.PublishTime, 0))
						fmt.Println("专栏创建日期: ", time.Unix(c.CardContent.T64.CTime, 0))
						fmt.Println("专栏字数: ", c.CardContent.T64.Words)
						for i, img := range c.CardContent.T64.ImageUrls {
							fmt.Printf("第%d张专栏图: %s\n", i+1, img)
						}
					case 256:
						fmt.Println("作品标题: ", c.CardContent.T256.Title)
						fmt.Println("作品介绍: ", c.CardContent.T256.Intro)
						fmt.Println("作品类型: ", c.CardContent.T256.TypeInfo)
						fmt.Println("作者: ", c.CardContent.T256.Author)
					case 512:
						fmt.Println("番剧ID: ", c.CardContent.T512.Aid)
						fmt.Println("番剧名: ", c.CardContent.T512.ApiSeasonInfo.Title)
						fmt.Println("番剧封面: ", c.CardContent.T512.Cover)
						fmt.Println("番剧第", c.CardContent.T512.Index, "集")
						fmt.Println("番剧最新一集描述: ", c.CardContent.T512.NewDesc)
						fmt.Println("番剧播放量: ", c.CardContent.T512.PlayCount)
					case 1024:
						fmt.Println("已失效")
					case 2048:
						fmt.Println("宣传动态内容", c.CardContent.T2048.Vest.Content)
						fmt.Println("宣传页标题", c.CardContent.T2048.Sketch.Title)
						fmt.Println("宣传页封面", c.CardContent.T2048.Sketch.CoverUrl)
						fmt.Println("宣传页描述", c.CardContent.T2048.Sketch.DescText)
						fmt.Println("宣传页链接", c.CardContent.T2048.Sketch.TargetUrl)

					default:
						logrus.Println("动态字符串", c.Desc.Type, c.CardString)
						panic("WAIT!")
					}
					fmt.Printf("转发数： %d\t评论数: %d\t点赞数: %d\n", c.Desc.RePost, c.Desc.Comment, c.Desc.Like)
					fmt.Println("/////////////////////////////////////////////////////////////////////////////////")
				}
				log.Println("下一页偏移: ", ret.Data.NextOffset)
				if ret.Data.NextOffset == 0 {
					log.Println("所有动态结束")
					return
				}
			} else {
				panic(err)
			}
			curs = ret.Data.NextOffset
			time.Sleep(time.Second * 2)
			break
		}
	})

}
