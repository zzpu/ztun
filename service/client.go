package service

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	l            net.Listener
	websocketURL string
	info         func(text string)
	running      bool
	cons         []net.Conn
}

func NewClient(port string, websocketURL string) (c *Client, err error) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return
	}
	c = &Client{
		l:            listener,
		websocketURL: websocketURL,
		cons:         make([]net.Conn, 0),
	}

	return
}
func (c *Client) DisableLog() {
	c.info = nil
}

func (c *Client) EnableLog(info func(text string)) {
	c.info = info
}

func (c *Client) Close() {
	c.l.Close()
	for _, con := range c.cons {
		con.Close()
	}
	c.running = false
}

func (c *Client) IsRunning() (b bool) {
	return c.running
}

func (c *Client) Start() {
	c.running = true
	for {
		tcp, err := c.l.Accept()
		if err != nil {
			c.info(fmt.Sprintf("连接出错,err=%v,", err))
			return
		}
		c.cons = append(c.cons, tcp)
		go func() {

			defer tcp.Close()
			ws, _, err := websocket.DefaultDialer.Dial(c.websocketURL, nil)
			if err != nil {
				log.Printf("连接失败,err=%v", err)
				return
			}
			defer ws.Close()

			go func() {
				for {
					//30秒发送一次心跳信号
					time.Sleep(time.Second * 30)
					if c.info != nil {
						go c.info(fmt.Sprintf("客户端>>>服务器,心跳..."))
					}
					err = ws.WriteMessage(websocket.PingMessage, []byte("ok"))
					if err != nil {
						log.Printf("读取失败,err=%v", err)
						tcp.Close()
						ws.Close()
						break
					}
				}

			}()

			//客户端流向服务端数据处理
			go func() {
				buf := make([]byte, 1024)
				for {
					len, err := tcp.Read(buf)
					if err != nil {
						log.Printf("读取失败,err=%v", err)
						tcp.Close()
						ws.Close()
						break
					}
					//通知函数不为空，则输出
					if c.info != nil {
						go c.info(fmt.Sprintf("客户端>>>服务器,数据:%d", len))
					}
					log.Printf("客户端>>>服务器,数据:%d", len)
					ws.WriteMessage(websocket.BinaryMessage, buf[0:len])
				}
			}()
			//服务端流向客户端数据处理
			for {
				msgType, buf, err := ws.ReadMessage()
				if err != nil {
					log.Printf("读取失败,err=%v", err)
					tcp.Close()
					ws.Close()
					break
				}
				if msgType != websocket.BinaryMessage {
					log.Println("unknown msgType")
				}
				//通知函数不为空，则输出
				if c.info != nil {
					go c.info(fmt.Sprintf("服务器>>>客户端,数据:%d", len(buf)))
				}
				log.Printf("服务器>>>客户端,数据:%d", len(buf))
				tcp.Write(buf)
			}
		}()
	}
}
