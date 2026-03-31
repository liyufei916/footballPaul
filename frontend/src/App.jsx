import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Navbar from './components/Navbar';
import ProtectedRoute from './components/ProtectedRoute';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import MatchDetailPage from './pages/MatchDetailPage';
import MyPredictionsPage from './pages/MyPredictionsPage';
import LeaderboardPage from './pages/LeaderboardPage';
import AdminPage from './pages/AdminPage';

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <div className="min-h-screen bg-slate-900">
          <Navbar />
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route path="/matches/:id" element={<MatchDetailPage />} />
            <Route path="/leaderboard" element={<LeaderboardPage />} />
            <Route
              path="/my-predictions"
              element={
                <ProtectedRoute>
                  <MyPredictionsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/admin"
              element={
                <ProtectedRoute adminOnly>
                  <AdminPage />
                </ProtectedRoute>
              }
            />
          </Routes>

          <footer className="border-t border-slate-800 mt-16">
            <div className="max-w-7xl mx-auto px-4 py-8 text-center text-slate-500 text-sm">
              FootballPaul - 足球竞猜系统 © 2024
            </div>
          </footer>
        </div>
      </BrowserRouter>
    </AuthProvider>
  );
}
