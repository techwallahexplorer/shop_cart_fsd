import React, { useState } from 'react';
import LoginScreen from './components/LoginScreen';
import RegisterScreen from './components/RegisterScreen';
import ItemsListScreen from './components/ItemsListScreen';

function App() {
  const [userToken, setUserToken] = useState(localStorage.getItem('userToken') || '');
  const [username, setUsername] = useState(localStorage.getItem('username') || '');
  const [showLogin, setShowLogin] = useState(true);

  // Routing: show login/register if not logged in, else show items
  return (
    <div>
      {!userToken ? (
        showLogin ? (
          <LoginScreen setUserToken={setUserToken} setUsername={setUsername} setShowLogin={setShowLogin} />
        ) : (
          <RegisterScreen setShowLogin={setShowLogin} />
        )
      ) : (
        <ItemsListScreen userToken={userToken} username={username} setUserToken={setUserToken} />
      )}
    </div>
  );
}

export default App;
