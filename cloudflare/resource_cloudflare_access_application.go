package cloudflare

import (
	"context"
	"fmt"
	"log"
	"strings"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudflareAccessApplication() *schema.Resource {
	return &schema.Resource{
		Schema: resourceCloudflareAccessApplicationSchema(),
		Create: resourceCloudflareAccessApplicationCreate,
		Read:   resourceCloudflareAccessApplicationRead,
		Update: resourceCloudflareAccessApplicationUpdate,
		Delete: resourceCloudflareAccessApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudflareAccessApplicationImport,
		},
<<<<<<< HEAD
=======

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"zone_id"},
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"account_id"},
			},
			"aud": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "self_hosted",
				ValidateFunc: validation.StringInSlice([]string{"self_hosted", "ssh", "vnc", "file"}, false),
			},
			"session_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "24h",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := time.ParseDuration(v)
					if err != nil {
						errs = append(errs, fmt.Errorf(`%q only supports "ns", "us" (or "µs"), "ms", "s", "m", or "h" as valid units.`, key))
					}
					return
				},
			},
			"cors_headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_methods": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"allowed_origins": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"allowed_headers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"allow_all_methods": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"allow_all_origins": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"allow_all_headers": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"allow_credentials": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"max_age": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(-1, 86400),
						},
					},
				},
			},
			"auto_redirect_to_identity": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_binding_cookie": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allowed_idps": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"custom_deny_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"custom_deny_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_only_cookie_attribute": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"same_site_cookie_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logo_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"skip_interstitial": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"app_launcher_visible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
>>>>>>> 425dd626 (Add app_launcher_visible to Access Application)
	}
}

func resourceCloudflareAccessApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	allowedIDPList := expandInterfaceToStringList(d.Get("allowed_idps"))
	appType := d.Get("type").(string)

	newAccessApplication := cloudflare.AccessApplication{
		Name:                    d.Get("name").(string),
		Domain:                  d.Get("domain").(string),
		Type:                    cloudflare.AccessApplicationType(appType),
		SessionDuration:         d.Get("session_duration").(string),
		AutoRedirectToIdentity:  d.Get("auto_redirect_to_identity").(bool),
		EnableBindingCookie:     d.Get("enable_binding_cookie").(bool),
		CustomDenyMessage:       d.Get("custom_deny_message").(string),
		CustomDenyURL:           d.Get("custom_deny_url").(string),
		HttpOnlyCookieAttribute: d.Get("http_only_cookie_attribute").(bool),
		SameSiteCookieAttribute: d.Get("same_site_cookie_attribute").(string),
		LogoURL:                 d.Get("logo_url").(string),
		SkipInterstitial:        d.Get("skip_interstitial").(bool),
		AppLauncherVisible:      d.Get("app_launcher_visible").(bool),
	}

	if len(allowedIDPList) > 0 {
		newAccessApplication.AllowedIdps = allowedIDPList
	}

	if _, ok := d.GetOk("cors_headers"); ok {
		CORSConfig, err := convertCORSSchemaToStruct(d)
		if err != nil {
			return err
		}
		newAccessApplication.CorsHeaders = CORSConfig
	}

	log.Printf("[DEBUG] Creating Cloudflare Access Application from struct: %+v", newAccessApplication)

	identifier, err := initIdentifier(d)
	if err != nil {
		return err
	}

	var accessApplication cloudflare.AccessApplication
	if identifier.Type == AccountType {
		accessApplication, err = client.CreateAccessApplication(context.Background(), identifier.Value, newAccessApplication)
	} else {
		accessApplication, err = client.CreateZoneLevelAccessApplication(context.Background(), identifier.Value, newAccessApplication)
	}
	if err != nil {
		return fmt.Errorf("error creating Access Application for %s %q: %s", identifier.Type, identifier.Value, err)
	}

	d.SetId(accessApplication.ID)

	return resourceCloudflareAccessApplicationRead(d, meta)
}

func resourceCloudflareAccessApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	identifier, err := initIdentifier(d)
	if err != nil {
		return err
	}

	var accessApplication cloudflare.AccessApplication
	if identifier.Type == AccountType {
		accessApplication, err = client.AccessApplication(context.Background(), identifier.Value, d.Id())
	} else {
		accessApplication, err = client.ZoneLevelAccessApplication(context.Background(), identifier.Value, d.Id())
	}

	if err != nil {
		if strings.Contains(err.Error(), "HTTP status 404") {
			log.Printf("[INFO] Access Application %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error finding Access Application %q: %s", d.Id(), err)
	}

	d.Set("name", accessApplication.Name)
	d.Set("aud", accessApplication.AUD)
	d.Set("session_duration", accessApplication.SessionDuration)
	d.Set("domain", accessApplication.Domain)
	d.Set("type", accessApplication.Type)
	d.Set("auto_redirect_to_identity", accessApplication.AutoRedirectToIdentity)
	d.Set("enable_binding_cookie", accessApplication.EnableBindingCookie)
	d.Set("custom_deny_message", accessApplication.CustomDenyMessage)
	d.Set("custom_deny_url", accessApplication.CustomDenyURL)
	d.Set("allowed_idps", accessApplication.AllowedIdps)
	d.Set("http_only_cookie_attribute", accessApplication.HttpOnlyCookieAttribute)
	d.Set("same_site_cookie_attribute", accessApplication.SameSiteCookieAttribute)
	d.Set("skip_interstitial", accessApplication.SkipInterstitial)
	d.Set("logo_url", accessApplication.LogoURL)
	d.Set("app_launcher_visible", accessApplication.AppLauncherVisible)

	corsConfig := convertCORSStructToSchema(d, accessApplication.CorsHeaders)
	if corsConfigErr := d.Set("cors_headers", corsConfig); corsConfigErr != nil {
		return fmt.Errorf("error setting Access Application CORS header configuration: %s", corsConfigErr)
	}

	return nil
}

func resourceCloudflareAccessApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	allowedIDPList := expandInterfaceToStringList(d.Get("allowed_idps"))
	appType := d.Get("type").(string)

	updatedAccessApplication := cloudflare.AccessApplication{
		ID:                      d.Id(),
		Name:                    d.Get("name").(string),
		Domain:                  d.Get("domain").(string),
		Type:                    cloudflare.AccessApplicationType(appType),
		SessionDuration:         d.Get("session_duration").(string),
		AutoRedirectToIdentity:  d.Get("auto_redirect_to_identity").(bool),
		EnableBindingCookie:     d.Get("enable_binding_cookie").(bool),
		CustomDenyMessage:       d.Get("custom_deny_message").(string),
		CustomDenyURL:           d.Get("custom_deny_url").(string),
		HttpOnlyCookieAttribute: d.Get("http_only_cookie_attribute").(bool),
		SameSiteCookieAttribute: d.Get("same_site_cookie_attribute").(string),
		LogoURL:                 d.Get("logo_url").(string),
		SkipInterstitial:        d.Get("skip_interstitial").(bool),
		AppLauncherVisible:      d.Get("app_launcher_visible").(bool),
	}

	if len(allowedIDPList) > 0 {
		updatedAccessApplication.AllowedIdps = allowedIDPList
	}

	if _, ok := d.GetOk("cors_headers"); ok {
		CORSConfig, err := convertCORSSchemaToStruct(d)
		if err != nil {
			return err
		}
		updatedAccessApplication.CorsHeaders = CORSConfig
	}

	log.Printf("[DEBUG] Updating Cloudflare Access Application from struct: %+v", updatedAccessApplication)

	identifier, err := initIdentifier(d)
	if err != nil {
		return err
	}

	var accessApplication cloudflare.AccessApplication
	if identifier.Type == AccountType {
		accessApplication, err = client.UpdateAccessApplication(context.Background(), identifier.Value, updatedAccessApplication)
	} else {
		accessApplication, err = client.UpdateZoneLevelAccessApplication(context.Background(), identifier.Value, updatedAccessApplication)
	}
	if err != nil {
		return fmt.Errorf("error updating Access Application for %s %q: %s", identifier.Type, identifier.Value, err)
	}

	if accessApplication.ID == "" {
		return fmt.Errorf("failed to find Access Application ID in update response; resource was empty")
	}

	return resourceCloudflareAccessApplicationRead(d, meta)
}

func resourceCloudflareAccessApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	appID := d.Id()

	log.Printf("[DEBUG] Deleting Cloudflare Access Application using ID: %s", appID)

	identifier, err := initIdentifier(d)
	if err != nil {
		return err
	}

	if identifier.Type == AccountType {
		err = client.DeleteAccessApplication(context.Background(), identifier.Value, appID)
	} else {
		err = client.DeleteZoneLevelAccessApplication(context.Background(), identifier.Value, appID)
	}
	if err != nil {
		return fmt.Errorf("error deleting Access Application for %s %q: %s", identifier.Type, identifier.Value, err)
	}

	resourceCloudflareAccessApplicationRead(d, meta)

	return nil
}

func resourceCloudflareAccessApplicationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	attributes := strings.SplitN(d.Id(), "/", 2)

	if len(attributes) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"accountID/accessApplicationID\"", d.Id())
	}

	accountID, accessApplicationID := attributes[0], attributes[1]

	log.Printf("[DEBUG] Importing Cloudflare Access Application: id %s for account %s", accessApplicationID, accountID)

	d.Set("account_id", accountID)
	d.SetId(accessApplicationID)

	resourceCloudflareAccessApplicationRead(d, meta)

	return []*schema.ResourceData{d}, nil
}
