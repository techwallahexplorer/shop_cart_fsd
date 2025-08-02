import React, { useState } from 'react';
import axios from 'axios';

function LoginScreen({ setUserToken, setUsername, setShowLogin }) {
  const [usernameInput, setUsernameInput] = useState('');
  const [passwordInput, setPasswordInput] = useState('');
  const [loading, setLoading] = useState(false);

  const handleUserLogin = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await axios.post('/api/users/login', {
        username: usernameInput,
        password: passwordInput,
      });
      setUserToken(res.data.token);
      setUsername(usernameInput);
      localStorage.setItem('userToken', res.data.token);
      localStorage.setItem('username', usernameInput);
    } catch (err) {
      window.alert('Invalid username/password');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '20px'
    }}>
      <div style={{
        maxWidth: 400,
        width: '100%',
        background: 'white',
        padding: '40px',
        borderRadius: '16px',
        boxShadow: '0 20px 40px rgba(0,0,0,0.1)',
        textAlign: 'center'
      }}>
        <div style={{
          width: '60px',
          height: '60px',
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          borderRadius: '50%',
          margin: '0 auto 20px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '24px',
          color: 'white'
        }}>🛒</div>
        <h2 style={{
          color: '#333',
          marginBottom: '30px',
          fontSize: '28px',
          fontWeight: '600'
        }}>Welcome Back!</h2>
        <form onSubmit={handleUserLogin}>
          <div style={{ marginBottom: '20px' }}>
            <input
              type="text"
              placeholder="Username"
              value={usernameInput}
              onChange={(e) => setUsernameInput(e.target.value)}
              required
              style={{
                width: '100%',
                padding: '15px',
                border: '2px solid #f0f0f0',
                borderRadius: '10px',
                fontSize: '16px',
                outline: 'none',
                transition: 'border-color 0.3s',
                boxSizing: 'border-box'
              }}
              onFocus={(e) => e.target.style.borderColor = '#667eea'}
              onBlur={(e) => e.target.style.borderColor = '#f0f0f0'}
            />
          </div>
          <div style={{ marginBottom: '25px' }}>
            <input
              type="password"
              placeholder="Password"
              value={passwordInput}
              onChange={(e) => setPasswordInput(e.target.value)}
              required
              style={{
                width: '100%',
                padding: '15px',
                border: '2px solid #f0f0f0',
                borderRadius: '10px',
                fontSize: '16px',
                outline: 'none',
                transition: 'border-color 0.3s',
                boxSizing: 'border-box'
              }}
              onFocus={(e) => e.target.style.borderColor = '#667eea'}
              onBlur={(e) => e.target.style.borderColor = '#f0f0f0'}
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            style={{
              width: '100%',
              padding: '15px',
              background: loading ? '#ccc' : 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              color: 'white',
              border: 'none',
              borderRadius: '10px',
              fontSize: '16px',
              fontWeight: '600',
              cursor: loading ? 'not-allowed' : 'pointer',
              transition: 'transform 0.2s',
              marginBottom: '20px'
            }}
            onMouseOver={(e) => !loading && (e.target.style.transform = 'translateY(-2px)')}
            onMouseOut={(e) => !loading && (e.target.style.transform = 'translateY(0)')}
          >
            {loading ? '🔄 Logging in...' : '🚀 Login'}
          </button>
        </form>
        <div style={{
          padding: '20px 0',
          borderTop: '1px solid #f0f0f0',
          marginTop: '10px'
        }}>
          <p style={{ color: '#666', marginBottom: '10px' }}>New to our store?</p>
          <button
            type="button"
            onClick={() => setShowLogin(false)}
            style={{
              background: 'transparent',
              border: '2px solid #667eea',
              color: '#667eea',
              padding: '10px 20px',
              borderRadius: '10px',
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
            Create Account
          </button>
        </div>
      </div>
    </div>
  );
}

export default LoginScreen;
