package helper

import "encoding/json"

const (
	TaskMessageGuestBookingSubmit        = float64(20001)
	TaskMessageGuestInviteNotify         = float64(20002)
	TaskMessageGuestCheckResultNotify    = float64(20003)
	TaskMessageGuestVisitNotify          = float64(20004)
	TaskMessageGuestQuestionnaire        = float64(20005)
	TaskMessageGuestLeaveNotify          = float64(20006)
	TaskMessageStaffBookingNotify        = float64(30001)
	TaskMessageStaffInviteSubmit         = float64(30002)
	TaskMessageStaffCheckTodoNotify      = float64(30003)
	TaskMessageStaffCheckResultNotify    = float64(30004)
	TaskMessageStaffCheckRecipientNotify = float64(30005)
	TaskMessageStaffVisitNotify          = float64(30006)
	TaskMessageStaffLeaveNotify          = float64(30007)
	TaskDevicePassage                    = float64(40001)
	TaskWorkFlow                         = float64(50001)
)

type message struct {
	Tenant string                 `json:"tenant"`
	Uuid   string                 `json:"uuid"`
	Cmd    float64                `json:"cmd"`
	Data   map[string]interface{} `json:"data"`
}

func buildDeviceTask(tenant string, pk string, cmd float64, data map[string]interface{}) (string, []byte) {
	payload, _ := json.Marshal(message{
		Tenant: tenant,
		Uuid:   pk,
		Cmd:    cmd,
		Data:   data,
	})
	return "device", payload
}

func buildMessageTask(tenant, uuid string, cmd float64, data map[string]interface{}) (string, []byte) {
	payload, _ := json.Marshal(message{
		Tenant: tenant,
		Uuid:   uuid,
		Cmd:    cmd,
		Data:   data,
	})
	return "message", payload
}

func buildWorkFlowTask(tenant, pk string, cmd float64, data map[string]interface{}) (string, []byte) {
	payload, _ := json.Marshal(message{
		Tenant: tenant,
		Uuid:   pk,
		Cmd:    cmd,
		Data:   data,
	})
	return "workflow", payload
}

func BWorkFlowSubmit(tenant, pk, contact, sponsor string, types float64) (string, []byte) {
	return buildWorkFlowTask(tenant, pk, TaskWorkFlow, map[string]interface{}{
		"type":    types,
		"contact": contact,
		"sponsor": sponsor,
	})
}

func BDBookingAndInviteData(tenant, pk, name, cardNo, visitorCompany, visitorSex, qrCode, carNumber, time, endTime, imgUrl, keyWord, visitorPhone, visitedPhone, visitedName,
	visitedCompany, visitedSex, reason, childCompany, childCompanyArea, ticket string, visitedUnit []string) (string, []byte) {
	return buildDeviceTask(tenant, pk, TaskDevicePassage, map[string]interface{}{
		"personId":         pk,
		"name":             name,
		"cardNo":           cardNo,
		"visitorCompany":   visitorCompany,
		"visitorSex":       visitorSex,
		"qrCode":           qrCode,
		"carNumber":        carNumber,
		"time":             time,
		"endTime":          endTime,
		"imgUrl":           imgUrl,
		"keyWord":          keyWord,
		"visitorPhone":     visitorPhone,
		"visitedPhone":     visitedPhone,
		"visitedName":      visitedName,
		"visitedCompany":   visitedCompany,
		"visitedSex":       visitedSex,
		"visitedUnit":      visitedUnit,
		"reason":           reason,
		"childCompany":     childCompany,
		"childCompanyArea": childCompanyArea,
		"ticket":           ticket,
	})
}

func BTMGuestBookingSubmit(tenant, uuid, pk, staffName, reason, startTime, guestName, guestPhone, remark string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageGuestBookingSubmit, map[string]interface{}{
		"pk":         pk,
		"staffName":  staffName,
		"reason":     reason,
		"startTime":  startTime,
		"guestName":  guestName,
		"guestPhone": guestPhone,
		"remark":     remark,
	})
}

