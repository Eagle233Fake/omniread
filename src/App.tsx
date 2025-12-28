import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/auth/Login';
import Register from './pages/auth/Register';
import Library from './pages/library/Library';
import Insight from './pages/insight/Insight';
import Profile from './pages/profile/Profile';
import Layout from './layout/Layout';
import { AuthProvider, useAuth } from './context/AuthContext';
import Reader from './pages/reader/Reader';

const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
};

const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
        <Route index element={<Navigate to="/library" replace />} />
        <Route path="library" element={<Library />} />
        <Route path="insight" element={<Insight />} />
        <Route path="profile" element={<Profile />} />
      </Route>
      <Route path="/read/:id" element={<ProtectedRoute><Reader /></ProtectedRoute>} />
    </Routes>
  );
};

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;
