CREATE TABLE IF NOT EXISTS sharing_formulas (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  app INT NOT NULL,         -- percentage to the app
  vendor INT NOT NULL,      -- percentage to the vendor
  basedOn ENUM('sale', 'profit') DEFAULT 'sale',
  description TEXT,
  createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS orders (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  quantity INT NOT NULL,
  unitPrice DECIMAL(12,2) NOT NULL,
  status ENUM('pending', 'processing', 'fulfilled', 'cancelled') NOT NULL DEFAULT 'pending',
  productId BIGINT NOT NULL,
  storeId BIGINT NOT NULL,
  fulfillingInventoryId BIGINT, //TODO: should have a table relationship
  sharingFormulaId BIGINT,
  createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (productId) REFERENCES products(id) ON DELETE CASCADE,
  FOREIGN KEY (storeId) REFERENCES stores(id) ON DELETE CASCADE,
  FOREIGN KEY (fulfillingInventoryId) REFERENCES inventories(id) ON DELETE SET NULL,
  FOREIGN KEY (sharingFormulaId) REFERENCES sharing_formulas(id) ON DELETE SET NULL
);
