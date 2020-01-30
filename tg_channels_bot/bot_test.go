package tg_channels_bot

import (
	"log"
	"testing"
)

func TestSendMsgToChannel(t *testing.T) {
	type args struct {
		channel *ChannelStruct
		msgText string
	}
	testName := "(TestSendMsgToChannel) Send to channel sanqiu"
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: testName + "0",
			args: args{
				channel: Config.Channels["sanqiu"],
				msgText: "hello from " + Config.Bots["0"].Name + testName + "0",
			},
			wantErr: false, // case 0: no error
		},
		{
			name: testName + "1",
			args: args{
				channel: Config.Channels["gadio"],
				msgText: "hello from " + Config.Bots["0"].Name + testName + "1",
			},
			wantErr: false, // case1: no error
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if resp, err := tt.args.channel.Notify(tt.args.msgText); (err != nil) != tt.wantErr {
				t.Errorf("Case#%2d: SendMsgToChannel() error = %v, wantErr %v", i, err, tt.wantErr)
			} else if err != nil {
				//log.Printf("Case#%2d: hit expected error: %v", i, err)
			} else {
				if _, err := tt.args.channel.DelMsg(resp.MessageID, resp.Chat.ID); err != nil {
					t.Errorf("dele msg[id: %v; chat_id: %v] error: %v\n", resp.MessageID, resp.Chat.ID, err)
				}
			}
		})
	}
}

func TestInitConfigFromToml(t *testing.T) {
	t.Run("", func(t *testing.T) {
		log.Println("Over!")
	})
}

func TestChannelStruct_QueryGadio(t *testing.T) {
	t.Run("", func(t *testing.T) {
		err := Config.Channels["gadio"].QueryGadio(5)
		if err != nil {
			log.Println(err)
		}
	})
}
