# App Flow

- Order

  - event_handlers: VendorConfirmOrderItemForDelivery
  - events_owned: Order Created, OrderCanceled,
  - grpc handlers: GetOrderById, GetOrders | filter - userId,
  - repo: createOrder, getOrders, updateOrderStatus, getOrderById,
  - http_handlers: repo + cancelOrderStatus
  - tables/models: order - id, common, paidAt, status, isPaid, isCanceled; order_items: product_id, common, store_id, quantity
