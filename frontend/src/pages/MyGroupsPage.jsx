import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getMyGroups } from '../api/apiClient';
import { Users, Plus, Trophy, Crown } from 'lucide-react';

export default function MyGroupsPage() {
  const [groups, setGroups] = useState([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchGroups();
  }, []);

  const fetchGroups = async () => {
    setIsLoading(true);
    try {
      const res = await getMyGroups();
      setGroups(res.data.groups || []);
    } catch (err) {
      console.error('Failed to fetch groups:', err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 bg-blue-600 rounded-xl flex items-center justify-center">
            <Users className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-white">我的组</h1>
            <p className="text-slate-400 text-sm">和好友一起竞猜排行</p>
          </div>
        </div>
        <div className="flex gap-3">
          <Link
            to="/groups/join"
            className="flex items-center gap-2 bg-slate-700 hover:bg-slate-600 text-white px-4 py-2 rounded-lg transition-colors"
          >
            加入组
          </Link>
          <Link
            to="/groups/new"
            className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors"
          >
            <Plus className="w-4 h-4" />
            创建组
          </Link>
        </div>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-slate-800 rounded-xl p-6 animate-pulse">
              <div className="h-6 bg-slate-700 rounded w-1/3 mb-3"></div>
              <div className="h-4 bg-slate-700 rounded w-1/4"></div>
            </div>
          ))}
        </div>
      ) : groups.length === 0 ? (
        <div className="text-center py-20">
          <Users className="w-20 h-20 text-slate-600 mx-auto mb-4" />
          <p className="text-slate-400 text-lg mb-2">你还没有加入任何组</p>
          <p className="text-slate-500 text-sm mb-6">创建或加入组，和朋友一起比拼预测实力</p>
          <div className="flex justify-center gap-4">
            <Link
              to="/groups/join"
              className="bg-slate-700 hover:bg-slate-600 text-white px-6 py-2 rounded-lg transition-colors"
            >
              加入组
            </Link>
            <Link
              to="/groups/new"
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg transition-colors"
            >
              创建组
            </Link>
          </div>
        </div>
      ) : (
        <div className="space-y-4">
          {groups.map((group) => (
            <Link
              key={group.id}
              to={`/groups/${group.id}`}
              className="block bg-slate-800 border border-slate-700 hover:border-slate-600 rounded-xl p-5 transition-colors"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className={`w-12 h-12 rounded-xl flex items-center justify-center ${
                    group.role === 'admin' ? 'bg-gold-600' : 'bg-blue-600'
                  }`}>
                    {group.role === 'admin' ? (
                      <Crown className="w-6 h-6 text-white" />
                    ) : (
                      <Users className="w-6 h-6 text-white" />
                    )}
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <h3 className="font-semibold text-white">{group.name}</h3>
                      {group.role === 'admin' && (
                        <span className="text-xs bg-gold-600/20 text-gold-400 px-2 py-0.5 rounded">组长</span>
                      )}
                    </div>
                    <div className="flex items-center gap-4 text-slate-400 text-sm">
                      <span>{group.member_count} 名成员</span>
                      <span className="text-slate-600">·</span>
                      <span>{group.competition_count} 个追踪赛事</span>
                    </div>
                  </div>
                </div>
                <div className="text-slate-500">
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                  </svg>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
