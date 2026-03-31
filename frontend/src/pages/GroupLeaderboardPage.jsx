import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getGroupLeaderboard, getGroup, getGroupCompetitions } from '../api/apiClient';
import { ArrowLeft, Trophy, Medal, User } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

export default function GroupLeaderboardPage() {
  const { id, competitionId } = useParams();
  const { user } = useAuth();
  const [group, setGroup] = useState(null);
  const [competition, setCompetition] = useState(null);
  const [competitions, setCompetitions] = useState([]);
  const [leaderboard, setLeaderboard] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedCompId, setSelectedCompId] = useState(parseInt(competitionId));

  useEffect(() => {
    fetchGroupData();
  }, [id]);

  useEffect(() => {
    if (selectedCompId) {
      fetchLeaderboard();
    }
  }, [selectedCompId]);

  const fetchGroupData = async () => {
    try {
      const [groupRes, compsRes] = await Promise.all([
        getGroup(id),
        getGroupCompetitions(id),
      ]);
      setGroup(groupRes.data.group);
      setCompetitions(compsRes.data.competitions || []);
      if (compsRes.data.competitions?.length > 0 && !selectedCompId) {
        setSelectedCompId(compsRes.data.competitions[0].id);
      }
    } catch (err) {
      console.error('Failed to fetch group data:', err);
    }
  };

  const fetchLeaderboard = async () => {
    setIsLoading(true);
    try {
      const res = await getGroupLeaderboard(id, selectedCompId);
      setLeaderboard(res.data.leaderboard || []);
      const comp = res.data.competition;
      if (comp && comp.id) {
        setCompetition(comp);
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

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-8">
        <Link to={`/groups/${id}`} className="text-slate-400 hover:text-white transition-colors">
          <ArrowLeft className="w-6 h-6" />
        </Link>
        <div>
          <h1 className="text-xl font-bold text-white">
            {competition?.name || '组排行榜'}
          </h1>
          {group && <p className="text-slate-400 text-sm">{group.name} - 组内排行</p>}
        </div>
      </div>

      {/* Competition Selector */}
      {competitions.length > 1 && (
        <div className="mb-6">
          <div className="flex flex-wrap gap-2">
            {competitions.map((c) => (
              <button
                key={c.id}
                onClick={() => setSelectedCompId(c.id)}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  selectedCompId === c.id
                    ? 'bg-blue-600 text-white'
                    : 'bg-slate-800 text-slate-400 hover:bg-slate-700'
                }`}
              >
                {c.name}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Leaderboard */}
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
      ) : leaderboard.length === 0 ? (
        <div className="text-center py-16">
          <Trophy className="w-16 h-16 text-slate-600 mx-auto mb-4" />
          <p className="text-slate-400 text-lg">暂无排行数据</p>
        </div>
      ) : (
        <div className="space-y-3">
          {leaderboard.map((entry) => (
            <div
              key={entry.user_id}
              className={`flex items-center gap-4 p-4 rounded-xl border ${getRankBg(entry.rank)} ${
                user?.id === entry.user_id ? 'ring-2 ring-blue-500' : ''
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
                    <span className="text-xs text-blue-400">(你)</span>
                  )}
                </div>
                <div className="text-slate-400 text-sm flex items-center gap-3">
                  <span>{entry.predictions_count} 次预测</span>
                  <span>·</span>
                  <span>{entry.exact_scores} 场全对</span>
                </div>
              </div>
              <div className="text-right">
                <div className="text-2xl font-bold text-blue-400">{entry.total_points}</div>
                <div className="text-slate-400 text-xs">总积分</div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
