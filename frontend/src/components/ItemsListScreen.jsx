import React, { useEffect, useState } from 'react';
import axios from 'axios';

function ItemsListScreen({ userToken, username, setUserToken }) {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [checkoutLoading, setCheckoutLoading] = useState(false);

  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    setLoading(true);
    try {
      const res = await axios.get('/api/items');
      setItems(res.data);
    } catch (err) {
      window.alert('Failed to fetch items');
    } finally {
      setLoading(false);
    }
  };

  const addItemToCart = async (itemId) => {
    try {
      await axios.post(
        '/api/carts',
        { itemId },
        { headers: { Authorization: `Bearer ${userToken}` } }
      );
      window.alert('Item added to cart');
    } catch (err) {
      window.alert('Failed to add item to cart');
    }
  };

  const fetchCartItems = async () => {
    try {
      const res = await axios.get('/api/carts', {
        headers: { Authorization: `Bearer ${userToken}` },
      });
      if (res.data.items && res.data.items.length > 0) {
        const msg = res.data.items.map((ci) => `cart_id: ${res.data.cartId}, item_id: ${ci.itemId}`).join('\n');
        window.alert(msg);
      } else {
        window.alert('Cart is empty');
      }
    } catch (err) {
      window.alert('Failed to fetch cart items');
    }
  };

  const fetchOrderHistory = async () => {
    try {
      const res = await axios.get('/api/orders', {
        headers: { Authorization: `Bearer ${userToken}` },
      });
      if (res.data.length > 0) {
        const msg = res.data.map((o) => `Order id: ${o.id}`).join('\n');
        window.alert(msg);
      } else {
        window.alert('No orders found');
      }
    } catch (err) {
      window.alert('Failed to fetch orders');
    }
  };

  const handleCheckout = async () => {
    setCheckoutLoading(true);
    try {
      // Get cartId first
      const res = await axios.get('/api/carts', {
        headers: { Authorization: `Bearer ${userToken}` },
      });
      if (!res.data.cartId) {
        window.alert('No cart found');
        setCheckoutLoading(false);
        return;
      }
      await axios.post(
        '/api/orders',
        { cartId: res.data.cartId },
        { headers: { Authorization: `Bearer ${userToken}` } }
      );
      window.alert('Order successful');
      fetchItems(); // reload items (if needed)
    } catch (err) {
      window.alert('Failed to checkout');
    } finally {
      setCheckoutLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('userToken');
    localStorage.removeItem('username');
    setUserToken('');
  };

  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)'
    }}>
      {/* Header */}
      <div style={{
        background: 'white',
        boxShadow: '0 2px 10px rgba(0,0,0,0.1)',
        padding: '20px 0',
        marginBottom: '30px'
      }}>
        <div style={{
          maxWidth: '1200px',
          margin: '0 auto',
          padding: '0 20px',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center'
        }}>
          <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '15px'
          }}>
            <div style={{
              width: '40px',
              height: '40px',
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              borderRadius: '50%',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: '20px',
              color: 'white'
            }}>🛒</div>
            <div>
              <h1 style={{
                margin: 0,
                fontSize: '24px',
                color: '#333',
                fontWeight: '700'
              }}>Shopping Store</h1>
              <p style={{
                margin: 0,
                color: '#666',
                fontSize: '14px'
              }}>Welcome back, {username}! 👋</p>
            </div>
          </div>
          <div style={{
            display: 'flex',
            gap: '10px',
            alignItems: 'center'
          }}>
            <button
              onClick={fetchCartItems}
              style={{
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                color: 'white',
                border: 'none',
                padding: '10px 20px',
                borderRadius: '25px',
                cursor: 'pointer',
                fontSize: '14px',
                fontWeight: '600',
                transition: 'transform 0.2s'
              }}
              onMouseOver={(e) => e.target.style.transform = 'translateY(-2px)'}
              onMouseOut={(e) => e.target.style.transform = 'translateY(0)'}
            >
              🛒 View Cart
            </button>
            <button
              onClick={handleCheckout}
              disabled={checkoutLoading}
              style={{
                background: checkoutLoading ? '#ccc' : 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
                color: 'white',
                border: 'none',
                padding: '10px 20px',
                borderRadius: '25px',
                cursor: checkoutLoading ? 'not-allowed' : 'pointer',
                fontSize: '14px',
                fontWeight: '600',
                transition: 'transform 0.2s'
              }}
              onMouseOver={(e) => !checkoutLoading && (e.target.style.transform = 'translateY(-2px)')}
              onMouseOut={(e) => !checkoutLoading && (e.target.style.transform = 'translateY(0)')}
            >
              {checkoutLoading ? '⏳ Processing...' : '💳 Checkout'}
            </button>
            <button
              onClick={fetchOrderHistory}
              style={{
                background: 'transparent',
                color: '#667eea',
                border: '2px solid #667eea',
                padding: '8px 16px',
                borderRadius: '25px',
                cursor: 'pointer',
                fontSize: '14px',
                fontWeight: '600',
                transition: 'all 0.3s'
              }}
              onMouseOver={(e) => {
                e.target.style.background = '#667eea';
                e.target.style.color = 'white';
              }}
              onMouseOut={(e) => {
                e.target.style.background = 'transparent';
                e.target.style.color = '#667eea';
              }}
            >
              📋 Orders
            </button>
            <button
              onClick={handleLogout}
              style={{
                background: 'transparent',
                color: '#e74c3c',
                border: '2px solid #e74c3c',
                padding: '8px 16px',
                borderRadius: '25px',
                cursor: 'pointer',
                fontSize: '14px',
                fontWeight: '600',
                transition: 'all 0.3s'
              }}
              onMouseOver={(e) => {
                e.target.style.background = '#e74c3c';
                e.target.style.color = 'white';
              }}
              onMouseOut={(e) => {
                e.target.style.background = 'transparent';
                e.target.style.color = '#e74c3c';
              }}
            >
              🚪 Logout
            </button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto',
        padding: '0 20px'
      }}>
        <div style={{
          textAlign: 'center',
          marginBottom: '40px'
        }}>
          <h2 style={{
            fontSize: '32px',
            color: '#333',
            marginBottom: '10px',
            fontWeight: '700'
          }}>🛍️ Our Products</h2>
          <p style={{
            color: '#666',
            fontSize: '18px',
            margin: 0
          }}>Discover amazing products just for you!</p>
        </div>

        {loading ? (
          <div style={{
            textAlign: 'center',
            padding: '60px 20px',
            background: 'white',
            borderRadius: '16px',
            boxShadow: '0 10px 30px rgba(0,0,0,0.1)'
          }}>
            <div style={{
              fontSize: '48px',
              marginBottom: '20px'
            }}>⏳</div>
            <h3 style={{
              color: '#333',
              marginBottom: '10px'
            }}>Loading Products...</h3>
            <p style={{
              color: '#666',
              margin: 0
            }}>Please wait while we fetch the latest items</p>
          </div>
        ) : (
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
            gap: '25px',
            marginBottom: '40px'
          }}>
            {items.map((item) => (
              <div
                key={item.id}
                style={{
                  background: 'white',
                  borderRadius: '16px',
                  padding: '25px',
                  boxShadow: '0 10px 30px rgba(0,0,0,0.1)',
                  transition: 'transform 0.3s, box-shadow 0.3s',
                  cursor: 'pointer'
                }}
                onMouseOver={(e) => {
                  e.currentTarget.style.transform = 'translateY(-5px)';
                  e.currentTarget.style.boxShadow = '0 20px 40px rgba(0,0,0,0.15)';
                }}
                onMouseOut={(e) => {
                  e.currentTarget.style.transform = 'translateY(0)';
                  e.currentTarget.style.boxShadow = '0 10px 30px rgba(0,0,0,0.1)';
                }}
              >
                <div style={{
                  width: '60px',
                  height: '60px',
                  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  borderRadius: '50%',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: '24px',
                  color: 'white',
                  marginBottom: '20px'
                }}>📦</div>
                <h3 style={{
                  color: '#333',
                  fontSize: '20px',
                  fontWeight: '600',
                  marginBottom: '10px',
                  margin: '0 0 10px 0'
                }}>{item.name}</h3>
                <p style={{
                  color: '#666',
                  fontSize: '14px',
                  lineHeight: '1.5',
                  marginBottom: '20px',
                  margin: '0 0 20px 0'
                }}>{item.description || 'No description available'}</p>
                <div style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center'
                }}>
                  <div style={{
                    fontSize: '24px',
                    fontWeight: '700',
                    color: '#11998e'
                  }}>₹{item.price}</div>
                  <button
                    onClick={() => addItemToCart(item.id)}
                    style={{
                      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                      color: 'white',
                      border: 'none',
                      padding: '12px 24px',
                      borderRadius: '25px',
                      cursor: 'pointer',
                      fontSize: '14px',
                      fontWeight: '600',
                      transition: 'transform 0.2s'
                    }}
                    onMouseOver={(e) => e.target.style.transform = 'translateY(-2px)'}
                    onMouseOut={(e) => e.target.style.transform = 'translateY(0)'}
                  >
                    🛒 Add to Cart
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export default ItemsListScreen;
