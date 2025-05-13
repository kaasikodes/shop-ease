package products

type ProductRepo interface {
	Save(id int, payload Product)
	GetById(id int) Product
}

// in memory repo
type InMemoryProductRepo struct {
	products map[int]Product
}

func NewInMemoryProductRepo() ProductRepo {
	return &InMemoryProductRepo{
		products: make(map[int]Product),
	}

}

func (p *InMemoryProductRepo) Save(id int, payload Product) {
	p.products[id] = payload

}
func (p *InMemoryProductRepo) GetById(id int) Product {
	return p.products[id]

}
