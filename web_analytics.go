package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	ErrMissingWebAnalyticsSiteTag   = errors.New("missing required web analytics site ID")
	ErrMissingWebAnalyticsRulesetID = errors.New("missing required web analytics ruleset ID")
	ErrMissingWebAnalyticsRuleID    = errors.New("missing required web analytics rule ID")
	ErrMissingWebAnalyticsSiteHost  = errors.New("missing required web analytics host or zone tag")
)

// WebAnalyticsSite describes a Web Analytics Site object.
type WebAnalyticsSite struct {
	SiteTag   string    `json:"site_tag"`
	SiteToken string    `json:"site_token"`
	Created   time.Time `json:"created,omitempty"`
	// Snippet is an encoded JS script to insert into your site HTML.
	Snippet string `json:"snippet"`
	// AutoInstall defines whether Cloudflare will inject the JS snippet automatically for orange-clouded sites.
	AutoInstall bool                `json:"auto_install"`
	Ruleset     WebAnalyticsRuleset `json:"ruleset"`
	Rules       []WebAnalyticsRule  `json:"rules"`
}

// WebAnalyticsRule describes a Web Analytics Rule object.
type WebAnalyticsRule struct {
	ID    string   `json:"id,omitempty"`
	Host  string   `json:"host"`
	Paths []string `json:"paths"`
	// Inclusive defines whether the rule includes or excludes the matched traffic from being measured in web analytics.
	Inclusive bool      `json:"inclusive"`
	Created   time.Time `json:"created,omitempty"`
	// IsPaused defines whether the rule is paused (inactive) or not.
	IsPaused bool `json:"is_paused"`
	Priority int  `json:"priority"`
}

// WebAnalyticsRuleset describes a Web Analytics Ruleset object.
type WebAnalyticsRuleset struct {
	ID       string `json:"id"`
	ZoneTag  string `json:"zone_tag"`
	ZoneName string `json:"zone_name"`
	Enabled  bool   `json:"enabled"`
}

// WebAnalyticsSiteResponse is the API response, containing a single WebAnalyticsSite.
type WebAnalyticsSiteResponse struct {
	Response
	Result WebAnalyticsSite `json:"result"`
}

// WebAnalyticsSitesResponse is the API response, containing an array of WebAnalyticsSite.
type WebAnalyticsSitesResponse struct {
	Response
	ResultInfo ResultInfo         `json:"result_info"`
	Result     []WebAnalyticsSite `json:"result"`
}

// WebAnalyticsRuleResponse is the API response, containing a single WebAnalyticsRule.
type WebAnalyticsRuleResponse struct {
	Response
	Result WebAnalyticsRule `json:"result"`
}

type WebAnalyticsRulesetRules struct {
	Ruleset WebAnalyticsRuleset `json:"ruleset"`
	Rules   []WebAnalyticsRule  `json:"rules"`
}

// WebAnalyticsRulesResponse is the API response, containing a WebAnalyticsRuleset and array of WebAnalyticsRule.
type WebAnalyticsRulesResponse struct {
	Response
	Result WebAnalyticsRulesetRules `json:"result"`
}

// WebAnalyticsIDResponse is the API response, containing a single ID.
type WebAnalyticsIDResponse struct {
	Response
	Result string `json:"result"`
}

// WebAnalyticsSiteTagResponse is the API response, containing a single ID.
type WebAnalyticsSiteTagResponse struct {
	Response
	Result struct {
		SiteTag string `json:"site_tag"`
	} `json:"result"`
}

type CreateWebAnalyticsSiteParams struct {
	// Host is the host to measure traffic for.
	Host string `json:"host,omitempty"`
	// ZoneTag is the zone tag to measure traffic for.
	ZoneTag string `json:"zone_tag,omitempty"`
	// AutoInstall defines whether Cloudflare will inject the JS snippet automatically for orange-clouded sites.
	AutoInstall bool `json:"auto_install"`
}

// CreateWebAnalyticsSite creates a new Web Analytics Site for an Account.
//
// API reference: https://api.cloudflare.com/#web-analytics-create-site
func (api *API) CreateWebAnalyticsSite(ctx context.Context, rc *ResourceContainer, params CreateWebAnalyticsSiteParams) (*WebAnalyticsSite, error) {
	if params.Host == "" && params.ZoneTag == "" {
		return nil, ErrMissingWebAnalyticsSiteHost
	}
	uri := fmt.Sprintf("/accounts/%s/rum/site_info", rc.Identifier)
	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, params)
	if err != nil {
		return nil, err
	}
	var r WebAnalyticsSiteResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return &r.Result, nil
}

