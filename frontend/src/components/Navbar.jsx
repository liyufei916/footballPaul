import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Trophy, LogOut, User, Settings, Users } from 'lucide-react';

export default function Navbar() {
  const { user, isAuthenticated, isAdmin, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <nav className="bg-slate-800 border-b border-slate-700 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <Link to="/" className="flex items-center gap-2">
            <div className="w-10 h-10 bg-pitch-600 rounded-full flex items-center justify-center">
              <Trophy className="w-6 h-6 text-white" />
            </div>
            <span className="text-xl font-bold text-white">FootballPaul</span>
          </Link>

          <div className="flex items-center gap-6">
            <Link
              to="/"
              className="text-slate-300 hover:text-pitch-400 transition-colors"
            >
              比赛
            </Link>
            <Link
              to="/leaderboard"
              className="text-slate-300 hover:text-pitch-400 transition-colors"
            >
              排行榜
            </Link>

            {isAuthenticated ? (
              <>
                <Link
                  to="/my-predictions"
                  className="text-slate-300 hover:text-pitch-400 transition-colors"
                >
                  我的预测
                </Link>
                <Link
                  to="/groups"
                  className="text-slate-300 hover:text-pitch-400 transition-colors flex items-center gap-1.5"
                >
                  <Users className="w-4 h-4" />
                  我的组
                </Link>
                {isAdmin && (
                  <Link
                    to="/admin"
                    className="text-gold-500 hover:text-gold-400 transition-colors flex items-center gap-1"
                  >
                    <Settings className="w-4 h-4" />
                    管理
                  </Link>
                )}
                <div className="flex items-center gap-4 pl-4 border-l border-slate-700">
                  <Link
                    to="/profile"
                    className="flex items-center gap-2 text-slate-300 hover:text-pitch-400"
                  >
                    <User className="w-4 h-4" />
                    <span>{user?.username}</span>
                    <span className="text-pitch-400 text-sm">({user?.total_points}分)</span>
                  </Link>
                  <button
                    onClick={handleLogout}
                    className="text-slate-400 hover:text-red-400 transition-colors"
                  >
                    <LogOut className="w-5 h-5" />
                  </button>
                </div>
              </>
            ) : (
              <div className="flex items-center gap-4">
                <Link
                  to="/login"
                  className="text-slate-300 hover:text-pitch-400 transition-colors"
                >
                  登录
                </Link>
                <Link
                  to="/register"
                  className="bg-pitch-600 hover:bg-pitch-700 text-white px-4 py-2 rounded-lg transition-colors"
                >
                  注册
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
