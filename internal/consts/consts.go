package consts

// DeliveredST - статуc заказа, что он доставлен курьером в ПВЗ
const DeliveredST = "delivered"

// AcceptedSt - статус заказа, что клиент забрал его из ПЗВ
const AcceptedSt = "accepted"

// ReturnedSt - статус заказа, что клиент его вернул
const ReturnedSt = "returned"

// IssueAction - используется, чтобы клиент взял заказы из пвз
const IssueAction = "issue"

// ReturnAction - используется, чтобы клиент вернул заказы в пвз
const ReturnAction = "return"

// CreatedTask задача создана
const CreatedTask = "CREATED"

// ProcessingTask задача в процессе выполнения
const ProcessingTask = "PROCESSING"

// FailedTask задача не обработана
const FailedTask = "FAILED"

// CompletedTask задача выполнена
const CompletedTask = "COMPLETED"

// NoAttemptsTask попытки отправления закончились
const NoAttemptsTask = "NO_ATTEMPTS_LEFT"
