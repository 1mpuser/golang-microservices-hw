package iam

// Реализовать (неделя 6): Login (контракт).
//   пустые креды → ErrEmptyCredential
//   UserRepository.GetByLogin + bcrypt.CompareHashAndPassword
//   любое несовпадение → ErrInvalidCredentials (не раскрываем, существует ли логин)
//   SessionRepository.Create (TTL 24h) → session_uuid
