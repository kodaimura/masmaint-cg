package jwt

import (
	"time"
)

const COOKIE_KEY_JWT string = "access_token"
const CONTEXT_KEY_PAYLOAD string = "payload"
const JWT_EXPIRES time.Duration = 30 * 24 * 60 * 60