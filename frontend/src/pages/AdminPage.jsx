import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { getCompetitions, createMatch, enterResult, getMatches } from '../api/apiClient';
import { Settings, Plus, Trophy, AlertCircle, CheckCircle, Edit2 } from 'lucide-react';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

export default function AdminPage() {
  const { isAdmin } = useAuth();
  const [competitions, setCompetitions] = useState([]);
  const [matches, setMatches] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  // Create match form
  const [selectedCompetition, setSelectedCompetition] = useState('');
  const [homeTeam, setHomeTeam] = useState('');
  const [awayTeam, setAwayTeam] = useState('');
  const [matchDate, setMatchDate] = useState('');
  const [deadline, setDeadline] = useState('');
  const [createError, setCreateError] = useState('');
  const [createSuccess, setCreateSuccess] = useState('');

  // Edit match form
  const [editingMatch, setEditingMatch] = useState(null);
  const [editHomeScore, setEditHomeScore] = useState('');
  const [editAwayScore, setEditAwayScore] = useState('');
  const [editError, setEditError] = useState('');
  const [editSuccess, setEditSuccess] = useState('');

  // Enter result form
  const [selectedMatch, setSelectedMatch] = useState('');
  const [homeScore, setHomeScore] = useState('');
  const [awayScore, setAwayScore] = useState('');
  const [resultError, setResultError] = useState('');
  const [resultSuccess, setResultSuccess] = useState('');

  useEffect(() => {
    fetchCompetitions();
    fetchMatches();
  }, []);

  const fetchCompetitions = async () => {
    try {
      const response = await getCompetitions();
      setCompetitions(response.data.competitions || []);
      if (response.data.competitions?.length > 0) {
        setSelectedCompetition(response.data.competitions[0].id);
      }
    } catch (err) {
      console.error('Failed to fetch competitions:', err);
    }
  };

  const fetchMatches = async () => {
    try {
      const response = await getMatches({ limit: 100 });
      setMatches(response.data.matches || []);
    } catch (err) {
      console.error('Failed to fetch matches:', err);
    }
  };

  const handleCreateMatch = async (e) => {
    e.preventDefault();
    setCreateError('');
    setCreateSuccess('');
    setIsLoading(true);

    try {
      await createMatch({
        competition_id: parseInt(selectedCompetition),
        home_team: homeTeam,
        away_team: awayTeam,
        match_date: new Date(matchDate).toISOString(),
        deadline: new Date(deadline).toISOString(),
      });
      setCreateSuccess('比赛创建成功！');
      setHomeTeam('');
      setAwayTeam('');
      setMatchDate('');
      setDeadline('');
      fetchMatches();
    } catch (err) {
      setCreateError(err.response?.data?.error || '创建失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleEnterResult = async (e) => {
    e.preventDefault();
    setResultError('');
    setResultSuccess('');
    setIsLoading(true);

    try {
      await enterResult(selectedMatch, {
        home_score: parseInt(homeScore),
        away_score: parseInt(awayScore),
      });
      setResultSuccess('比分录入成功！评分已完成。');
      setSelectedMatch('');
      setHomeScore('');
      setAwayScore('');
      fetchMatches();
    } catch (err) {
      setResultError(err.response?.data?.error || '录入失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleEditMatch = async (e) => {
    e.preventDefault();
    setEditError('');
    setEditSuccess('');
    setIsLoading(true);

    try {
      await enterResult(editingMatch.id, {
        home_score: parseInt(editHomeScore),
        away_score: parseInt(editAwayScore),
      });
      setEditSuccess('比分更新成功！评分已完成。');
      setEditingMatch(null);
      setEditHomeScore('');
      setEditAwayScore('');
      fetchMatches();
    } catch (err) {
      setEditError(err.response?.data?.error || '更新失败');
    } finally {
      setIsLoading(false);
    }
  };

  const openEditModal = (match) => {
    setEditingMatch(match);
    setEditHomeScore(match.home_score?.toString() || '');
    setEditAwayScore(match.away_score?.toString() || '');
    setEditError('');
    setEditSuccess('');
  };

  const pendingMatches = matches.filter((m) => m.status === 'pending');

  if (!isAdmin) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8 text-center">
        <AlertCircle className="w-16 h-16 text-red-400 mx-auto mb-4" />
        <h2 className="text-xl font-bold text-white mb-2">无权限访问</h2>
        <p className="text-slate-400">您不是管理员，无法访问此页面。</p>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-8">
        <div className="w-12 h-12 bg-gold-500 rounded-xl flex items-center justify-center">
          <Settings className="w-6 h-6 text-white" />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-white">管理面板</h1>
          <p className="text-slate-400 text-sm">管理赛事、比赛和比分</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Create Match */}
        <div className="bg-slate-800 border border-slate-700 rounded-xl p-6">
          <div className="flex items-center gap-2 mb-6">
            <Plus className="w-5 h-5 text-pitch-400" />
            <h2 className="text-lg font-semibold text-white">创建比赛</h2>
          </div>

          {createError && (
            <div className="mb-4 p-3 bg-red-500/10 border border-red-500/30 rounded-lg flex items-center gap-2 text-red-400 text-sm">
              <AlertCircle className="w-4 h-4" />
              {createError}
            </div>
          )}

          {createSuccess && (
            <div className="mb-4 p-3 bg-green-500/10 border border-green-500/30 rounded-lg flex items-center gap-2 text-green-400 text-sm">
              <CheckCircle className="w-4 h-4" />
              {createSuccess}
            </div>
          )}

          <form onSubmit={handleCreateMatch} className="space-y-4">
            <div>
              <label className="block text-slate-300 text-sm mb-1">赛事</label>
              <select
                value={selectedCompetition}
                onChange={(e) => setSelectedCompetition(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-pitch-500"
                required
              >
                <option value="">选择赛事...</option>
                {competitions.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.name}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-slate-300 text-sm mb-1">主队</label>
              <input
                type="text"
                value={homeTeam}
                onChange={(e) => setHomeTeam(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-pitch-500"
                placeholder="如: 巴西"
                required
              />
            </div>

            <div>
              <label className="block text-slate-300 text-sm mb-1">客队</label>
              <input
                type="text"
                value={awayTeam}
                onChange={(e) => setAwayTeam(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-pitch-500"
                placeholder="如: 阿根廷"
                required
              />
            </div>

            <div>
              <label className="block text-slate-300 text-sm mb-1">比赛时间</label>
              <input
                type="datetime-local"
                value={matchDate}
                onChange={(e) => setMatchDate(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-pitch-500"
                required
              />
            </div>

            <div>
              <label className="block text-slate-300 text-sm mb-1">预测截止时间</label>
              <input
                type="datetime-local"
                value={deadline}
                onChange={(e) => setDeadline(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-pitch-500"
                required
              />
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-pitch-600 hover:bg-pitch-700 text-white font-medium py-2 rounded-lg transition-colors disabled:opacity-50"
            >
              {isLoading ? '创建中...' : '创建比赛'}
            </button>
          </form>
        </div>

        {/* Enter Result */}
        <div className="bg-slate-800 border border-slate-700 rounded-xl p-6">
          <div className="flex items-center gap-2 mb-6">
            <Trophy className="w-5 h-5 text-gold-500" />
            <h2 className="text-lg font-semibold text-white">录入比分</h2>
          </div>

          {resultError && (
            <div className="mb-4 p-3 bg-red-500/10 border border-red-500/30 rounded-lg flex items-center gap-2 text-red-400 text-sm">
              <AlertCircle className="w-4 h-4" />
              {resultError}
            </div>
          )}

          {resultSuccess && (
            <div className="mb-4 p-3 bg-green-500/10 border border-green-500/30 rounded-lg flex items-center gap-2 text-green-400 text-sm">
              <CheckCircle className="w-4 h-4" />
              {resultSuccess}
            </div>
          )}

          <form onSubmit={handleEnterResult} className="space-y-4">
            <div>
              <label className="block text-slate-300 text-sm mb-1">选择比赛</label>
              <select
                value={selectedMatch}
                onChange={(e) => setSelectedMatch(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-pitch-500"
                required
              >
                <option value="">选择一场比赛...</option>
                {pendingMatches.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.competition?.name}: {m.home_team} vs {m.away_team} ({format(new Date(m.match_date), 'MM/dd HH:mm', { locale: zhCN })})
                  </option>
                ))}
              </select>
            </div>

            <div className="flex items-center gap-4">
              <div className="flex-1">
                <label className="block text-slate-300 text-sm mb-1">主队得分</label>
                <input
                  type="number"
                  min="0"
                  value={homeScore}
                  onChange={(e) => setHomeScore(e.target.value)}
                  className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white text-center focus:outline-none focus:border-pitch-500"
                  required
                />
              </div>
              <div className="text-slate-500 pt-6">-</div>
              <div className="flex-1">
                <label className="block text-slate-300 text-sm mb-1">客队得分</label>
                <input
                  type="number"
                  min="0"
                  value={awayScore}
                  onChange={(e) => setAwayScore(e.target.value)}
                  className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white text-center focus:outline-none focus:border-pitch-500"
                  required
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="w-full bg-gold-500 hover:bg-gold-600 text-white font-medium py-2 rounded-lg transition-colors disabled:opacity-50"
            >
              {isLoading ? '录入中...' : '录入比分'}
            </button>
          </form>
        </div>
      </div>

      {/* Recent Matches */}
      <div className="mt-8">
        <h2 className="text-lg font-semibold text-white mb-4">所有比赛</h2>
        <div className="bg-slate-800 border border-slate-700 rounded-xl overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-slate-900">
                <tr>
                  <th className="px-4 py-3 text-left text-xs font-medium text-slate-400 uppercase">赛事</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-slate-400 uppercase">比赛</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-slate-400 uppercase">时间</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-slate-400 uppercase">状态</th>
                  <th className="px-4 py-3 text-left text-xs font-medium text-slate-400 uppercase">比分</th>
                  <th className="px-4 py-3 text-right text-xs font-medium text-slate-400 uppercase">操作</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-700">
                {matches.map((m) => (
                  <tr key={m.id} className="hover:bg-slate-700/50">
                    <td className="px-4 py-3 text-sm text-slate-400">{m.competition?.name || '-'}</td>
                    <td className="px-4 py-3 text-sm text-white font-medium">
                      {m.home_team} vs {m.away_team}
                    </td>
                    <td className="px-4 py-3 text-sm text-slate-400">
                      {format(new Date(m.match_date), 'MM/dd HH:mm', { locale: zhCN })}
                    </td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-0.5 rounded text-xs font-medium ${
                        m.status === 'finished' ? 'bg-slate-600 text-slate-300' :
                        m.status === 'ongoing' ? 'bg-green-600 text-white' :
                        'bg-amber-600 text-white'
                      }`}>
                        {m.status === 'finished' ? '已结束' : m.status === 'ongoing' ? '进行中' : '待开始'}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-white font-bold">
                      {m.status === 'finished' ? `${m.home_score} - ${m.away_score}` : '-'}
                    </td>
                    <td className="px-4 py-3 text-right">
                      {m.status !== 'finished' && (
                        <button
                          onClick={() => openEditModal(m)}
                          className="text-pitch-400 hover:text-pitch-300 p-1"
                          title="录入/编辑比分"
                        >
                          <Edit2 className="w-4 h-4" />
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Edit Modal */}
      {editingMatch && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold text-white mb-4">编辑比分</h3>
            <p className="text-slate-400 text-sm mb-4">
              {editingMatch.home_team} vs {editingMatch.away_team}
            </p>

            {editError && (
              <div className="mb-4 p-3 bg-red-500/10 border border-red-500/30 rounded-lg flex items-center gap-2 text-red-400 text-sm">
                <AlertCircle className="w-4 h-4" />
                {editError}
              </div>
            )}

            {editSuccess && (
              <div className="mb-4 p-3 bg-green-500/10 border border-green-500/30 rounded-lg flex items-center gap-2 text-green-400 text-sm">
                <CheckCircle className="w-4 h-4" />
                {editSuccess}
              </div>
            )}

            <form onSubmit={handleEditMatch} className="space-y-4">
              <div className="flex items-center gap-4">
                <div className="flex-1">
                  <label className="block text-slate-300 text-sm mb-1">{editingMatch.home_team}</label>
                  <input
                    type="number"
                    min="0"
                    value={editHomeScore}
                    onChange={(e) => setEditHomeScore(e.target.value)}
                    className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white text-center text-xl focus:outline-none focus:border-pitch-500"
                    required
                  />
                </div>
                <div className="text-slate-500 pt-6 text-xl">-</div>
                <div className="flex-1">
                  <label className="block text-slate-300 text-sm mb-1">{editingMatch.away_team}</label>
                  <input
                    type="number"
                    min="0"
                    value={editAwayScore}
                    onChange={(e) => setEditAwayScore(e.target.value)}
                    className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white text-center text-xl focus:outline-none focus:border-pitch-500"
                    required
                  />
                </div>
              </div>

              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={() => setEditingMatch(null)}
                  className="flex-1 bg-slate-700 hover:bg-slate-600 text-white font-medium py-2 rounded-lg transition-colors"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={isLoading}
                  className="flex-1 bg-gold-500 hover:bg-gold-600 text-white font-medium py-2 rounded-lg transition-colors disabled:opacity-50"
                >
                  {isLoading ? '保存中...' : '保存'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
