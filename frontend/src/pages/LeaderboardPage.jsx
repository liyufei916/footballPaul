import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { getCompetitions, getLeaderboard, getMyRank } from '../api/apiClient';
import { Trophy, Medal, User, Filter } from 'lucide-react';

export default function LeaderboardPage() {
  const { user, isAuthenticated } = useAuth();
  const [rankings, setRankings] = useState([]);
  const [myRank, setMyRank] = useState(null);
  const [competitions, setCompetitions] = useState([]);
  const [selectedCompetition, setSelectedCompetition] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchCompetitions();
  }, []);

  useEffect(() => {
    fetchLeaderboard();
  }, [selectedCompetition]);

  const fetchCompetitions = async () => {
    try {
      const response = await getCompetitions();
      setCompetitions(response.data.competitions || []);
    } catch (err) {
      console.error('Failed to fetch competitions:', err);
    }
  };

  const fetchLeaderboard = async () => {
    setIsLoading(true);
    try {
      const response = await getLeaderboard(selectedCompetition);
      setRankings(response.data.rankings || []);

      if (isAuthenticated) {
        try {
          const rankRes = await getMyRank(selectedCompetition);
          setMyRank(rankRes.data.rank);
        } catch {
          setMyRank(null);
        }
      }
    } catch (err) {
      console.error('Failed to fetch leaderboard:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const getRankIcon = (rank) => {
    if (rank === 1) return <Medal className="w-6 h-6 text-gold-500" />;
    if (rank === 2) return <Medal className="w-6 h-6 text-slate-300" />;
    if (rank === 3) return <Medal className="w-6 h-6 text-amber-600" />;
    return <span className="text-slate-400 font-medium">#{rank}</span>;
  };

  const getRankBg = (rank) => {
    if (rank === 1) return 'bg-gold-500/10 border-gold-500/30';
    if (rank === 2) return 'bg-slate-500/10 border-slate-500/30';
    if (rank === 3) return 'bg-amber-600/10 border-amber-600/30';
    return 'bg-slate-800 border-slate-700';
  };

  const selectedCompetitionName = selectedCompetition
    ? competitions.find((c) => c.id === selectedCompetition)?.name
    : null;

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-8">
        <div className="w-12 h-12 gold-gradient rounded-xl flex items-center justify-center">
          <Trophy className="w-6 h-6 text-white" />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-white">
            {selectedCompetitionName ? `${selectedCompetitionName} 排行榜` : '排行榜'}
          </h1>
          <p className="text-slate-400 text-sm">
            {selectedCompetitionName ? '赛事积分排行' : '全局积分排行'}
          </p>
        </div>
      </div>

      {/* Competition Filter */}
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-3">
          <Filter className="w-5 h-5 text-slate-400" />
          <span className="text-slate-300 text-sm">赛事筛选</span>
        </div>
        <div className="flex flex-wrap gap-2">
          <button
            onClick={() => setSelectedCompetition(null)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              selectedCompetition === null
                ? 'bg-gold-600 text-white'
                : 'bg-slate-800 text-slate-400 hover:bg-slate-700'
            }`}
          >
            全局排行
          </button>
          {competitions.map((c) => (
            <button
              key={c.id}
              onClick={() => setSelectedCompetition(c.id)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                selectedCompetition === c.id
                  ? 'bg-gold-600 text-white'
                  : 'bg-slate-800 text-slate-400 hover:bg-slate-700'
              }`}
            >
              {c.name}
            </button>
          ))}
        </div>
      </div>

      {myRank && (
        <div className="mb-6 p-4 bg-pitch-900/50 border border-pitch-600/30 rounded-xl">
          <div className="flex items-center gap-3">
            <User className="w-5 h-5 text-pitch-400" />
            <span className="text-slate-300">你的排名:</span>
            <span className="text-pitch-400 font-bold text-xl">第 {myRank} 名</span>
          </div>
        </div>
      )}

      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="bg-slate-800 rounded-xl p-6 animate-pulse">
              <div className="flex items-center gap-4">
                <div className="w-10 h-10 bg-slate-700 rounded-full"></div>
                <div className="flex-1">
                  <div className="h-5 bg-slate-700 rounded w-1/3 mb-2"></div>
                  <div className="h-4 bg-slate-700 rounded w-1/4"></div>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : rankings.length === 0 ? (
        <div className="text-center py-16">
          <Trophy className="w-16 h-16 text-slate-600 mx-auto mb-4" />
          <p className="text-slate-400 text-lg">
            {selectedCompetition ? '该赛事暂无排行数据' : '暂无排行数据'}
          </p>
        </div>
      ) : (
        <div className="space-y-3">
          {rankings.map((entry) => (
            <div
              key={entry.user_id}
              className={`flex items-center gap-4 p-4 rounded-xl border ${getRankBg(entry.rank)} ${
                user?.id === entry.user_id ? 'ring-2 ring-pitch-500' : ''
              }`}
            >
              <div className="w-10 h-10 flex items-center justify-center">
                {getRankIcon(entry.rank)}
              </div>

              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <span className="font-semibold text-white">{entry.username}</span>
                  {entry.rank <= 3 && <span className="text-xs">🏆</span>}
                  {user?.id === entry.user_id && (
                    <span className="text-xs text-pitch-400">(你)</span>
                  )}
                </div>
                <div className="text-slate-400 text-sm">
                  {entry.predictions_count} 次预测
                </div>
              </div>

              <div className="text-right">
                <div className="text-2xl font-bold text-pitch-400">{entry.total_points}</div>
                <div className="text-slate-400 text-xs">总积分</div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
