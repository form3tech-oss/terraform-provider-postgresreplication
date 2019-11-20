package postgresreplication

import "github.com/hashicorp/terraform/helper/schema"

const (
	replicationSlotName = "postgresreplication_slot"
)

const (
	portKey 		 = "port"
	hostKey          = "host"
	userKey          = "user"
	passwordKey      = "password"
)

const (
	defaultPort = 5432
	defaultHost = "localhost"
	defaultUser = "postgres"
	defaultPassword = ""
)

type providerConfiguration struct {
	port              uint16
	host          	  string
	user              string
	password          string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			return &providerConfiguration{
				port:      uint16(d.Get(portKey).(int)),
				host:      d.Get(hostKey).(string),
				user:      d.Get(userKey).(string),
				password:  d.Get(passwordKey).(string),
			}, nil
		},
		ResourcesMap: map[string]*schema.Resource{
			replicationSlotName: resourceReplicationSlot(),
		},
		Schema: map[string]*schema.Schema{
			hostKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   false,
				Default:	 defaultHost,
				Description: "The server host to connect to.",
			},
			portKey: {
				Type:        schema.TypeInt,
				Optional:    true,
				Sensitive:   false,
				Default:	 defaultPort,
				Description: "The server port to connect to.",
			},
			userKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:	 defaultUser,
				Description: "The user to use to connect.",
			},
			passwordKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     defaultPassword,
				Description: "The password to use to connect.",
			},
		},
	}
}