package amqplb

type SubscribeOptions struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	Args       map[string]interface{}
}

func NewDefaultSubscribeOptions() *SubscribeOptions {
	s := new(SubscribeOptions)

	s.Durable = true
	s.Exclusive = false
	s.AutoDelete = false
	s.Args = make(map[string]interface{})

	return s
}

func (o *SubscribeOptions) SetDurable(d bool) *SubscribeOptions {
	o.Durable = d
	return o
}

func (o *SubscribeOptions) SetExclusive(e bool) *SubscribeOptions {
	o.Exclusive = e
	return o
}

func (o *SubscribeOptions) SetAutoDelete(ad bool) *SubscribeOptions {
	o.AutoDelete = ad
	return o
}

func (o *SubscribeOptions) SetArgs(args map[string]interface{}) *SubscribeOptions {
	o.Args = args
	return o
}

func (o *SubscribeOptions) SetArg(name string, value interface{}) *SubscribeOptions {
	o.Args[name] = value
	return o
}

type PublisherOptions struct {
	Durable    bool
	AutoDelete bool
	Type       PublishType
}

type PublishType string

const (
	PublishTypeDirect          PublishType = "direct"
	PublishTypeTopic           PublishType = "topic"
	PublishTypeFanout          PublishType = "fanout"
	PublishTypeXConsistentHash PublishType = "x-consistent-hash"
)

func NewDefaultPublisherOptions() *PublisherOptions {
	s := new(PublisherOptions)

	s.Durable = true
	s.Type = PublishTypeDirect
	s.AutoDelete = false

	return s
}

func (o *PublisherOptions) SetDurable(d bool) *PublisherOptions {
	o.Durable = d
	return o
}

func (o *PublisherOptions) SetType(t PublishType) *PublisherOptions {
	o.Type = t
	return o
}

func (o *PublisherOptions) SetAutoDelete(ad bool) *PublisherOptions {
	o.AutoDelete = ad
	return o
}
