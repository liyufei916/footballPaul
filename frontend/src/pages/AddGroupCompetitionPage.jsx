import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { getCompetitions, getGroupCompetitions, addGroupCompetition } from '../api/apiClient';
import { ArrowLeft, Trophy, Check } from 'lucide-react';

export default function AddGroupCompetitionPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [allCompetitions, setAllCompetitions] = useState([]);
  const [trackedIds, setTrackedIds] = useState(new Set());
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [selected, setSelected] = useState(new Set());

  useEffect(() => {
    fetchData();
  }, [id]);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const [compsRes, trackedRes] = await Promise.all([
        getCompetitions(),
        getGroupCompetitions(id),
      ]);
      const all = compsRes.data.competitions || [];
      const tracked = trackedRes.data.competitions || [];
      setAllCompetitions(all);
      setTrackedIds(new Set(tracked.map((c) => c.id)));
    } catch (err) {
      console.error('Failed to fetch competitions:', err);
    } finally {
      setIsLoading(false);
    }
  };

  const toggleCompetition = (compId) => {
    const newSelected = new Set(selected);
    if (newSelected.has(compId)) {
      newSelected.delete(compId);
    } else {
      newSelected.add(compId);
    }
    setSelected(newSelected);
  };

  const handleAdd = async () => {
    if (selected.size === 0) return;
    setIsSaving(true);
    try {
      for (const compId of selected) {
        await addGroupCompetition(id, compId);
      }
      navigate(`/groups/${id}`);
    } catch (err) {
      alert(err.response?.data?.error || '添加失败');
    } finally {
      setIsSaving(false);
    }
  };

  const availableCompetitions = allCompetitions.filter((c) => !trackedIds.has(c.id));

  return (
    <div className="max-w-md mx-auto px-4 py-8">
      <div className="flex items-center gap-3 mb-8">
        <Link to={`/groups/${id}`} className="text-slate-400 hover:text-white transition-colors">
          <ArrowLeft className="w-6 h-6" />
        </Link>
        <h1 className="text-xl font-bold text-white">添加追踪赛事</h1>
      </div>

      {isLoading ? (
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-slate-800 rounded-xl p-4 animate-pulse">
              <div className="h-5 bg-slate-700 rounded w-2/3"></div>
            </div>
          ))}
        </div>
      ) : availableCompetitions.length === 0 ? (
        <div className="text-center py-16">
          <Trophy className="w-16 h-16 text-slate-600 mx-auto mb-4" />
          <p className="text-slate-400">所有赛事都已在追踪列表中</p>
        </div>
      ) : (
        <>
          <div className="space-y-3 mb-6">
            {availableCompetitions.map((comp) => (
              <button
                key={comp.id}
                onClick={() => toggleCompetition(comp.id)}
                className={`w-full flex items-center gap-3 bg-slate-800 border rounded-xl p-4 transition-colors ${
                  selected.has(comp.id)
                    ? 'border-blue-500 bg-blue-600/10'
                    : 'border-slate-700 hover:border-slate-600'
                }`}
              >
                <div className={`w-6 h-6 rounded-full border-2 flex items-center justify-center ${
                  selected.has(comp.id)
                    ? 'border-blue-500 bg-blue-500'
                    : 'border-slate-600'
                }`}>
                  {selected.has(comp.id) && <Check className="w-4 h-4 text-white" />}
                </div>
                <div className="flex-1 text-left">
                  <span className="font-medium text-white">{comp.name}</span>
                </div>
              </button>
            ))}
          </div>

          <button
            onClick={handleAdd}
            disabled={selected.size === 0 || isSaving}
            className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-slate-700 disabled:text-slate-500 text-white font-medium py-3 rounded-lg transition-colors"
          >
            {isSaving ? '添加中...' : `添加 ${selected.size > 0 ? `(${selected.size})` : ''}`}
          </button>
        </>
      )}
    </div>
  );
}
