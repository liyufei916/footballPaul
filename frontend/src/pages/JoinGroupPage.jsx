import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { joinGroup } from '../api/apiClient';
import { ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function JoinGroupPage() {
  const navigate = useNavigate();
  const [inviteCode, setInviteCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!inviteCode.trim() || inviteCode.trim().length !== 6) {
      setError('请输入6位邀请码');
      return;
    }

    setIsLoading(true);
    setError('');
    try {
      const res = await joinGroup(inviteCode.trim().toUpperCase());
      navigate(`/groups/${res.data.group.id}`);
    } catch (err) {
      const msg = err.response?.data?.error;
      if (msg === '邀请码无效') {
        setError('邀请码无效，请检查后重新输入');
      } else if (msg === '你已经在该组中') {
        setError('你已经在该组中');
      } else {
        setError('加入失败，请重试');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-md mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-8">
        <Link to="/groups" className="text-slate-400 hover:text-white transition-colors">
          <ArrowLeft className="w-6 h-6" />
        </Link>
        <h1 className="text-xl font-bold text-white">加入组</h1>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label className="block text-slate-300 text-sm font-medium mb-2">邀请码</label>
          <input
            type="text"
            value={inviteCode}
            onChange={(e) => setInviteCode(e.target.value.toUpperCase().replace(/[^A-Z0-9]/g, ''))}
            placeholder="6位字母或数字"
            maxLength={6}
            className="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-3 text-white text-center text-2xl tracking-widest placeholder-slate-600 focus:outline-none focus:border-blue-500 transition-colors font-mono"
            autoFocus
          />
          {error && <p className="mt-2 text-red-400 text-sm text-center">{error}</p>}
        </div>

        <button
          type="submit"
          disabled={isLoading || inviteCode.length !== 6}
          className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-slate-700 disabled:text-slate-500 text-white font-medium py-3 rounded-lg transition-colors"
        >
          {isLoading ? '加入中...' : '加入组'}
        </button>
      </form>
    </div>
  );
}
