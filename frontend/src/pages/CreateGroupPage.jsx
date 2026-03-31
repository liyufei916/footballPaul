import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createGroup } from '../api/apiClient';
import { ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function CreateGroupPage() {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!name.trim()) {
      setError('请输入组名');
      return;
    }
    if (name.trim().length > 50) {
      setError('组名最长50个字符');
      return;
    }

    setIsLoading(true);
    setError('');
    try {
      const res = await createGroup({ name: name.trim() });
      navigate(`/groups/${res.data.group.id}`);
    } catch (err) {
      setError(err.response?.data?.error || '创建失败，请重试');
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
        <h1 className="text-xl font-bold text-white">创建组</h1>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label className="block text-slate-300 text-sm font-medium mb-2">组名</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="例如：欧冠竞猜群"
            maxLength={50}
            className="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-3 text-white placeholder-slate-500 focus:outline-none focus:border-blue-500 transition-colors"
            autoFocus
          />
          {error && <p className="mt-2 text-red-400 text-sm">{error}</p>}
        </div>

        <button
          type="submit"
          disabled={isLoading}
          className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-slate-700 text-white font-medium py-3 rounded-lg transition-colors"
        >
          {isLoading ? '创建中...' : '创建组'}
        </button>
      </form>
    </div>
  );
}
