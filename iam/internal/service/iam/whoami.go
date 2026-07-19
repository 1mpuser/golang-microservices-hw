package iam

// Реализовать (неделя 6): Whoami (контракт).
//   пустой session_uuid → ErrEmptySessionID
//   SessionRepository.Get → нет/истекла → ErrSessionNotFound
//   вернуть session + user (UserRepository.GetByUUID)