type ListWebAnalyticsSitesParams struct {
	// Page offset for pagination.
	Page int
	// Items per page limit for pagination.
	PerPage int
	// Property to order Sites by, "host" or "created".
	OrderBy string
}

// ListWebAnalyticsSites returns all Web Analytics Sites for an Account.
//
// API reference: https://api.cloudflare.com/#web-analytics-list-sites
func (api *API) ListWebAnalyticsSites(ctx context.Context, rc *ResourceContainer, params ListWebAnalyticsSitesParams) ([]WebAnalyticsSite, *ResultInfo, error) {
	v := url.Values{}
	if params.Page > 0 {
		v.Set("page", strconv.Itoa(params.Page))
	}
	if params.PerPage > 0 {
		v.Set("per_page", strconv.Itoa(params.PerPage))
	}
	if params.OrderBy != "" {
		v.Set("order_by", params.OrderBy)
	}
	uri := fmt.Sprintf("/accounts/%s/rum/site_info/list?"+v.Encode(), rc.Identifier)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil, err
	}
	var r WebAnalyticsSitesResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return r.Result, &r.ResultInfo, nil
}

type WebAnalyticsSiteParams struct {
	SiteTag string
}

// WebAnalyticsSite fetches detail about one Web Analytics Site for an Account.
//
// API reference: https://api.cloudflare.com/#web-analytics-get-site
func (api *API) WebAnalyticsSite(ctx context.Context, rc *ResourceContainer, params WebAnalyticsSiteParams) (*WebAnalyticsSite, error) {
	if params.SiteTag == "" {
		return nil, ErrMissingWebAnalyticsSiteTag
	}
	uri := fmt.Sprintf("/accounts/%s/rum/site_info/%s", rc.Identifier, params.SiteTag)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	var r WebAnalyticsSiteResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return &r.Result, nil
}

type UpdateWebAnalyticsSiteParams struct {
	SiteTag string `json:"-"`
	// Host is the host to measure traffic for.
	Host string `json:"host,omitempty"`
	// ZoneTag is the zone tag to measure traffic for.
	ZoneTag string `json:"zone_tag,omitempty"`
	// AutoInstall defines whether Cloudflare will inject the JS snippet automatically for orange-clouded sites.
	AutoInstall bool `json:"auto_install"`
}

// UpdateWebAnalyticsSite lets you replace a Web Analytics Site for an Account.
//
// API reference: https://api.cloudflare.com/#web-analytics-update-site
func (api *API) UpdateWebAnalyticsSite(ctx context.Context, rc *ResourceContainer, params UpdateWebAnalyticsSiteParams) (*WebAnalyticsSite, error) {
	if params.SiteTag == "" {
		return nil, ErrMissingWebAnalyticsSiteTag
	}
	uri := fmt.Sprintf("/accounts/%s/rum/site_info/%s", rc.Identifier, params.SiteTag)
	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params)
	if err != nil {
		return nil, err
	}
	var r WebAnalyticsSiteResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return &r.Result, nil
}

type DeleteWebAnalyticsSiteParams struct {
	SiteTag string
}

// DeleteWebAnalyticsSite deletes a Web Analytics Site for an Account.
//
// API reference: https://api.cloudflare.com/#web-analytics-delete-site
func (api *API) DeleteWebAnalyticsSite(ctx context.Context, rc *ResourceContainer, params DeleteWebAnalyticsSiteParams) (*string, error) {
	if params.SiteTag == "" {
		return nil, ErrMissingWebAnalyticsSiteTag
	}
	uri := fmt.Sprintf("/accounts/%s/rum/site_info/%s", rc.Identifier, params.SiteTag)
	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	var r WebAnalyticsSiteResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return &r.Result.SiteTag, nil
}

type CreateWebAnalyticsRuleParams struct {
	RulesetID string
	Rule      WebAnalyticsRule
}

// CreateWebAnalyticsRule creates a new Web Analytics Rule for a WebAnalyticsRuleset.
//
// API reference: https://api.cloudflare.com/#web-analytics-create-rule
func (api *API) CreateWebAnalyticsRule(ctx context.Context, rc *ResourceContainer, params CreateWebAnalyticsRuleParams) (*WebAnalyticsRule, error) {
	if params.RulesetID == "" {
		return nil, ErrMissingWebAnalyticsRulesetID
	}
	uri := fmt.Sprintf("/accounts/%s/rum/v2/%s/rule", rc.Identifier, params.RulesetID)
	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, params.Rule)
	if err != nil {
		return nil, err
	}
	var r WebAnalyticsRuleResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return &r.Result, nil
}

type WebAnalyticsRuleParams struct {
	RulesetID string
	RuleID    string
}

