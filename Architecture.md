https://www.omg.org/spec/UML/#issues

## Архитектура(устаревшая)

::: mermaid
classDiagram
    class Repository
    Repository : +Orders map[int64]models.OrderInfo
    Repository : +Users  map[int64]map[int64]struct{}
    Repository : -path   string
    Repository : +AddOrder(orderID, userID, weight, price int64, shelflife time.Time, packType string, extraPack bool) error
    Repository : +RemoveOrder(orderID int64) error
    Repository : +GetOrderInfo(orderID int64) (models.OrderInfo, error)
    Repository : +GetOrdersForUser(userID int64) ([]models.Order, error)
    Repository : +GetAllOrders() ([]models.Order, error)
    Repository : +UpdateStatus(orderID int64, newStatus string) error
    Repository : +ReadFile(path string) ([]models.Order, error)
	Repository : -saveToFile() error

    class Repository_Interface
    Repository_Interface : + AddOrder(orderID, userID, weight, price int64, shelflife time.Time, packType string, extraPack bool) error
    Repository_Interface : + RemoveOrder(orderID int64) error
    Repository_Interface : + GetOrderInfo(orderID int64) (models.OrderInfo, error)
    Repository_Interface : + GetOrdersForUser(userID int64) ([]models.Order, error)
    Repository_Interface : + GetAllOrders() ([]models.Order, error)
    Repository_Interface : + UpdateStatus(orderID int64, newStatus string) error
    Repository_Interface : + ReadFile(path string) ([]models.Order, error)

    class Accept_Service
    Accept_Service
    Accept_Service : - storageRepository StorageRepository
    Accept_Service : + AcceptOrder(orderID, userID, weight, price int64, packType string, shelflife time.Time, extraPack bool) error
    Accept_Service : + BulkAcceptOrders(path string) (int, error)
    Accept_Service : + ProcessClientIssue(userID int64, orderIDs []int64) (int64, error)

    class Returns_Service
    Returns_Service : - storageRepository StorageRepository
    Returns_Service : + ProcessClientReturn(userID int64, orderIDs []int64) (int64, error)
    Returns_Service : + ReturnOrder(id int64) error


    class Show_Service
    Show_Service : - storageRepository StorageRepository
    Show_Service : + ListOrdersUserAllN(userID int64, opt int64) ([]models.Order, error)
    Show_Service : + ListOrdersUserPVZ(userID int64) ([]models.Order, error)
    Show_Service : + ListReturns(limit, page int64) ([]models.Order, int, error)
    Show_Service : + ListHistory() ([]models.Order, error)


    class Accept_Service_Interface
    Accept_Service_Interface : + AcceptOrder(orderID, userID, weight, price int64, packType string, shelflife time.Time, extraPack bool) error
    Accept_Service_Interface : + BulkAcceptOrders(path string) (int, error)
    Accept_Service_Interface : + ProcessClientIssue(userID int64, orderIDs []int64) (int64, error)

    class Returns_Service_Interface
    Returns_Service_Interface : + ReturnOrder(id int64) error
    Returns_Service_Interface : + ProcessClientReturn(userID int64, orderIDs []int64) (int64, error)

    class Show_Service_Interface
    Show_Service_Interface : + ListOrdersUserAllN(userID int64, opt int64) ([]models.Order, error)
    Show_Service_Interface : + ListOrdersUserPVZ(userID int64) ([]models.Order, error)
    Show_Service_Interface : + ListReturns(limit, page int64) ([]models.Order, int, error)
    Show_Service_Interface : + ListHistory() ([]models.Order, error)


    class Handler
    Handler : - acceptServ  AcceptService
	Handler : - returnsServ ReturnsService
	Handler : -	showServ    ShowService

    Handler : + AcceptOrder(args []string)
    Handler : + BulkAcceptOrders(args []string)
    Handler : + Help()
    Handler : + ListHistory()
    Handler : + ListOrdersUser(args []string)
    Handler : + ListReturns(args []string)
    Handler : + ProcessClient(args []string)
    Handler : + ReturnOrderToCourier(args []string)
    Handler : + Run() error


    class Order
    Order : + OrderID       int64 
    Order : + UserID        int64 
    Order : + Weight        int64 
    Order : + Price         int64 `
    Order : + Packaging     Pack     
    Order : + InStorageFrom time.Time
    Order : + ShelfLife     time.Time
    Order : + Status        string   
    Order : + TwoDaysOfLife time.Time
    Order : + LastUpdate    time.Time

    class Pack
    Pack : + PackType  string
    Pack : + ExtraPack string

    %% Interface
    Repository_Interface <|-- Repository : implement
    Accept_Service_Interface <|-- Accept_Service : implement
    Returns_Service_Interface <|-- Returns_Service : implement
    Show_Service_Interface <|-- Show_Service : implement

    %% Repository
    Order ..> Pack : uses
    Repository ..> Order : uses
    
    %% Services
    Accept_Service ..> Repository_Interface : uses
    Returns_Service ..> Repository_Interface : uses
    Show_Service ..> Repository_Interface : uses

    %% Handler
    Handler ..> Accept_Service_Interface : uses
    Handler ..> Returns_Service_Interface : uses
    Handler ..> Show_Service_Interface : uses
:::
