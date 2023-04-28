package routes

var GET_EVENT_BY_ID string = "/api/v1/get/event/:id"
var CREATE_EVENT string = "/api/v1/event"
var GET_EVENT_BY_CITY string = "/api/v1/event/city"
var ADD_USER_TO_EVENT string = "/api/v1/event/:userId/:eventId"
var POST_MESSAGE_TO_EVENT string = "/api/v1/message/post/:eventId/:userId"
var GET_MESSAGE_IN_EVENT string = "/api/v1/message/event/:eventId/:userId"
var LIKE_MESSAGE_IN_EVENT string = "/api/v1/message/event/upvote/:eventId/:userId/:messageId"
var GET_EVENTS_MEMBER_OF string = `/api/v1/event/member/user/:userId`
var GET_USRS_LIKE_MESSAGE string = "/api/v1/message/event/whoupvote/:eventId/:userId/:messageId"
