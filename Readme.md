Проект написано для личных нужд.

Основная функциональность - получение всех комментариев по виджету из сервиса cackle.me. Написано на Go.

Описание API сервиса можно всегда посмотреть на странице [API Cackle.me для разработчиков](https://cackle.me/help/widget-api).

Для управления процессом импорта используются переменные окружения. Создать их можно через файл .env или любым другим способом.

TIMEOUT=число, время в секундах между запросами к API. Рекомендация сервиса - 5 секунд.

ID=число, номер виджета, который можно увидеть в админке сервиса.

SITE_API_KEY=строковый токен, для доступа к API, смотреть в админке сервиса.

ACCOUNT_API_KEY=строковый токен, для доступа к API, смотреть в админке сервиса.

>Чтобы получить ID и токены, надо выбрать виджет в админке сервиса Cackle.me, перейти в раздел "Установить", затем "CMS Платформа", выбрать Wordpress, там покажутся необходимые данные. Выделенные жирным значения и нужно вставить после =.