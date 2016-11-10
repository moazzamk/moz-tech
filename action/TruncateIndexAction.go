package action

import "gopkg.in/olivere/elastic.v3"

/**
 * Truncate index action
 *
 * Deletes and re-creates an ES index
 */

type TruncateIndexAction struct {
	client *elastic.Client
}

func NewTruncateIndexAction(client *elastic.Client) *TruncateIndexAction {
	return &TruncateIndexAction{client: client}
}

func (r *TruncateIndexAction) Run(name string) {
	r.client.DeleteIndex(name).Do()
	r.client.CreateIndex(name).Do()
}
