package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/go-redis/redis"
)


type Repository struct {
	Db    *gorm.DB
	Redis *redis.Client
}


func (r *Repository) CreateTables(db *gorm.DB) error {
	
	var created bool
	created = r.SetUpUser(db)
	created = r.SetUpCart(db)
	created = r.SetUpItem(db)
	created = r.SetUpOrder(db)
	created = r.SetUpCartItemRelation(db)
	
	if created == false {
	// add foreign keys
		if err := db.Model(&User{}).AddForeignKey("cart_id", "carts(id)", "CASCADE", "CASCADE").Error; err != nil {
			return err
		}
		if err := db.Model(&Order{}).AddForeignKey("cart_id", "carts(id)", "CASCADE", "CASCADE").Error; err != nil {
			return	err
		}
		if err := db.Model(&Order{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error; err != nil {
			return	err
		}
		if err := db.Model(&CartItemRelation{}).AddForeignKey("cart_id", "carts(id)", "CASCADE", "CASCADE").Error; err != nil {
			return	err
		}
		if err := db.Model(&CartItemRelation{}).AddForeignKey("item_id", "items(id)", "CASCADE", "CASCADE").Error; err != nil {
			return	err
		}
	}
	return nil
}

func (r *Repository) SetUpUser(db *gorm.DB) bool {
	if db.HasTable(&User{}) == false {
		db.CreateTable(&User{})
		return false
	}
	return true
}

func (r *Repository) SetUpCart(db *gorm.DB) bool {
	if db.HasTable(&Cart{}) == false {
		db.CreateTable(&Cart{})
		return false
	}
	return true
}

func (r *Repository) SetUpItem(db *gorm.DB) bool {
	if db.HasTable(&Item{}) == false {
		db.CreateTable(&Item{})
		return false
	}	
	return true
}

func (r *Repository) SetUpOrder(db *gorm.DB) bool {
	if db.HasTable(&Order{}) == false {
		db.CreateTable(&Order{})
		return false
	}	
	return true
}


func (r *Repository) SetUpCartItem(db *gorm.DB) bool {
	if db.HasTable(&CartItemRelation{}) == false {
		db.CreateTable(&CartItemRelation{})
		return false
	} 
	return true
}


func (r *Repository) SetUpCartItemRelation(db *gorm.DB) bool {
	if db.HasTable(&CartItemRelation{}) == false {
		db.CreateTable(&CartItemRelation{})
		return false
	}	
	return true
}


//Ensures tables exist
func (r *Repository) Init() error {
	if err := r.Db.AutoMigrate(
		new(User),
		new(Cart),
		new(Item),
		new(Order),
		new(CartItemRelation),
	).Error; err != nil {
		return err
	}

	return nil
}


//Authenticate enterprise with provided credentials
func (r *Repository)  AuthByPassword(username, password string) (*User, error) {
	//query Enterprise by credentials
	user := new(User)
	if result := r.Db.Take(user, &User{UserName: username, Password: password}); result.Error != nil {
		return nil, result.Error
	} else {
		return user, nil
	}
}


//get user by the name
func (r *Repository) GetUserByName(userName string) (*User, error) {
        user := new(User)
        if result := r.Db.Take(user, &User{UserName: userName}); result.Error != nil {
                return nil, result.Error
        }
        return user, nil
}


//get user by the id
func (r *Repository) GetUserById(id int) (*User, error) {
        user := new(User)
        if result := r.Db.Take(user, &User{Id: uint64(id)}); result.Error != nil {
                return nil, result.Error
        }
        return user, nil
}


//get user by toke
func (r *Repository) GetUserByToken(token string) (*User, error) {
        user := new(User)
        if result := r.Db.Take(user, &User{Token: token}); result.Error != nil {
                return nil, result.Error
        }
        return user, nil
}

func (r *Repository) CreateUser(name, userName, password, token string, cartId int) error {
	user := &User{
		Name : name,
		UserName : userName,
		Password :  password,
		Token	: token,
		CartId	: uint64(cartId),
	}
	return r.Db.Create(user).Error
}

// list all users
func (r *Repository) ListAllUsers() (*[]User, error) {
	var users []User
	result := r.Db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}


// creat item
func (r *Repository) CreateItem(itemName string ) error {
	item := &Item{
		Name : itemName,
	}
	return r.Db.Create(item).Error
}

//get item by name
func (r *Repository) GetItemByName(itemName string) (*Item, error) {
        item := new(Item)
        if result := r.Db.Take(item, &Item{Name: itemName}); result.Error != nil {
                return nil, result.Error
        }
        return item, nil
}

//list all item
func (r *Repository) ListAllItems() (*[]Item, error) {
	var items []Item
	result := r.Db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}


//get cart by the id
func (r *Repository) GetCartById(id int) (*Cart, error) {
    cart := new(Cart)
    if result := r.Db.Take(cart, &Cart{Id: uint64(id)}); result.Error != nil {
            return nil, result.Error
    }
    return cart, nil
}

//add item to cart (item cart relation)
func (r *Repository) AddItemToCart(cartId, itemId int)  error {
	cartItem := &CartItemRelation{
		CartId : uint64(cartId),
		ItemId : uint64(itemId),
	}
	return r.Db.Create(cartItem).Error
}

//mark item purchased in 
func (r *Repository) UpdateCartItemRelation(cartId int,isPurchased bool)  error {
	result := r.Db.Model(&CartItemRelation{}).Where("cart_id = ?", cartId).Update("is_purchased", isPurchased)
    if result.Error != nil {
    	return result.Error
    }
    return nil
}

//create cart
func (r *Repository) CreateCart() (int, error) {
	cart := &Cart{
		IsPurchased : true,
	}
	err := r.Db.Create(cart).Error
	if err != nil {
		return 0,err
	} else {
		return int(cart.Id), nil
	}
}

//update cart
func(r *Repository) UpdateCart(isPurchased bool,cartId int) error {
	result := r.Db.Model(&Cart{}).Where("	id = ?", uint64(cartId)).Update("is_purchased", isPurchased)
	return result.Error
}

//list all carts
func (r *Repository) ListAllCarts() (*[]Cart, error) {
	var carts []Cart
	result := r.Db.Find(&carts)
	if result.Error != nil {
		return nil, result.Error
	}
	return &carts, nil
}


//create order
func (r *Repository) CreateOrder(cartId, userId int)  error {
	oder := &Order{
		CartId : uint64(cartId),
		UserId : uint64(userId),
	}
	return r.Db.Create(oder).Error
}

//list all orders
func (r *Repository) ListAllOrders() (*[]Order, error) {
	var orders []Order
	result := r.Db.Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return &orders, nil
}


func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		Db:    db,
		Redis: redis,
	}
}