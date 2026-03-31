import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { getGroup, getGroupCompetitions, leaveGroup, deleteGroup, getMyGroups } from '../api/apiClient';
import { Users, Trophy, ArrowLeft, Crown, Plus, Trash2, LogOut } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

export default function GroupDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [group, setGroup] = useState(null);
  const [competitions, setCompetitions] = useState([]);
  const [myGroups, setMyGroups] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isAdmin, setIsAdmin] = useState(false);
  const [isOwner, setIsOwner] = useState(false);

  useEffect(() => {
    fetchData();
  }, [id]);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const [groupRes, compsRes, myGroupsRes] = await Promise.all([
        getGroup(id),
        getGroupCompetitions(id),
        getMyGroups(),
      ]);
      setGroup(groupRes.data.group);
      setCompetitions(compsRes.data.competitions || []);
      setMyGroups(myGroupsRes.data.groups || []);
    } catch (err) {
      console.error('Failed to fetch group:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (group && user) {
      const membership = myGroups.find((g) => g.id === parseInt(id));
      setIsAdmin(membership?.role === 'admin');
      setIsOwner(group.owner_id === user.id);
    }
  }, [group, user, myGroups]);

  const handleLeave = async () => {
    if (!window.confirm('确定要离开该组吗？')) return;
    try {
      await leaveGroup(id);
      navigate('/groups');
    } catch (err) {
      alert(err.response?.data?.error || '离开失败');
    }
  };

  const handleDelete = async () => {
    if (!window.confirm('确定要解散该组吗？此操作不可恢复！')) return;
    try {
      await deleteGroup(id);
      navigate('/groups');
    } catch (err) {
      alert(err.response?.data?.error || '解散失败');
    }
  };

  if (isLoading) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-8">
        <div className="bg-slate-800 rounded-xl p-8 animate-pulse">
          <div className="h-8 bg-slate-700 rounded w-1/3 mb-4"></div>
          <div className="h-4 bg-slate-700 rounded w-1/2"></div>
        </div>
      </div>
    );
  }

  if (!group) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-8 text-center">
        <p className="text-slate-400">组队不存在或你无权访问</p>
        <Link to="/groups" className="text-blue-400 hover:underline mt-4 inline-block">返回我的组</Link>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-6">
        <Link to="/groups" className="text-slate-400 hover:text-white transition-colors">
          <ArrowLeft className="w-6 h-6" />
        </Link>
      </div>

      {/* Group Header */}
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-6 mb-6">
        <div className="flex items-center gap-4 mb-4">
          <div className={`w-14 h-14 rounded-xl flex items-center justify-center ${isOwner ? 'bg-gold-600' : 'bg-blue-600'}`}>
            {isOwner ? <Crown className="w-7 h-7 text-white" /> : <Users className="w-7 h-7 text-white" />}
          </div>
          <div className="flex-1">
            <div className="flex items-center gap-2">
              <h1 className="text-2xl font-bold text-white">{group.name}</h1>
              {isOwner && <span className="text-xs bg-gold-600/20 text-gold-400 px-2 py-0.5 rounded">组长</span>}
              {isAdmin && !isOwner && <span className="text-xs bg-blue-600/20 text-blue-400 px-2 py-0.5 rounded">管理员</span>}
            </div>
            <p className="text-slate-400 text-sm">邀请码：<span className="font-mono text-white">{group.invite_code}</span></p>
          </div>
        </div>

        <div className="flex gap-3">
          <Link
            to={`/groups/${id}/members`}
            className="flex-1 bg-slate-700 hover:bg-slate-600 text-white text-center py-2.5 rounded-lg transition-colors text-sm"
          >
            查看成员
          </Link>
          {isOwner && (
            <button
              onClick={handleDelete}
              className="flex items-center gap-1.5 bg-red-600/20 hover:bg-red-600 text-red-400 px-4 py-2.5 rounded-lg transition-colors text-sm"
            >
              <Trash2 className="w-4 h-4" />
              解散组
            </button>
          )}
          {!isOwner && (
            <button
              onClick={handleLeave}
              className="flex items-center gap-1.5 bg-red-600/20 hover:bg-red-600 text-red-400 px-4 py-2.5 rounded-lg transition-colors text-sm"
            >
              <LogOut className="w-4 h-4" />
              离开组
            </button>
          )}
        </div>
      </div>

      {/* Competitions Section */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <Trophy className="w-5 h-5 text-slate-400" />
            <h2 className="text-lg font-semibold text-white">追踪赛事</h2>
          </div>
          {isAdmin && (
            <Link
              to={`/groups/${id}/add-competition`}
              className="flex items-center gap-1 text-blue-400 hover:text-blue-300 text-sm transition-colors"
            >
              <Plus className="w-4 h-4" />
              添加赛事
            </Link>
          )}
        </div>

        {competitions.length === 0 ? (
          <div className="bg-slate-800 border border-slate-700 rounded-xl p-8 text-center">
            <Trophy className="w-12 h-12 text-slate-600 mx-auto mb-3" />
            <p className="text-slate-400 mb-1">暂未追踪任何赛事</p>
            {isAdmin && (
              <Link to={`/groups/${id}/add-competition`} className="text-blue-400 hover:underline text-sm">
                添加一个赛事
              </Link>
            )}
          </div>
        ) : (
          <div className="space-y-3">
            {competitions.map((comp) => (
              <div key={comp.id} className="bg-slate-800 border border-slate-700 rounded-xl p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-medium text-white">{comp.name}</h3>
                  </div>
                  <Link
                    to={`/groups/${id}/leaderboard/${comp.id}`}
                    className="bg-pitch-600 hover:bg-pitch-700 text-white px-4 py-1.5 rounded-lg text-sm transition-colors"
                  >
                    组排行榜
                  </Link>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
