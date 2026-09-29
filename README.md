<div align="center">

# Payment System

**Платёжная система на Go: интерфейсы, полиморфизм и точные деньги**

[![Go](https://img.shields.io/badge/go-1.26%2B-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Decimal](https://img.shields.io/badge/money-shopspring%2Fdecimal-7D5CFF?style=flat-square)](https://github.com/shopspring/decimal)
[![Methods](https://img.shields.io/badge/payment%20methods-3-3FB950?style=flat-square)](#методы-оплаты)

<p align="center">
  <em>CreditCard · PayPal · Crypto — за интерфейсом <code>PaymentMethod</code></em>
</p>

</div>

---

## Что это

Учебная реализация платёжной системы: единый интерфейс `PaymentMethod`, три способа оплаты и процессор, который считает комиссию, ведёт историю и скрывает различия между методами.

Ключевая идея — **полиморфизм**: процессору неважно, карта это, PayPal или крипта. Он получает `PaymentMethod` и работает только с контрактом интерфейса. Добавить новый метод оплаты — значит написать один файл с пятью методами, ничего не меняя в процессоре.

Деньги считаются через [`shopspring/decimal`](https://github.com/shopspring/decimal) — никаких `float64` в денежных величинах.

## Архитектура

```
                    ┌──────────────────────────┐
                    │  PaymentProcessor         │
                    │  ProcessPayment(method,…) │
                    │  GetHistory()             │
                    └────────────┬─────────────┘
                                 │ принимает PaymentMethod
                                 ▼
        ┌────────────────────────┼────────────────────────┐
        │                        │                        │
   ┌────▼─────┐            ┌─────▼─────┐            ┌─────▼────┐
   │CreditCard│            │  PayPal   │            │  Crypto  │
   │  2%      │            │ 3% + $0.50│            │   50%    │
   └──────────┘            └───────────┘            └──────────┘
              все три реализуют PaymentMethod
```

## Запуск

```bash
go mod download
go run ./cmd/app
go vet ./...
```

Требуется Go **1.26+**. Единственная зависимость — `shopspring/decimal`.

## Пример вывода

```
Оплата через Credit Card
Processing credit card payment of $100

Оплата через PayPal
Processing PayPal payment of $50

Попытка оплаты через Credit Card с невалидными данными
Ошибка: Invalid credit card details

=== История платежей ===
Платеж 1: метод=Credit Card, сумма=100, комиссия=2, итого=102, успех=true
Платеж 2: метод=PayPal, сумма=50, комиссия=2, итого=52, успех=true
Платеж 3: метод=Credit Card, сумма=10, комиссия=0.2, итого=10.2, успех=false
```

## Использование

```go
import paymentsystem "PaymentSystemModule"

processor := paymentsystem.NewPaymentProcessor()

// Любой метод — один и тот же вызов
var method paymentsystem.PaymentMethod = paymentsystem.NewCreditCard(
    "1234567890123456",
    "12/12/2027",
)

err := processor.ProcessPayment(method, decimal.NewFromFloat(100.00))
if err != nil {
    log.Fatal(err)
}

for _, record := range processor.GetHistory() {
    fmt.Printf("%s: %s на %s (комиссия %s, итого %s)\n",
        record.ID, record.Method, record.Amount, record.Fee, record.Total)
}
```

## Методы оплаты

| Метод | Конструктор | Комиссия | Валидация |
| --- | --- | --- | --- |
| **CreditCard** | `NewCreditCard(number, expiryDate)` | 2% | 16 цифр + дата `ДД/ММ/ГГГГ` + год в будущем |
| **PayPal** | `NewPayPal(email)` | 3% + $0.50 | email по RFC-подобному regex |
| **Crypto** | `NewCrypto(walletAddress)` | 50% | префикс `0x` + ровно 42 символа |

### Формат даты карты

```
ДД/ММ/ГГГГ   ← строго

31/12/2027   ✓ валидно
12/2027      ✗ не хватает дня
31/13/2027   ✗ месяц вне 1–12
01/01/2020   ✗ год в прошлом
```

## API

### `PaymentMethod`

```go
type PaymentMethod interface {
    Pay(amount decimal.Decimal) error
    GetFree(amount decimal.Decimal) decimal.Decimal
    Validate() bool
    GetName() string
}
```

| Метод | Назначение |
| --- | --- |
| `Pay` | Выполняет оплату, возвращает ошибку при невалидных данных |
| `GetFree` | Рассчитывает комиссию для суммы |
| `Validate` | Проверяет реквизиты **без** выполнения платежа |
| `GetName` | Человекочитаемое имя для истории |

### `PaymentProcessor`

```go
type PaymentProcessor interface {
    ProcessPayment(method PaymentMethod, amount decimal.Decimal) error
    GetHistory() []PaymentRecord
}
```

### `PaymentRecord`

| Поле | Тип | Содержимое |
| --- | --- | --- |
| `ID` | `string` | порядковый номер платежа |
| `Method` | `string` | имя метода оплаты |
| `Amount` | `decimal.Decimal` | исходная сумма |
| `Fee` | `decimal.Decimal` | комиссия |
| `Total` | `decimal.Decimal` | сумма с комиссией |
| `Success` | `bool` | результат валидации |
| `Timestamp` | `time.Time` | время обработки |

## Файлы

| Файл | Содержимое |
| --- | --- |
| `payment-interfaces.go` | Контракты: `PaymentMethod`, `PaymentProcessor`, `PaymentRecord` |
| `payment-processor.go` | `DefaultPaymentProcessor` — расчёт комиссий и история |
| `CreditCard.go` | Банковская карта: проверка длины номера и срока действия |
| `PayPal.go` | PayPal с валидацией email и фиксированной комиссией |
| `cmd/app/main.go` | Демо: успешная оплата, PayPal, отказ по невалидной карте, вывод истории |

## Ключевые решения

**`decimal.Decimal` вместо `float64`.** `0.1 + 0.2 != 0.3` в binary floating point. Для денег это недопустимо: `decimal` хранит число как целое + экспонента и даёт точную арифметику. Ради этого в проекте одна внешняя зависимость.

**Проверки интерфейсов на этапе компиляции.** В `payment-interfaces.go` объявлено `var _ PaymentMethod = (*CreditCard)(nil)`. Если сигнатура метода разойдётся с интерфейсом, сборка упадёт сразу — без этого Go не проверяет, что тип реализует интерфейс, пока его не подставят в переменную этого типа.

**`Pay` и `Validate` разделены.** Метод оплаты может проверить реквизиты двумя способами: `Validate()` — чистая проверка без побочных эффектов, `Pay()` — проверка с выполнением транзакции. Это позволяет валидировать до оплаты.

**`decimal.NewFromFloat(0.02)` для комиссий.** Ставка задана константой, но результат сразу приводится к `decimal` — дальше вся арифметика точная.

## Гарантии и семантика

| Свойство | Поведение |
| --- | --- |
| Порядок платежей | сохраняется, `ID` растёт на 1 |
| Точность сумм | точная, без ошибок округления |
| Невалидные данные | `Pay` возвращает ошибку, платёж попадает в историю с `Success: false` |
| Комиссия PayPal | 3% **плюс** фиксированные $0.50 |
| Валюта | подразумевается USD, в коде не зафиксирована |

## Известные ограничения

- **Процессор не потокобезопасен.** `history` — обычный слайс без мьютекса; параллельные вызовы `ProcessPayment` приведут к гонке данных.
- **`GetHistory` отдаёт внутренний слайс.** Вызывающий код может изменить историю, повлияв на состояние процессора. Нужна копия.
- **`Success` и `err` вычисляются по-разному.** `err` берётся из `Pay`, а `Success` — из отдельного вызова `Validate`. Логика работает, но хрупка: она зависит от того, что `Pay` внутри себя тоже зовёт `Validate`.
- **Комиссия считается до проверки.** `GetFree` вызывается до `Pay`, поэтому невалидный платёж всё равно получает ненулевую комиссию в истории.
- **`ID` — просто счётчик.** Не уникален между перезапусками, не пригоден как бизнес-идентификатор.
- **Поля структур неэкспортиемые и без геттеров.** После `Validate` нельзя прочитать распарсенные `day` / `month` / `year`.
- **Нет тестов.** Проверка только через вывод демо.
- **Нет валюты, возвратов и чека.** Нет `Currency` в `PaymentRecord`, нет операции отмены платежа.

## Возможные улучшения

- `sync.Mutex` в процессоре + возврат копии истории.
- Реальная проверка номера карты по алгоритму Луна (сейчас только длина 16).
- Конфигурируемые ставки комиссий вместо констант в коде.
- Добавление валюты в `PaymentRecord` и конвертации.
- Тесты с `testify` и табличные кейсы на валидацию каждого метода.
- Отдельный пакет `refund` и журнал транзакций.

## Назначение

Учебный проект по **ООП в Go**: интерфейсы, полиморфизм, внедрение зависимостей через конструктор (`NewPaymentProcessor() PaymentProcessor` возвращает интерфейс, а не структуру). Показывает, как добавить новый способ оплаты, не меняя существующий код, и почему деньги в Go считают через `decimal`.

## Лицензия

Свободное использование.

---

<div align="center">
  <sub>Go · interfaces · shopspring/decimal</sub>
</div>
