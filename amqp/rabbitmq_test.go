package amqp

import ()

import (
	"github.com/tsuru/config"
	"launchpad.net/gocheck"
)

type RabbitMQSuite struct{}

var _ = gocheck.Suite(&RabbitMQSuite{})

func (s *RabbitMQSuite) SetUpSuite(c *gocheck.C) {
	//	config.Set("queue-server", "127.0.0.1:11300")
}

func (s *RabbitMQSuite) TestConnection(c *gocheck.C) {
	_, err := connection()
	c.Assert(err, gocheck.IsNil)
}

func (s *RabbitMQSuite) TestConnectionQueueServerUndefined(c *gocheck.C) {
	old, _ := config.Get("amqp:url")
	config.Unset("amqp:url")
	defer config.Set("amqp:url", old)
	conn, err := connection()
	c.Assert(err, gocheck.IsNil)
	c.Assert(conn, gocheck.NotNil)
}

func (s *RabbitMQSuite) TestConnectionResfused(c *gocheck.C) {
	old, _ := config.Get("amqp:url")
	config.Set("amqp:url", "127.0.0.1:11301")
	defer config.Set("amqp:url", old)
	conn, err := connection()
	c.Assert(conn, gocheck.IsNil)
	c.Assert(err, gocheck.NotNil)
}

/*func (s *RabbitMQSuite) TestPut(c *gocheck.C) {
	msg := Message{
		Action: "startapp",
		Args:   []string{"node1.megam.co"},
	}
	q := rabbitmqQ{name: "default"}
	err := q.Put(&msg, 0)
	c.Assert(err, gocheck.IsNil)
	c.Assert(msg.id, gocheck.Not(gocheck.Equals), 0)
	defer conn.Delete(msg.id)
	id, body, err := conn.Reserve(1e6)
	c.Assert(err, gocheck.IsNil)
	c.Assert(id, gocheck.Equals, msg.id)
	var got Message
	buf := bytes.NewBuffer(body)
	err = gob.NewDecoder(buf).Decode(&got)
	c.Assert(err, gocheck.IsNil)
	got.id = msg.id
	c.Assert(got, gocheck.DeepEquals, msg)
}*/

func (s *RabbitMQSuite) TestRabbitMQFactoryIsInFactoriesMap(c *gocheck.C) {
	f, ok := factories["rabbitmq"]
	c.Assert(ok, gocheck.Equals, true)
	_, ok = f.(rabbitmqFactory)
	c.Assert(ok, gocheck.Equals, true)
}
