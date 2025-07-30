// Package i18n provides centralized localization for consistent UI messaging.
package i18n

const (
	CommonWait        = "Пожалуйста, подождите..."
	CommonError       = "Ошибка: %s\nНажмите Enter для продолжения..."
	CommonPressEnter  = "Нажмите ENTER для подтверждения"
	CommonPressESC    = "Нажмите ESC для отмены"
	CommonPressAnyKey = "Нажмите любую клавишу..."

	InputDataPrompt     = "Введите данные:\n\n"
	InputSavePathPrompt = "Введите путь сохранения файла:\n\n>%s\n\n" + CommonPressEnter

	VaultTitle                 = "GophKeeper"
	VaultListTitleNameSingular = "Объект"
	VaultListTitleNamePlural   = "Объектов"
	VaultObjectTitle           = "Объект: %s"
	VaultTypeTitle             = "Тип: %s"
	VaultDescTitle             = "Описание: %s"
	VaultUpdatedTitle          = "Дата последнего обновления: %s\n\n"
	VaultActions               = "Нажмите ENTER для загрузки данных...\n" +
		"Нажмите DEL для удаления данных...\n" +
		"Нажмите INS для редактирования данных...\n" +
		"Нажмите ESC для возврата к списку..."
	VaultUpdateHelp  = "обновить"
	VaultAddItemHelp = "добавить"

	AuthLoginTitle     = "Вход в GophKeeper"
	AuthRegisterTitle  = "Регистрация"
	AuthUsernameLabel  = "Логин: %s"
	AuthPasswordLabel  = "Пароль: %s"
	AuthLoginButton    = "Вход"
	AuthRegisterButton = "Регистрация"
	AuthTabHint        = "Нажмите Enter для подтверждения, Tab для переключения"

	AddSelectPrompt   = "Выберите тип хранимой информации:\n\n"
	AddChoiceTemplate = "%s %s\n"

	InputName = "Название"
	InputInfo = "Информация"

	BinInputFilePath = "Путь к файлу"
	BinInputSuccess  = "Файл успешно сохранен по пути:"

	CardInputNumber     = "Номер"
	CardInputHolderName = "Имя владельца"
	CardInputExpiryDate = "Срок действия"
	CardInputCVV        = "CVV"

	LogPassInputLogin    = "Логин"
	LogPassInputPassword = "Пароль"

	TextInputContent = "Ваш текст"

	AddPassword = "Логин/Пароль"
	AddCard     = "Банковская карта"
	AddText     = "Текст"
	AddBinary   = "Бинарный файл"
)
