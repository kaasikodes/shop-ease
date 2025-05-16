# App Flow

- Conten Service
- Product

  - http handlers: most of the repo methods, Get Products (will show products from different stores), and combine the discounts for price
  - event_handlers: will react to inventory update to update its own product inventory count,
  - events_owned: Product Created
  - grpc handlers: Create Discount, Get Discounts(to be used to show the price of products), when displayng products - the product service is responsible
  - repo: createCategory, createProduct, updateCategory, updateProduct, getProducts(with meta data inventory: talk to inventory service), bulkAddProducts, create Discount, set Price for product,
  - other category, and sub category for effiecient querying
  - \*Need for arrival date, will prove to be helpful for e-commerce

- Order
- Search & Recommend
- Review & Rating
