import { Link } from 'react-router-dom';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import { Calendar, Clock, CheckCircle, Play, Trophy } from 'lucide-react';

const statusConfig = {
  pending: {
    label: '待开始',
    color: 'bg-amber-500/20 text-amber-400 border-amber-500/30',
    icon: Clock,
  },
  ongoing: {
    label: '进行中',
    color: 'bg-green-500/20 text-green-400 border-green-500/30',
    icon: Play,
  },
  finished: {
    label: '已结束',
    color: 'bg-slate-500/20 text-slate-400 border-slate-500/30',
    icon: Trophy,
  },
};

export default function MatchCard({ match }) {
  const status = statusConfig[match.status] || statusConfig.pending;
  const StatusIcon = status.icon;

  const isDeadlinePassed = match.deadline && new Date(match.deadline) < new Date();
  const canPredict = match.status === 'pending' && !isDeadlinePassed;

  return (
    <Link to={`/matches/${match.id}`}>
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-5 hover:border-pitch-600 transition-all hover:shadow-lg hover:shadow-pitch-900/20">
        {/* Header row: status badge + score (for finished) */}
        <div className="flex items-center justify-between mb-4">
          <span className={`flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border ${status.color}`}>
            <StatusIcon className="w-3.5 h-3.5" />
            {status.label}
          </span>
          {match.status === 'finished' && match.home_score != null && (
            <span className="text-white font-bold text-sm">
              {match.home_score} - {match.away_score}
            </span>
          )}
        </div>

        {/* Teams + score row */}
        <div className="flex items-center justify-between">
          <div className="flex-1 text-center">
            <div className="text-lg font-semibold text-white">{match.home_team}</div>
          </div>

          {match.status === 'finished' && match.home_score != null ? (
            <div className="text-center px-4">
              <div className="text-2xl font-bold text-white">
                {match.home_score} - {match.away_score}
              </div>
            </div>
          ) : (
            <div className="text-2xl font-bold text-slate-600 px-4">vs</div>
          )}

          <div className="flex-1 text-center">
            <div className="text-lg font-semibold text-white">{match.away_team}</div>
          </div>
        </div>

        {/* Date row */}
        <div className="flex items-center justify-end mt-3">
          <div className="flex items-center gap-1.5 text-slate-400 text-sm">
            <Calendar className="w-4 h-4" />
            {format(new Date(match.match_date), 'MM/dd HH:mm', { locale: zhCN })}
          </div>
          {match.status === 'pending' && (
            <span className={`ml-3 flex items-center gap-1.5 text-sm ${isDeadlinePassed ? 'text-red-400' : 'text-amber-400'}`}>
              <Clock className="w-4 h-4" />
              {isDeadlinePassed ? '已截止' : `截止 ${format(new Date(match.deadline), 'MM/dd HH:mm', { locale: zhCN })}`}
            </span>
          )}
        </div>

        {canPredict && (
          <div className="mt-4 pt-4 border-t border-slate-700">
            <span className="text-pitch-400 text-sm flex items-center gap-1">
              <CheckCircle className="w-4 h-4" />
              点击提交预测
            </span>
          </div>
        )}
      </div>
    </Link>
  );
}
