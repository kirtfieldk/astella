package routes

const GET_EVENT_BY_ID string = "/api/v1/event/:id"
const CREATE_EVENT string = "/api/v1/event"
const PIN_MESSAGE string = "/api/v1/event/pin/message/:eventId/:userId/:messageId"
const GET_PIN_MESSAGE string = "/api/v1/event/pin/message/:eventId/:userId"
const GET_EVENT_BY_CITY string = "/api/v1/event/city"
const ADD_USER_TO_EVENT string = "/api/v1/event/:eventId/:userId"
const POST_MESSAGE_TO_EVENT string = "/api/v1/message/post/:eventId/:userId"
const GET_MESSAGE_IN_EVENT string = "/api/v1/message/event/:eventId/:userId"
const LIKE_MESSAGE_IN_EVENT string = "/api/v1/message/event/upvote/:eventId/:userId/:messageId"
const UNLIKE_MESSAGE_IN_EVENT string = "/api/v1/message/event/downvote/:eventId/:userId/:messageId"
const GET_EVENTS_MEMBER_OF string = `/api/v1/user/event/member/:userId`
const GET_EVENTS_MEMBERS string = `/api/v1/member/user/event/:eventId`

const GET_USRS_LIKE_MESSAGE string = "/api/v1/message/event/whoupvote/:eventId/:userId/:messageId"
