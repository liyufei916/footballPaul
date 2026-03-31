import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getCompetitions, getMyPredictions } from '../api/apiClient';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import { Trophy, Clock, Star, ArrowRight, Filter } from 'lucide-react';

const pointsConfig = {
  10: { label: '完全正确', color: 'text-gold-500', bg: 'bg-gold-500/20' },
  7: { label: '净胜球正确', color: 'text-pitch-400', bg: 'bg-pitch-500/20' },
  5: { label: '胜负正确', color: 'text-blue-400', bg: 'bg-blue-500/20' },
  3: { label: '部分正确', color: 'text-amber-400', bg: 'bg-amber-500/20' },
  0: { label: '未得分', color: 'text-slate-400', bg: 'bg-slate-500/20' },
};

export default function MyPredictionsPage() {
  const [predictions, setPredictions] = useState([]);
  const [competitions, setCompetitions] = useState([]);
  const [selectedCompetition, setSelectedCompetition] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [filter, setFilter] = useState('all');

  useEffect(() => {
    fetchCompetitions();
  }, []);

  useEffect(() => {
    fetchPredictions();
  }, [selectedCompetition]);

  const fetchCompetitions = async () => {
    try {
      const response = await getCompetitions();
      setCompetitions(response.data.competitions || []);
    } catch (err) {
      console.error('Failed to fetch competitions:', err);
    }
  };

  const fetchPredictions = async () => {
    setIsLoading(true);
    try {
      const response = await getMyPredictions(selectedCompetition);
      setPredictions(response.data.predictions || []);
    } catch (err) {
      console.error('Failed to fetch predictions:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const filteredPredictions = predictions.filter((p) => {
    if (filter === 'scored') return p.is_scored;
    if (filter === 'pending') return !p.is_scored;
    return true;
  });

  const filterCounts = {
    all: predictions.length,
    scored: predictions.filter((p) => p.is_scored).length,
    pending: predictions.filter((p) => !p.is_scored).length,
  };

  const totalPoints = predictions.reduce((sum, p) => sum + (p.points_earned || 0), 0);

  // Get competition name from prediction
  const getCompetitionName = (prediction) => {
    if (prediction.match?.competition) {
      return prediction.match.competition.name;
    }
    return '未知赛事';
  };

  return (
    <div className="max-w-5xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 pitch-gradient rounded-xl flex items-center justify-center">
            <Trophy className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-white">我的预测</h1>
            <p className="text-slate-400 text-sm">
              共 {predictions.length} 条预测，获得 {totalPoints} 积分
            </p>
          </div>
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
                ? 'bg-pitch-600 text-white'
                : 'bg-slate-800 text-slate-400 hover:bg-slate-700'
            }`}
          >
            全部赛事
          </button>
          {competitions.map((c) => (
            <button
              key={c.id}
              onClick={() => setSelectedCompetition(c.id)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                selectedCompetition === c.id
                  ? 'bg-pitch-600 text-white'
                  : 'bg-slate-800 text-slate-400 hover:bg-slate-700'
              }`}
            >
              {c.name}
            </button>
          ))}
        </div>
      </div>

      {/* Status Filter */}
      <div className="flex items-center gap-2 mb-6">
        {[
          { value: 'all', label: `全部 ${filterCounts.all}` },
          { value: 'scored', label: `已评分 ${filterCounts.scored}` },
          { value: 'pending', label: `待评分 ${filterCounts.pending}` },
        ].map((f) => (
          <button
            key={f.value}
            onClick={() => setFilter(f.value)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              filter === f.value
                ? 'bg-pitch-600 text-white'
                : 'bg-slate-800 text-slate-400 hover:bg-slate-700'
            }`}
          >
            {f.label}
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-slate-800 rounded-xl p-6 animate-pulse">
              <div className="h-6 bg-slate-700 rounded w-1/4 mb-4"></div>
              <div className="h-8 bg-slate-700 rounded w-1/2"></div>
            </div>
          ))}
        </div>
      ) : filteredPredictions.length === 0 ? (
        <div className="text-center py-16">
          <Trophy className="w-16 h-16 text-slate-600 mx-auto mb-4" />
          <p className="text-slate-400 text-lg">暂无预测记录</p>
          <Link
            to="/"
            className="inline-flex items-center gap-2 text-pitch-400 hover:text-pitch-300 mt-4"
          >
            去预测比赛 <ArrowRight className="w-4 h-4" />
          </Link>
        </div>
      ) : (
        <div className="space-y-4">
          {filteredPredictions.map((prediction) => {
            const pointsStyle = pointsConfig[prediction.points_earned] || pointsConfig[0];
            const match = prediction.match || {};

            return (
              <Link
                key={prediction.id}
                to={`/matches/${prediction.match_id}`}
                className="block bg-slate-800 border border-slate-700 rounded-xl p-6 hover:border-pitch-600 transition-all"
              >
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-4">
                    <span className="text-pitch-400 text-sm px-2 py-1 bg-pitch-900/50 rounded">
                      {getCompetitionName(prediction)}
                    </span>
                    <div>
                      <div className="text-lg font-semibold text-white">
                        {match.home_team} vs {match.away_team}
                      </div>
                      <div className="flex items-center gap-2 text-slate-400 text-sm mt-1">
                        <Clock className="w-4 h-4" />
                        {format(new Date(match.match_date), 'MM/dd HH:mm', { locale: zhCN })}
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center gap-4">
                    <div className="text-right">
                      <div className="text-slate-400 text-sm">你的预测</div>
                      <div className="text-xl font-bold text-white">
                        {prediction.predicted_home_score} - {prediction.predicted_away_score}
                      </div>
                    </div>

                    {prediction.is_scored && match.status === 'finished' && (
                      <>
                        <div className="text-slate-500 text-2xl">→</div>
                        <div className="text-right">
                          <div className="text-slate-400 text-sm">实际比分</div>
                          <div className="text-xl font-bold text-white">
                            {match.home_score} - {match.away_score}
                          </div>
                        </div>
                      </>
                    )}

                    <div className={`px-3 py-2 rounded-lg ${pointsStyle.bg}`}>
                      <div className={`text-lg font-bold ${pointsStyle.color}`}>
                        +{prediction.points_earned}
                      </div>
                    </div>
                  </div>
                </div>

                {prediction.is_scored && (
                  <div className="flex items-center gap-2 text-sm">
                    <Star className={`w-4 h-4 ${pointsStyle.color}`} />
                    <span className={pointsStyle.color}>{pointsStyle.label}</span>
                  </div>
                )}
              </Link>
            );
          })}
        </div>
      )}
    </div>
  );
}
