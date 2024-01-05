// Package invoice 电子发票
package invoice

import (
	"context"

	"github.com/huimingz/wechatgo/wecom"
)

const (
	urlGetInfo           string = "/cgi-bin/card/invoice/reimburse/getinvoiceinfo"
	urlGetInfoBatch      string = "/cgi-bin/card/invoice/reimburse/getinvoiceinfobatch"
	urlUpdateStatus      string = "/cgi-bin/card/invoice/reimburse/updateinvoicestatus"
	urlUpdateStatusBatch string = "/cgi-bin/card/invoice/reimburse/updatestatusbatch"
)

type Product struct {
	Name  string `json:"name"`  // 项目（商品）名称
	Num   int    `json:"num"`   // 项目数量
	Unit  string `json:"unit"`  // 项目单位
	Fee   int    `json:"fee"`   // 发票加税合计金额，以分为单位
	Price int    `json:"price"` // 单价，以分为单位
}

type User struct {
	Fee                   int       `json:"fee"`                      // 发票加税合计金额，以分为单位
	Title                 string    `json:"title"`                    // 发票的抬头
	BillingTime           int       `json:"billing_time"`             // 开票时间，为十位时间戳
	BillingNo             string    `json:"billing_no"`               // 发票代码
	BillingCode           string    `json:"billing_code"`             // 发票号码
	Info                  []Product `json:"info"`                     // 商品信息结构
	FeeWithoutTax         int       `json:"fee_without_tax"`          // 不含税金额，以分为单位
	Tax                   int       `json:"tax"`                      // 分
	Detail                string    `json:"detail"`                   // 发票详情，一般描述的是发票的使用说明
	PdfUrl                string    `json:"pdf_url"`                  // 这张发票对应的PDF_URL
	TripPdfUrl            string    `json:"trip_pdf_url"`             // 其它消费凭证附件对应的URL，如行程单、水单等
	CheckCode             string    `json:"check_code"`               // 校验码
	BuyerNumber           string    `json:"buyer_number"`             // 购买方纳税人识别号
	BuyerAddressAndPhone  string    `json:"buyer_address_and_phone"`  // 购买方地址、电话
	BuyerBankAccount      string    `json:"buyer_bank_account"`       // 购买方开户行及账号
	SellerNumber          string    `json:"seller_number"`            // 销售方纳税人识别号
	SellerAddressAndPhone string    `json:"seller_address_and_phone"` // 销售方地址、电话
	SellerBankAccount     string    `json:"seller_bank_account"`      // 销售方开户行及账号
	Remarks               string    `json:"remarks"`                  // 备注
	Cashier               string    `json:"cashier"`                  // 收款人，发票左下角处
	Maker                 string    `json:"maker"`                    // 开票人，发票有下角处
	ReimburseStatus       string    `json:"reimburse_status"`         // 发报销状态INVOICE_REIMBURSE_INIT：发票初始状态，未锁定；INVOICE_REIMBURSE_LOCK：发票已锁定；INVOICE_REIMBURSE_CLOSURE：发票已核销
}

type Info struct {
	CardId    string `json:"card_id"`    // 发票id
	BeginTime int    `json:"begin_time"` // 发票的有效期起始时间
	EndTime   int    `json:"end_time"`   // 发票的有效期截止时间
	OpenId    string `json:"openid"`     // 用户标识
	Type      string `json:"type"`       // 发票类型，如广东增值税普通发票
	Payee     string `json:"payee"`      // 发票的收款方
	Detail    string `json:"detail"`     // 发票详情
	UserInfo  User   `json:"user_info"`  // 发票的用户信息
}

type CardInfo struct {
	CardId      string `json:"card_id"`      // 发票id
	EncryptCode string `json:"encrypt_code"` // 加密code
}

type WechatInvoice struct {
	Client *wecom.Client
}

func NewWechatInvoice(client *wecom.Client) *WechatInvoice {
	return &WechatInvoice{client}
}

func (w WechatInvoice) GetInfo(ctx context.Context, cardId, encryptCode string) (*Info, error) {
	data := struct {
		CardId      string `json:"card_id"`      // 发票id
		EncryptCode string `json:"encrypt_code"` // 加密code
	}{
		CardId:      cardId,
		EncryptCode: encryptCode,
	}

	out := Info{}
	err := w.Client.Post(ctx, urlGetInfo, nil, data, nil, &out)
	return &out, err
}

func (w WechatInvoice) GetInfoBatch(ctx context.Context, items []CardInfo) (invList []Info, err error) {
	data := struct {
		ItemList []CardInfo `json:"item_list"`
	}{
		ItemList: items,
	}

	out := struct {
		ItemList []Info `json:"item_list"`
	}{}

	err = w.Client.Post(ctx, urlGetInfoBatch, nil, data, nil, &out)
	invList = out.ItemList
	return
}

func (w WechatInvoice) UpdateStatus(ctx context.Context, cardId, encryptCode, reimburseStatus string) error {
	// TODO: 检查reimburseStatus是否有效
	data := struct {
		// 发票id
		CardId string `json:"card_id"`

		// 加密code
		EncryptCode string `json:"encrypt_code"`

		// 发报销状态
		// INVOICE_REIMBURSE_INIT：发票初始状态，未锁定；
		// INVOICE_REIMBURSE_LOCK：发票已锁定，无法重复提交报销;
		// INVOICE_REIMBURSE_CLOSURE:发票已核销，从用户卡包中移除
		ReimburseStatus string `json:"reimburse_status"`
	}{
		CardId:          cardId,
		EncryptCode:     encryptCode,
		ReimburseStatus: reimburseStatus,
	}

	return w.Client.Post(ctx, urlUpdateStatus, nil, data, nil, nil)
}

// 批量更新发票状态
//
// 接口说明：发票平台可以通过该接口对某个成员的一批发票进行锁定、解锁和报销操作。
// 注意，报销状态为不可逆状态，请开发者慎重调用。
//
// 企业微信API：https://work.weixin.qq.com/api/doc#90000/90135/90286
func (w WechatInvoice) UpdateStatusBatch(ctx context.Context, openId, reimburseStatus string, invoiceList []CardInfo) error {
	data := struct {
		// 用户openid，可用“userid与openid互换接口”获取
		OpenId string `json:"openid"`

		// 发报销状态
		// INVOICE_REIMBURSE_INIT：发票初始状态，未锁定；
		// INVOICE_REIMBURSE_LOCK：发票已锁定，无法重复提交报销;
		// INVOICE_REIMBURSE_CLOSURE:发票已核销，从用户卡包中移除
		ReimburseStatus string `json:"reimburse_status"`

		// 发票列表，必须全部属于同一个openid
		InvoiceList []CardInfo `json:"invoice_list"`
	}{
		OpenId:          openId,
		ReimburseStatus: reimburseStatus,
		InvoiceList:     invoiceList,
	}

	return w.Client.Post(ctx, urlUpdateStatusBatch, nil, data, nil, nil)
}
