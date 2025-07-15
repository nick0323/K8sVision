import React, { useState } from 'react';
import './LoginPage.css';

export default function LoginPage({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [remember, setRemember] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    if (remember) {
      localStorage.setItem('remembered_username', username);
      localStorage.setItem('remembered_password', password);
    } else {
      localStorage.removeItem('remembered_username');
      localStorage.removeItem('remembered_password');
    }
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    });
    if (res.ok) {
      const data = await res.json();
      localStorage.setItem('token', data.token);
      if (onLogin) onLogin();
      window.location.reload();
    } else {
      setError('Invalid username or password');
    }
  };

  React.useEffect(() => {
    const savedUser = localStorage.getItem('remembered_username') || '';
    const savedPwd = localStorage.getItem('remembered_password') || '';
    if (savedUser && savedPwd) {
      setUsername(savedUser);
      setPassword(savedPwd);
      setRemember(true);
    }
  }, []);

  return (
    <div className="login-bg">
      <form onSubmit={handleSubmit} className="login-form-card">
        <h2 style={{fontWeight:700, fontSize:28, marginBottom:24, letterSpacing:1, color:'#2563eb'}}>
          Vision For Kubernetes
        </h2>
        <input className="login-input" placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} autoFocus />
        <input className="login-input" type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
        <div className="login-remember-row">
          <input id="rememberMe" type="checkbox" checked={remember} onChange={e => setRemember(e.target.checked)} style={{marginRight:6}} />
          <label htmlFor="rememberMe" style={{fontSize:14, color:'#666', userSelect:'none'}}>Remember me</label>
        </div>
        <button className="login-btn" type="submit">Sign in</button>
        {error && <div className="login-error-tip">{error}</div>}
      </form>
    </div>
  );
} 