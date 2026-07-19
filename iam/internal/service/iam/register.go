package iam

// Реализовать (неделя 6): Register (контракт).
//   валидация: пустой логин → ErrInvalidLogin; пароль < 8 → ErrWeakPassword
//   bcrypt.GenerateFromPassword(DefaultCost) → UserRepository.Create → user_uuid
