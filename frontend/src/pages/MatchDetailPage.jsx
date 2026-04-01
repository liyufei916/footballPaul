import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import api from '../api/apiClient';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import { Calendar, Clock, Trophy, ArrowLeft, AlertCircle, CheckCircle } from 'lucide-react';

export default function MatchDetailPage() {
  const { id } = useParams();
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [match, setMatch] = useState(null);
  const [userPrediction, setUserPrediction] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const [homeScore, setHomeScore] = useState('');
  const [awayScore, setAwayScore] = useState('');

  useEffect(() => {
    fetchMatchData();
  }, [id]);

  const fetchMatchData = async () => {
    setIsLoading(true);
    try {
      const matchRes = await api.get(`/matches/${id}`);
      setMatch(matchRes.data);

      if (isAuthenticated) {
        try {
          const predRes = await api.get(`/predictions/my`);
          const prediction = predRes.data.predictions?.find((p) => p.match_id === parseInt(id));
          if (prediction) {
            setUserPrediction(prediction);
            setHomeScore(prediction.predicted_home_score.toString());
            setAwayScore(prediction.predicted_away_score.toString());
          }
        } catch {
          // No prediction found
        }
      }
    } catch (err) {
      setError('获取比赛信息失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmitPrediction = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setIsSubmitting(true);

    try {
      if (userPrediction) {
        await api.put(`/predictions/${userPrediction.id}`, {
          match_id: parseInt(id),
          predicted_home_score: parseInt(homeScore),
          predicted_away_score: parseInt(awayScore),
        });
        setSuccess('预测更新成功！');
      } else {
        await api.post('/predictions', {
          match_id: parseInt(id),
          predicted_home_score: parseInt(homeScore),
          predicted_away_score: parseInt(awayScore),
        });
        setSuccess('预测提交成功！');
      }
      fetchMatchData();
    } catch (err) {
      setError(err.response?.data?.error || '提交预测失败');
    } finally {
      setIsSubmitting(false);
    }
  };

  const isDeadlinePassed = match?.deadline && new Date(match.deadline) < new Date();
  const canPredict = match?.status === 'pending' && !isDeadlinePassed;

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="animate-pulse space-y-6">
          <div className="h-8 w-32 bg-slate-700 rounded"></div>
          <div className="h-64 bg-slate-800 rounded-xl"></div>
        </div>
      </div>
    );
  }

  if (!match) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8 text-center">
        <AlertCircle className="w-16 h-16 text-red-400 mx-auto mb-4" />
        <p className="text-white text-xl">比赛不存在</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <button
        onClick={() => navigate('/')}
        className="flex items-center gap-2 text-slate-400 hover:text-white mb-6 transition-colors"
      >
        <ArrowLeft className="w-5 h-5" />
        返回列表
      </button>

      <div className="bg-slate-800 border border-slate-700 rounded-xl overflow-hidden">
        {/* Header */}
        <div className="pitch-gradient p-6">
          <div className="flex items-center justify-between mb-4">
            <span className={`px-3 py-1 rounded-full text-sm font-medium ${
              match.status === 'pending' ? 'bg-amber-500/20 text-amber-300' :
              match.status === 'ongoing' ? 'bg-green-500/20 text-green-300' :
              'bg-slate-500/20 text-slate-300'
            }`}>
              {match.status === 'pending' ? '待开始' : match.status === 'ongoing' ? '进行中' : '已结束'}
            </span>
            <div className="flex items-center gap-2 text-white/80 text-sm">
              <Calendar className="w-4 h-4" />
              {format(new Date(match.match_date), 'yyyy年MM月dd日 HH:mm', { locale: zhCN })}
            </div>
          </div>

          <div className="flex items-center justify-between">
            <div className="text-center flex-1">
              <div className="text-3xl font-bold text-white mb-2">{match.home_team}</div>
              {match.status === 'finished' && (
                <div className="text-5xl font-bold text-white">
                  {match.home_score} - {match.away_score}
                </div>
              )}
            </div>
            <div className="text-2xl font-bold text-white/60 px-8">VS</div>
            <div className="text-center flex-1">
              <div className="text-3xl font-bold text-white mb-2">{match.away_team}</div>
            </div>
          </div>
        </div>

        {/* Prediction Section */}
        <div className="p-6">
          {match.status === 'finished' ? (
            <div className="text-center py-8">
              <Trophy className="w-12 h-12 text-gold-500 mx-auto mb-4" />
              <p className="text-slate-400">比赛已结束</p>
              {userPrediction && (
                <div className="mt-4 inline-block bg-slate-900 rounded-lg px-6 py-4">
                  <p className="text-slate-400 text-sm">你的预测</p>
                  <p className="text-2xl font-bold text-white mt-1">
                    {userPrediction.predicted_home_score} - {userPrediction.predicted_away_score}
                  </p>
                  <p className="text-pitch-400 font-medium mt-2">
                    获得 {userPrediction.points_earned} 积分
                  </p>
                </div>
              )}
            </div>
          ) : !isAuthenticated ? (
            <div className="text-center py-8">
              <p className="text-slate-400 mb-4">登录后可参与预测</p>
              <button
                onClick={() => navigate('/login')}
                className="bg-pitch-600 hover:bg-pitch-700 text-white px-6 py-3 rounded-lg transition-colors"
              >
                登录参与预测
              </button>
            </div>
          ) : canPredict ? (
            <div>
              <h3 className="text-lg font-semibold text-white mb-4">
                {userPrediction ? '修改你的预测' : '提交你的预测'}
              </h3>

              {error && (
                <div className="mb-4 p-3 bg-red-500/10 border border-red-500/30 rounded-lg flex items-center gap-2 text-red-400 text-sm">
                  <AlertCircle className="w-4 h-4" />
                  {error}
                </div>
              )}

              {success && (
                <div className="mb-4 p-3 bg-green-500/10 border border-green-500/30 rounded-lg flex items-center gap-2 text-green-400 text-sm">
                  <CheckCircle className="w-4 h-4" />
                  {success}
                </div>
              )}

              <form onSubmit={handleSubmitPrediction} className="space-y-4">
                <div className="flex items-center gap-4">
                  <div className="flex-1">
                    <label className="block text-slate-400 text-sm mb-2">{match.home_team}</label>
                    <input
                      type="number"
                      min="0"
                      value={homeScore}
                      onChange={(e) => setHomeScore(e.target.value)}
                      className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-3 text-white text-center text-2xl focus:outline-none focus:border-pitch-500"
                      required
                    />
                  </div>
                  <div className="text-2xl text-slate-500 pt-6">-</div>
                  <div className="flex-1">
                    <label className="block text-slate-400 text-sm mb-2">{match.away_team}</label>
                    <input
                      type="number"
                      min="0"
                      value={awayScore}
                      onChange={(e) => setAwayScore(e.target.value)}
                      className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-3 text-white text-center text-2xl focus:outline-none focus:border-pitch-500"
                      required
                    />
                  </div>
                </div>

                {canPredict && !userPrediction && (
                  <div className="flex items-center gap-2 text-amber-400 text-sm">
                    <Clock className="w-4 h-4" />
                    预测截止: {format(new Date(match.deadline), 'MM/dd HH:mm', { locale: zhCN })}
                  </div>
                )}

                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="w-full bg-pitch-600 hover:bg-pitch-700 text-white font-medium py-3 rounded-lg transition-colors disabled:opacity-50"
                >
                  {isSubmitting ? '提交中...' : userPrediction ? '更新预测' : '提交预测'}
                </button>
              </form>
            </div>
          ) : (
            <div className="text-center py-8">
              <AlertCircle className="w-12 h-12 text-amber-400 mx-auto mb-4" />
              <p className="text-amber-400">预测已截止</p>
              {userPrediction && (
                <p className="text-slate-400 mt-2">
                  你的预测: {userPrediction.predicted_home_score} - {userPrediction.predicted_away_score}
                </p>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
