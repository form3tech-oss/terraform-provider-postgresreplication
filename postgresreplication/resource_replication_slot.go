package postgresreplication

import (
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

const (
	slotNameAttributeName     = "slot_name"
	outputPluginAttributeName = "output_plugin"
	databaseAttributeName     = "database"
)

func resourceReplicationSlot() *schema.Resource {
	return &schema.Resource{
		Create: resourceReplicationSlotCreate,
		Read:   resourceReplicationSlotRead,
		Delete: resourceReplicationSlotDelete,
		Importer: &schema.ResourceImporter{
			State: resourceReplicationSlotImport,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			slotNameAttributeName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the slot to create. Must be a valid replication slot name.",
			},
			outputPluginAttributeName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the output plugin used for logical decoding.",
			},
			databaseAttributeName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the database this slot is associated with.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func connect(d *schema.ResourceData, m interface{}) (r *pgx.ReplicationConn, err error) {
	c := m.(*providerConfiguration)

	u, err := url.Parse(fmt.Sprintf("postgres://%s:%d/%s?sslmode=%s",
		c.host,
		c.port,
		d.Get(databaseAttributeName).(string),
		c.sslMode))
	if err != nil {
		return nil, errors.Wrap(err, "error contructing database connection uri.")
	}

	u.User = url.UserPassword(c.user, c.password)
	dbConfig, err := pgx.ParseURI(u.String())

	if err != nil {
		return nil, errors.Wrap(err, "error setting up database connection.")
	}

	replConn, err := pgx.ReplicationConnect(dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to database.")
	}

	return replConn, nil
}

func resourceReplicationSlotCreate(d *schema.ResourceData, m interface{}) error {
	replConn, err := connect(d, m)
	if err != nil {
		return err
	}
	defer replConn.Close()

	err = replConn.CreateReplicationSlot(d.Get(slotNameAttributeName).(string), d.Get(outputPluginAttributeName).(string))
	if err != nil {
		return errors.Wrap(err, "error creating replication slot.")
	}

	d.SetId(d.Get(slotNameAttributeName).(string))

	return nil
}

func resourceReplicationSlotImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	err := resourceReplicationSlotRead(d, m)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceReplicationSlotRead(d *schema.ResourceData, m interface{}) error {
	replConn, err := connect(d, m)
	if err != nil {
		return err
	}
	defer replConn.Close()

	r, err := replConn.Query("select slot_name, plugin, database from pg_replication_slots where slot_name=$1;", d.Id())
	if err != nil {
		return errors.Wrap(err, "error while trying to read existing replication slot")
	}
	defer r.Close()
	if r.Next() {
		v, _ := r.Values()
		err = d.Set(slotNameAttributeName, v[0])
		if err != nil {
			return errors.Wrap(err, "error reading slot name")
		}
		err = d.Set(outputPluginAttributeName, v[1])
		if err != nil {
			return errors.Wrap(err, "error reading output plugin")
		}
		err = d.Set(databaseAttributeName, v[2])
		if err != nil {
			return errors.Wrap(err, "error reading database")
		}
	}

	return nil
}

func resourceReplicationSlotDelete(d *schema.ResourceData, m interface{}) error {
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		replConn, err := connect(d, m)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		defer replConn.Close()

		err = replConn.DropReplicationSlot(d.Id())
		if err != nil {
			return resource.RetryableError(errors.Wrap(err, "error dropping replication slot."))
		}

		return nil
	})
}