func BTMGuestInviteNotify(tenant, uuid, pk, staffName, reason, startTime, guestName, remark string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageGuestInviteNotify, map[string]interface{}{
		"pk":        pk,
		"staffName": staffName,
		"reason":    reason,
		"startTime": startTime,
		"guestName": guestName,
		"remark":    remark,
	})
}

func BTMGuestCheckResultNotify(tenant, uuid, pk, staffName, status, startTime, endTime, guestName, remark string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageGuestCheckResultNotify, map[string]interface{}{
		"pk":        pk,
		"staffName": staffName,
		"startTime": startTime,
		"endTime":   endTime,
		"status":    status,
		"guestName": guestName,
		"remark":    remark,
	})
}

func BTMGuestVisitNotify(tenant, uuid, pk, staffName, staffPhone, reason string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageGuestVisitNotify, map[string]interface{}{
		"pk":         pk,
		"staffName":  staffName,
		"staffPhone": staffPhone,
		"reason":     reason,
	})
}

func BTMGuestQuestionnaire(tenant, uuid, pk, staffName string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageGuestQuestionnaire, map[string]interface{}{
		"pk":        pk,
		"staffName": staffName,
	})
}

func BTMGuestLeaveNotify(tenant, uuid, guestName, leaveTime string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageGuestLeaveNotify, map[string]interface{}{
		"guestName": guestName,
		"leaveTime": leaveTime,
	})
}

func BTMStaffBookingNotify(tenant, uuid, pk, staffName, reason, startTime, guestName, guestPhone, remark string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageStaffBookingNotify, map[string]interface{}{
		"pk":         pk,
		"staffName":  staffName,
		"reason":     reason,
		"startTime":  startTime,
		"guestName":  guestName,
		"guestPhone": guestPhone,
		"remark":     remark,
	})
}

func BTMStaffInviteSubmit(tenant, uuid, pk, staffName, reason, startTime, guestName, guestPhone, remark string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageStaffInviteSubmit, map[string]interface{}{
		"pk":         pk,
		"staffName":  staffName,
		"reason":     reason,
		"startTime":  startTime,
		"guestName":  guestName,
		"guestPhone": guestPhone,
		"remark":     remark,
	})
}

func BTMStaffCheckTodoNotify(tenant, uuid, pk, staffName, reason, startTime, guestName, guestPhone, remark string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageStaffCheckTodoNotify, map[string]interface{}{
		"pk":         pk,
		"staffName":  staffName,
		"reason":     reason,
		"startTime":  startTime,
		"guestName":  guestName,
		"guestPhone": guestPhone,
		"remark":     remark,
	})
}

func BTMStaffCheckResultNotify(tenant, uuid, pk, staffName, status, startTime, guestName, remark string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageStaffCheckResultNotify, map[string]interface{}{
		"pk":        pk,
		"staffName": staffName,
		"status":    status,
		"startTime": startTime,
		"guestName": guestName,
		"remark":    remark,
	})
}

func BTMStaffCheckRecipientNotify(tenant, uuid, pk, staffName, startTime, guestName string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageStaffCheckRecipientNotify, map[string]interface{}{
		"pk":        pk,
		"staffName": staffName,
		"startTime": startTime,
		"guestName": guestName,
	})
}

func BTMStaffVisitNotify(tenant, uuid, pk, guestName, guestPhone, visitTime, deviceName string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageStaffVisitNotify, map[string]interface{}{
		"pk":         pk,
		"guestName":  guestName,
		"guestPhone": guestPhone,
		"visitTime":  visitTime,
		"deviceName": deviceName,
	})
}

func BTMStaffLeaveNotify(tenant, uuid, pk, guestName, leaveTime, deviceName string) (string, []byte) {
	return buildMessageTask(tenant, uuid, TaskMessageStaffLeaveNotify, map[string]interface{}{
		"pk":         pk,
		"guestName":  guestName,
		"leaveTime":  leaveTime,
		"deviceName": deviceName,
	})
}
