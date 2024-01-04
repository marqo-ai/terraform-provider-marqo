package marqo

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMarqoIndex() *schema.Resource {
    return &schema.Resource{
        Create: resourceMarqoIndexCreate,
        Update: resourceMarqoIndexUpdate,
        Delete: resourceMarqoIndexDelete,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
            },
        },
    }
}

func resourceMarqoIndexCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*MarqoClient)
	indexName := d.Get("name").(string)

	// Define the settings for the index, modify as needed
	settings := IndexSettings{
		IndexDefaults: IndexDefaults{
			// Set defaults
		},
		// other settings
	}

	err := client.CreateIndex(indexName, settings)
	if err != nil {
		return err
	}

	d.SetId(indexName)
	return resourceMarqoIndexRead(d, m)
}

func resourceMarqoIndexUpdate(d *schema.ResourceData, m interface{}) error {
    // update a Marqo index 
	// use update document and delete document functions from client.go
    return nil
}

func resourceMarqoIndexDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*MarqoClient)
	indexName := d.Id()

	return client.DeleteIndex(indexName)
}