// WebAnalyticsRule fetches detail about one Web Analytics Rule for a WebAnalyticsRuleset.
//
// API reference: https://api.cloudflare.com/#web-analytics-get-rule
func (api *API) WebAnalyticsRule(ctx context.Context, rc *ResourceContainer, params WebAnalyticsRuleParams) (*WebAnalyticsRule, error) {
	if params.RulesetID == "" {
		return nil, ErrMissingWebAnalyticsRulesetID
	}
	if params.RuleID == "" {
		return nil, ErrMissingWebAnalyticsRuleID
	}
	uri := fmt.Sprintf("/accounts/%s/rum/v2/%s/rule/%s", rc.Identifier, params.RulesetID, params.RuleID)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	var r WebAnalyticsRuleResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return &r.Result, nil
}

type ListWebAnalyticsRulesParams struct {
	RulesetID string
}

// ListWebAnalyticsRules fetches all Web Analytics Rules for a WebAnalyticsRuleset.
//
// API reference: https://api.cloudflare.com/#web-analytics-list-rules
func (api *API) ListWebAnalyticsRules(ctx context.Context, rc *ResourceContainer, params ListWebAnalyticsRulesParams) (*WebAnalyticsRulesetRules, error) {
	if params.RulesetID == "" {
		return nil, ErrMissingWebAnalyticsRulesetID
	}
	uri := fmt.Sprintf("/accounts/%s/rum/v2/%s/rules", rc.Identifier, params.RulesetID)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	var r WebAnalyticsRulesResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return &r.Result, nil
}

type DeleteWebAnalyticsRuleParams struct {
	RulesetID string
	RuleID    string
}

// DeleteWebAnalyticsRule deletes a Web Analytics Rule for a WebAnalyticsRuleset.
//
// API reference: https://api.cloudflare.com/#web-analytics-delete-rule
func (api *API) DeleteWebAnalyticsRule(ctx context.Context, rc *ResourceContainer, params DeleteWebAnalyticsRuleParams) (*string, error) {
	if params.RulesetID == "" {
		return nil, ErrMissingWebAnalyticsRulesetID
	}
	if params.RuleID == "" {
		return nil, ErrMissingWebAnalyticsRuleID
	}
	uri := fmt.Sprintf("/accounts/%s/rum/v2/%s/rule/%s", rc.Identifier, params.RulesetID, params.RuleID)
	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	var r WebAnalyticsIDResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return &r.Result, nil
}

type UpdateWebAnalyticsRuleParams struct {
	RulesetID string
	Rule      WebAnalyticsRule
}

// UpdateWebAnalyticsRule lets you replace a Web Analytics Rule for a WebAnalyticsRuleset.
//
// API reference: https://api.cloudflare.com/#web-analytics-update-rule
func (api *API) UpdateWebAnalyticsRule(ctx context.Context, rc *ResourceContainer, params UpdateWebAnalyticsRuleParams) (*WebAnalyticsRule, error) {
	if params.RulesetID == "" {
		return nil, ErrMissingWebAnalyticsRulesetID
	}
	if params.Rule.ID == "" {
		return nil, ErrMissingWebAnalyticsRuleID
	}
	uri := fmt.Sprintf("/accounts/%s/rum/v2/%s/rule/%s", rc.Identifier, params.RulesetID, params.Rule.ID)
	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, params.Rule)
	if err != nil {
		return nil, err
	}
	var r WebAnalyticsRuleResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return &r.Result, nil
}

type ModifyWebAnalyticsRulesParams struct {
	RulesetID string
	// Rules is a list of rules to create if no id is specified, or modify if an existing id is specified.
	Rules []WebAnalyticsRule
	// DeleteRules is a list of rule ids to be deleted
	DeleteRules []string
}

// ModifyWebAnalyticsRules lets you create, updated, or delete multiple Web Analytics Rule for a WebAnalyticsRuleset.
//
// API reference: https://api.cloudflare.com/#web-analytics-modify-rules
func (api *API) ModifyWebAnalyticsRules(ctx context.Context, rc *ResourceContainer, params ModifyWebAnalyticsRulesParams) (*WebAnalyticsRulesetRules, error) {
	if params.RulesetID == "" {
		return nil, ErrMissingWebAnalyticsRulesetID
	}
	uri := fmt.Sprintf("/accounts/%s/rum/v2/%s/rules", rc.Identifier, params.RulesetID)
	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, struct {
		Rules       []WebAnalyticsRule `json:"rules"`
		DeleteRules []string           `json:"deleteRules"`
	}{
		Rules:       params.Rules,
		DeleteRules: params.DeleteRules,
	})
	if err != nil {
		return nil, err
	}
	var r WebAnalyticsRulesResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}
	return &r.Result, nil
}
