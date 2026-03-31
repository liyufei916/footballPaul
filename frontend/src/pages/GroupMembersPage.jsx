import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getGroupMembers, getGroup } from '../api/apiClient';
import { ArrowLeft, Crown, User } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

export default function GroupMembersPage() {
  const { id } = useParams();
  const { user } = useAuth();
  const [members, setMembers] = useState([]);
  const [group, setGroup] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, [id]);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const [membersRes, groupRes] = await Promise.all([
        getGroupMembers(id),
        getGroup(id),
      ]);
      setMembers(membersRes.data.members || []);
      setGroup(groupRes.data.group);
    } catch (err) {
      console.error('Failed to fetch members:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const isOwner = group?.owner_id === user?.id;

  return (
    <div className="max-w-md mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-8">
        <Link to={`/groups/${id}`} className="text-slate-400 hover:text-white transition-colors">
          <ArrowLeft className="w-6 h-6" />
        </Link>
        <h1 className="text-xl font-bold text-white">组成员</h1>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-slate-800 rounded-xl p-4 animate-pulse">
              <div className="h-5 bg-slate-700 rounded w-1/3"></div>
            </div>
          ))}
        </div>
      ) : (
        <div className="space-y-3">
          {members.map((member) => (
            <div
              key={member.user_id}
              className="flex items-center gap-3 bg-slate-800 border border-slate-700 rounded-xl p-4"
            >
              <div className={`w-10 h-10 rounded-full flex items-center justify-center ${
                member.role === 'admin' ? 'bg-gold-600' : 'bg-slate-600'
              }`}>
                {member.role === 'admin' ? (
                  <Crown className="w-5 h-5 text-white" />
                ) : (
                  <User className="w-5 h-5 text-white" />
                )}
              </div>
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <span className="font-medium text-white">{member.username}</span>
                  {member.role === 'admin' && (
                    <span className="text-xs bg-gold-600/20 text-gold-400 px-2 py-0.5 rounded">
                      {member.user_id === group?.owner_id ? '组长' : '管理员'}
                    </span>
                  )}
                </div>
                <p className="text-slate-400 text-xs">
                  加入于 {new Date(member.joined_at).toLocaleDateString('zh-CN')}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
