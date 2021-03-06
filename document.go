package esssdk

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

// DocType resprents doc type.
type DocType int

// DocType enum
const (
	PDFDocType DocType = 2
	OFDDocType DocType = 3
)

// Document represents an ess document.
type Document struct {
	DocName          string  `json:"docName"`
	DocType          DocType `json:"docType"`
	DocContentBase64 string  `json:"docContentBase64"`
	UserID           string  `json:"userId"`
	UUID             string  `json:"uuid"`
}

// CreateDocument creates an ess document.
func (c *Client) CreateDocument(doc *Document) (err error) {
	url := fmt.Sprintf("http://%s/ess/api/user/doc/create", c.opts.Addr)

	formData := map[string]string{
		"docName": doc.DocName,
		"docType": fmt.Sprintf("%v", doc.DocType),
		"userId":  doc.UserID,
		"uuid":    doc.UUID,
	}

	var fileName string
	switch doc.DocType {
	case PDFDocType:
		fileName = fmt.Sprintf("%s.pdf", doc.DocName)
	case OFDDocType:
		fileName = fmt.Sprintf("%s.ofd", doc.DocName)
	default:
		fileName = doc.DocName
	}

	fileContent, err := base64.StdEncoding.DecodeString(doc.DocContentBase64)
	if err != nil {
		err = fmt.Errorf("文档内容格式无效: %v", err)
		return
	}
	fileReader := bytes.NewReader(fileContent)

	data, err := c.postFileForm(url, formData, "file", fileName, fileReader)
	if err != nil {
		return
	}

	v, ok := data.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("response data invalid: %v", data)
		return
	}

	tv, found := v["id"]
	if !found {
		err = fmt.Errorf("response data invalid: %v", v)
		return
	}

	docID, ok := tv.(string)
	if !ok {
		err = fmt.Errorf("response data invalid: %v", tv)
		return
	}

	if docID == "" {
		err = fmt.Errorf("response data invalid: id is empty")
		return
	}

	doc.UUID = docID

	return
}

// KeywordSignInfo represents the keyword signing info to an ess document.
type KeywordSignInfo struct {
	DocID   string `json:"docId"`
	SealID  string `json:"sealId"`
	Keyword string `json:"keyword"`
	Scope   int    `json:"scope"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
	Zoom    int    `json:"zoom"`
}

// SignDocumentForKeyword signs an ess document for specified keyword.
func (c *Client) SignDocumentForKeyword(si *KeywordSignInfo) (signedDocURL string, err error) {
	url := fmt.Sprintf("http://%s/ess/api/user/doc/sign/keyword", c.opts.Addr)

	params := map[string]string{
		"docId":   si.DocID,
		"sealId":  si.SealID,
		"keyword": si.Keyword,
		"scope":   fmt.Sprintf("%v", si.Scope),
		"start":   fmt.Sprintf("%v", si.Start),
		"end":     fmt.Sprintf("%v", si.End),
		"zoom":    fmt.Sprintf("%v", si.Zoom),
	}
	data, err := c.postParams(url, params)
	if err != nil {
		return
	}

	v, ok := data.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("response data invalid: %v", data)
		return
	}

	tv, found := v["id"]
	if !found {
		err = fmt.Errorf("response data invalid: %v", v)
		return
	}

	docID, ok := tv.(string)
	if !ok {
		err = fmt.Errorf("response data invalid: %v", tv)
		return
	}

	if docID != si.DocID {
		err = fmt.Errorf("response data invalid: id is unmatched")
		return
	}

	tv, found = v["url"]
	if !found {
		err = fmt.Errorf("response data invalid: %v", v)
		return
	}

	signedDocURL, ok = tv.(string)
	if !ok {
		err = fmt.Errorf("response data invalid: %v", tv)
		return
	}

	return
}

// PositionSignInfo represents the position signing info to an ess document.
type PositionSignInfo struct {
	DocID      string `json:"docId"`
	SealID     string `json:"sealId"`
	PageNumber string `json:"pageNumber"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
	Zoom       int    `json:"zoom"`
	Reason     string `json:"reason"`
	Remark     string `json:"remark"`
}

// SignDocumentForPosition signs an ess document for specified postion.
func (c *Client) SignDocumentForPosition(si *PositionSignInfo) (signedDocURL string, err error) {
	url := fmt.Sprintf("http://%s/ess/api/user/doc/sign/position", c.opts.Addr)

	params := map[string]string{
		"docId":               si.DocID,
		"remark":              si.Remark,
		"signs[0].pageNumber": si.PageNumber,
		"signs[0].sealId":     si.SealID,
		"signs[0].x":          fmt.Sprintf("%v", si.X),
		"signs[0].y":          fmt.Sprintf("%v", si.Y),
		"signs[0].zoom":       fmt.Sprintf("%v", si.Zoom),
		"signs[0].reason":     si.Reason,
	}
	data, err := c.postParams(url, params)
	if err != nil {
		return
	}

	v, ok := data.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("response data invalid: %v", data)
		return
	}

	tv, found := v["id"]
	if !found {
		err = fmt.Errorf("response data invalid: %v", v)
		return
	}

	docID, ok := tv.(string)
	if !ok {
		err = fmt.Errorf("response data invalid: %v", tv)
		return
	}

	if docID != si.DocID {
		err = fmt.Errorf("response data invalid: id is unmatched")
		return
	}

	tv, found = v["url"]
	if !found {
		err = fmt.Errorf("response data invalid: %v", v)
		return
	}

	signedDocURL, ok = tv.(string)
	if !ok {
		err = fmt.Errorf("response data invalid: %v", tv)
		return
	}

	return
}
