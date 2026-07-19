package redisview

type SessionRedisView struct {
	UUID      string `redis:"uuid"`
	UserUUID  string `redis:"user_uuid"`
	Login     string `redis:"login"`
	CreatedAt string `redis:"created_at"` // ← см. развилку ниже
	ExpiresAt string `redis:"expires_at"`
}
