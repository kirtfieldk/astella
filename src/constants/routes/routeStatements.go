package routes

var GET_EVENT_BY_ID string = "/event/:id"
var GET_EVENT_BY_CITY string = "/event/city"
var ADD_USER_TO_EVENT string = "/event/:userId/:eventId"
var POST_MESSAGE_TO_EVENT string = "/message/:userId/:eventId"
var GET_MESSAGE_IN_EVENT string = "/message/event/:eventId/:userId"
var LIKE_MESSAGE_IN_EVENT string = "/message/event/upvote/:eventId/:userId/:messageId"
var GET_USRS_LIKE_MESSAGE string = "/message/event/whoupvote/:eventId/:userId/:messageId"
