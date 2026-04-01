import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getGroup, getGroupCompetitions, getGroupMatchPredictions } from '../api/apiClient';
import { Users, Trophy, ArrowLeft, Calendar, Clock, CheckCircle, XCircle } from 'lucide-react';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import { useAuth } from '../context/AuthContext';

export default function GroupPredictionsPage() {
  const { id, competitionId } = useParams();
  const { user } = useAuth();
  const [group, setGroup] = useState(null);
  const [competition, setCompetition] = useState(null);
  const [predictions, setPredictions] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, [id, competitionId]);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const [groupRes, predictionsRes] = await Promise.all([
        getGroup(id),
        getGroupMatchPredictions(id, competitionId),
      ]);
      setGroup(groupRes.data.group);
      setPredictions(predictionsRes.data.predictions || []);
      
      // Find competition info
      if (predictionsRes.data.predictions && predictionsRes.data.predictions.length > 0) {
        setCompetition(predictionsRes.data.predictions[0].match.competition);
      }
    } catch (err) {
      console.error('Failed to fetch data:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const renderPoints = (points, isScored) => {
    if (!isScored) {
      return <span className="text-slate-500 text-sm">待评分</span>;
    }
    return (
      <span className={`font-bold ${points > 0 ? 'text-pitch-400' : 'text-slate-500'}`}>
        {points}分
      </span>
    );
  };

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="bg-slate-800 rounded-xl p-8 animate-pulse">
          <div className="h-8 bg-slate-700 rounded w-1/3 mb-4"></div>
          <div className="h-4 bg-slate-700 rounded w-1/2"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-6">
        <Link to={`/groups/${id}`} className="text-slate-400 hover:text-white transition-colors">
          <ArrowLeft className="w-6 h-6" />
        </Link>
        <div>
          <h1 className="text-2xl font-bold text-white">
            {group?.name} - 预测一览
          </h1>
          {competition && (
            <p className="text-slate-400 text-sm mt-1">
              {competition.name}
            </p>
          )}
        </div>
      </div>

      {predictions.length === 0 ? (
        <div className="bg-slate-800 border border-slate-700 rounded-xl p-8 text-center">
          <Trophy className="w-12 h-12 text-slate-600 mx-auto mb-3" />
          <p className="text-slate-400 mb-1">暂无预测数据</p>
          <p className="text-slate-500 text-sm">组内成员尚未提交任何预测</p>
        </div>
      ) : (
        <div className="space-y-6">
          {predictions.map((item) => (
            <div key={item.match.id} className="bg-slate-800 border border-slate-700 rounded-xl p-5">
              {/* Match Header */}
              <div className="flex items-center justify-between mb-4 pb-4 border-b border-slate-700">
                <div className="flex items-center gap-2">
                  <span className={`px-2 py-0.5 rounded text-xs font-medium ${
                    item.match.status === 'finished' ? 'bg-slate-600 text-slate-300' :
                    item.match.status === 'ongoing' ? 'bg-green-600 text-white' :
                    'bg-amber-600 text-white'
                  }`}>
                    {item.match.status === 'finished' ? '已结束' : 
                     item.match.status === 'ongoing' ? '进行中' : '待开始'}
                  </span>
                  {item.match.status === 'finished' && item.match.home_score != null && (
                    <span className="text-white font-bold">
                      {item.match.home_score} - {item.match.away_score}
                    </span>
                  )}
                </div>
                <div className="flex items-center gap-1.5 text-slate-400 text-sm">
                  <Calendar className="w-4 h-4" />
                  {format(new Date(item.match.match_date), 'MM/dd HH:mm', { locale: zhCN })}
                </div>
              </div>

              {/* Teams */}
              <div className="flex items-center justify-between mb-4">
                <div className="text-lg font-semibold text-white text-center flex-1">
                  {item.match.home_team}
                </div>
                <div className="text-slate-500 text-xl px-4">vs</div>
                <div className="text-lg font-semibold text-white text-center flex-1">
                  {item.match.away_team}
                </div>
              </div>

              {/* Predictions */}
              {item.predictions.length === 0 ? (
                <div className="text-center py-4 text-slate-500 text-sm">
                  暂无成员预测
                </div>
              ) : (
                <div className="space-y-2">
                  <div className="grid grid-cols-3 gap-4 text-xs text-slate-400 px-2 mb-2">
                    <div>成员</div>
                    <div className="text-center">预测比分</div>
                    <div className="text-right">积分</div>
                  </div>
                  {item.predictions.map((pred) => (
                    <div
                      key={pred.user_id}
                      className="grid grid-cols-3 gap-4 items-center py-2 px-2 rounded-lg hover:bg-slate-700/30"
                    >
                      <div className="flex items-center gap-2">
                        <Users className="w-4 h-4 text-slate-500" />
                        <span className="text-white truncate">
                          {pred.username}
                          {pred.user_id === user?.id && (
                            <span className="text-pitch-400 text-xs ml-1">(你)</span>
                          )}
                        </span>
                      </div>
                      <div className="text-center text-white font-mono">
                        {pred.predicted_home_score} - {pred.predicted_away_score}
                      </div>
                      <div className="text-right">
                        {renderPoints(pred.points, pred.is_scored)}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
