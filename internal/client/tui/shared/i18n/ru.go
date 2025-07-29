package i18n

const (
	// Общие фразы
	CommonWait        = "Пожалуйста, подождите..."
	CommonError       = "Ошибка: %s\nНажмите Enter для продолжения..."
	CommonPressEnter  = "Нажмите ENTER для подтверждения"
	CommonPressESC    = "Нажмите ESC для отмены"
	CommonPressAnyKey = "Нажмите любую клавишу..."

	// Формы ввода
	InputDataPrompt     = "Введите данные:\n\n"
	InputSavePathPrompt = "Введите путь сохранения файла:\n\n>%s\n\n" + CommonPressEnter

	// Vault
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

	// Auth
	AuthLoginTitle     = "Вход в GophKeeper"
	AuthRegisterTitle  = "Регистрация"
	AuthUsernameLabel  = "Логин: %s"
	AuthPasswordLabel  = "Пароль: %s"
	AuthLoginButton    = "Вход"
	AuthRegisterButton = "Регистрация"
	AuthTabHint        = "Нажмите Enter для подтверждения, Tab для переключения"

	// Add
	AddSelectPrompt   = "Выберите тип хранимой информации:\n\n"
	AddChoiceTemplate = "%s %s\n"
)
