import { useState, useEffect } from 'react';
import { getCompetitions, getMatches } from '../api/apiClient';
import MatchCard from '../components/MatchCard';
import { Trophy, ArrowLeft, Filter } from 'lucide-react';

const statusFilters = [
  { value: '', label: '全部' },
  { value: 'pending', label: '待开始' },
  { value: 'ongoing', label: '进行中' },
  { value: 'finished', label: '已结束' },
];

export default function HomePage() {
  const [competitions, setCompetitions] = useState([]);
  const [selectedCompetition, setSelectedCompetition] = useState(null);
  const [matches, setMatches] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [statusFilter, setStatusFilter] = useState('');

  useEffect(() => {
    fetchCompetitions();
  }, []);

  useEffect(() => {
    if (selectedCompetition) {
      fetchMatches();
    }
  }, [selectedCompetition, statusFilter]);

  const fetchCompetitions = async () => {
    setIsLoading(true);
    try {
      const response = await getCompetitions();
      const comps = response.data.competitions || [];
      setCompetitions(comps);
      // Auto-select first competition if available
      if (comps.length > 0 && !selectedCompetition) {
        setSelectedCompetition(comps[0]);
      }
    } catch (err) {
      console.error('Failed to fetch competitions:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchMatches = async () => {
    try {
      const params = {};
      if (statusFilter) {
        params.status = statusFilter;
      }
      if (selectedCompetition) {
        params.competition_id = selectedCompetition.id;
      }
      const response = await getMatches(params);
      setMatches(response.data.matches || []);
    } catch (err) {
      console.error('Failed to fetch matches:', err);
    }
  };

  const handleCompetitionClick = (competition) => {
    setSelectedCompetition(competition);
    setStatusFilter('');
  };

  const handleBackToCompetitions = () => {
    setSelectedCompetition(null);
    setMatches([]);
  };

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 pitch-gradient rounded-xl flex items-center justify-center">
            <Trophy className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-white">
              {selectedCompetition ? selectedCompetition.name : '赛事列表'}
            </h1>
            <p className="text-slate-400 text-sm">
              {selectedCompetition ? '选择比赛，提交你的预测' : '点击赛事查看比赛'}
            </p>
          </div>
        </div>
      </div>

      {selectedCompetition && (
        <button
          onClick={handleBackToCompetitions}
          className="flex items-center gap-2 text-slate-400 hover:text-white mb-6 transition-colors"
        >
          <ArrowLeft className="w-5 h-5" />
          返回赛事列表
        </button>
      )}

      {/* Competition Grid - Show when no competition selected or on mobile */}
      {!selectedCompetition ? (
        <>
          {isLoading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {[1, 2, 3].map((i) => (
                <div key={i} className="bg-slate-800 rounded-xl p-6 animate-pulse">
                  <div className="h-8 bg-slate-700 rounded w-3/4 mb-4"></div>
                  <div className="h-4 bg-slate-700 rounded w-1/2"></div>
                </div>
              ))}
            </div>
          ) : competitions.length === 0 ? (
            <div className="text-center py-16">
              <Trophy className="w-16 h-16 text-slate-600 mx-auto mb-4" />
              <p className="text-slate-400 text-lg">暂无可用赛事</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {competitions.map((competition) => (
                <button
                  key={competition.id}
                  onClick={() => handleCompetitionClick(competition)}
                  className="bg-slate-800 border border-slate-700 rounded-xl p-6 hover:border-pitch-600 transition-all text-left group"
                >
                  <div className="flex items-center justify-between mb-4">
                    <div className="w-16 h-16 pitch-gradient rounded-xl flex items-center justify-center">
                      <Trophy className="w-8 h-8 text-white" />
                    </div>
                    <span className="text-pitch-400 text-sm group-hover:translate-x-1 transition-transform">
                      查看比赛 →
                    </span>
                  </div>
                  <h3 className="text-xl font-bold text-white mb-1">{competition.name}</h3>
                  <p className="text-slate-400 text-sm">{competition.code}</p>
                </button>
              ))}
            </div>
          )}
        </>
      ) : (
        <>
          {/* Competition Header */}
          <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 mb-6">
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 pitch-gradient rounded-xl flex items-center justify-center">
                <Trophy className="w-8 h-8 text-white" />
              </div>
              <div>
                <h2 className="text-xl font-bold text-white">{selectedCompetition.name}</h2>
                <p className="text-slate-400 text-sm">{selectedCompetition.code}</p>
              </div>
            </div>
          </div>

          {/* Status Filter */}
          <div className="flex items-center gap-2 mb-6">
            <Filter className="w-5 h-5 text-slate-400" />
            <div className="flex gap-2">
              {statusFilters.map((f) => (
                <button
                  key={f.value}
                  onClick={() => setStatusFilter(f.value)}
                  className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                    statusFilter === f.value
                      ? 'bg-pitch-600 text-white'
                      : 'bg-slate-800 text-slate-400 hover:bg-slate-700'
                  }`}
                >
                  {f.label}
                </button>
              ))}
            </div>
          </div>

          {/* Matches Grid */}
          {matches.length === 0 ? (
            <div className="text-center py-16">
              <Trophy className="w-16 h-16 text-slate-600 mx-auto mb-4" />
              <p className="text-slate-400 text-lg">暂无比赛</p>
              <p className="text-slate-500 text-sm mt-2">敬请期待更多比赛</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {matches.map((match) => (
                <MatchCard key={match.id} match={match} />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );
}
