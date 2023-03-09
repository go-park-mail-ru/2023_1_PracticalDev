# 2023_1_PracticalDev [![Docs](https://godoc.org/github.com/go-park-mail-ru/2023_1_PracticalDev?status.svg)](http://pkg.go.dev/github.com/go-park-mail-ru/2023_1_PracticalDev)

Проект Pinterest команды "Practical Dev"

[Ссылка на репозиторий по фронтенду](https://github.com/frontend-park-mail-ru/2023_1_PracticalDev)

Ссылка на деплой - [сайт](http://pickpin.ru), [api](http://pickpin.ru/api), [OpenAPI](http://pickpin.ru/api/docs)

## Авторы 

- [Гречко Георгий](https://github.com/geogreck)

- [Кирилл Киселев](https://github.com/t1d333)

- [Евгений Димитриев](https://github.com/UjinIaly)

- [Вячеслав Шагалов](https://github.com/SlavaShagalov)

## Менторы

- [Дмитрий Ильин](https://github.com/Neytrinoo) - Backend

- [Фарис Набиев](https://github.com/rflban) - Frontend


## Руководство по использованию:

| | |
|---------------------|--------------------------------------------|
| make build          | Собирает весь `backend`                    |
| make run            | Запускает `backend` на порту 8080          |
| make unit-test      | Запускает юнит тесты для всего `backend`'a |
| make print-coverage | Выыводит отчет о покрытии                  |
| make coverage       | Генерирует coverage.out                    |
| make lint           | Запускает проверку линтером                |
| make db             | Поднимает базу данных                      |
| make deploy         | Поднимает фронт, бек, доку и бд            |

Фронтенд доступен на [localhost](http://localhost) или [localhost:8000](http://localhost:8000)

Бэкенд доступен на [localhost/api](http://localhost/api) или [localhost:8080](http://localhost:8080)

OpenAPI доступен на [localhost/api/docs](http://localhost/api/docs)
