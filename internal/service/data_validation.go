package service

import (
	"errors"
	"net/http"
	"regexp"
)

func validateEmail(email string) error {
	// Проверка электронной почты (пример: должна соответствовать формату "example@example.com")
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

//func validatePassword(password string) error {
//	// Проверка пароля (пример: длина не менее 8 символов, содержит цифры и буквы)
//	pattern := `^(?=.*\d)(?=.*[a-zA-Z]).{8,}$`
//	matched, err := regexp.MatchString(pattern, password)
//	if err != nil {
//		fmt.Printf(" what %v\n", err)
//		return err
//	}
//	if !matched {
//		return errors.New("invalid password format")
//	}
//	return nil
//}

func validateGender(genUser string) bool {
	var genderArray [3]string = [3]string{"female", "male", "other"}

	valid := false
	for _, gender := range genderArray {
		if genUser == gender {
			valid = true
			break
		}
	}

	return valid
}

func validateProfilePhoto(photoURL *string) error {
	// Проверка фотографии профиля
	if *photoURL == "" {
		// Если фотография не была загружена, устанавливаем значение по умолчанию
		*photoURL = "https://anavara.com/wp-content/uploads/2020/05/blank-profile-pic.png"
		return nil
	}
	// Проверка типа файла (допустимые типы: .jpg, .jpeg, .png)
	allowedTypes := []string{".jpg", ".jpeg", ".png"}

	fileType := http.DetectContentType([]byte(*photoURL))
	allowed := false
	for _, t := range allowedTypes {
		if fileType == "image/"+t[1:] {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("invalid file type")
	}

	// Здесь можно добавить дополнительные проверки, например, максимального размера файла

	return nil
}

func validatePhone(phone string) error {
	// Проверка номера телефона (пример: должен соответствовать формату "+7XXXXXXXXXX")
	pattern := `^\+7\d{10}$`
	matched, err := regexp.MatchString(pattern, phone)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid phone number format")
	}
	return nil
}

func validateUsername(nickname string) error {
	// Проверка никнейма (пример: не должен быть пустым, длина от 3 до 20 символов)
	if len(nickname) < 3 || len(nickname) > 20 {
		return errors.New("invalid nickname length")
	}
	return nil
}
