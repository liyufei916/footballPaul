import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;

// Competition APIs
export const getCompetitions = () => api.get('/competitions');
export const getCompetition = (id) => api.get(`/competitions/${id}`);
export const createCompetition = (data) => api.post('/competitions', data);

// Match APIs
export const getMatches = (params) => api.get('/matches', { params });
export const getMatch = (id) => api.get(`/matches/${id}`);
export const createMatch = (data) => api.post('/matches', data);
export const enterResult = (id, data) => api.put(`/matches/${id}/result`, data);

// Prediction APIs
export const createPrediction = (data) => api.post('/predictions', data);
export const updatePrediction = (id, data) => api.put(`/predictions/${id}`, data);
export const getMyPredictions = (competitionId) => {
  const params = competitionId ? { competition_id: competitionId } : {};
  return api.get('/predictions/my', { params });
};
export const getMatchPredictions = (matchId) => api.get(`/matches/${matchId}/predictions`);

// Leaderboard APIs
export const getLeaderboard = (competitionId, limit = 50) => {
  const params = { limit };
  if (competitionId) {
    params.competition_id = competitionId;
  }
  return api.get('/leaderboard', { params });
};
export const getMyRank = (competitionId) => {
  const params = competitionId ? { competition_id: competitionId } : {};
  return api.get('/leaderboard/my-rank', { params });
};

// User APIs
export const getProfile = () => api.get('/users/profile');

// Group APIs
export const getMyGroups = () => api.get('/groups');
export const getGroup = (id) => api.get(`/groups/${id}`);
export const createGroup = (data) => api.post('/groups', data);
export const joinGroup = (inviteCode) => api.post('/groups/join', { invite_code: inviteCode });
export const leaveGroup = (id) => api.delete(`/groups/${id}/leave`);
export const deleteGroup = (id) => api.delete(`/groups/${id}`);
export const getGroupMembers = (id) => api.get(`/groups/${id}/members`);
export const getGroupCompetitions = (id) => api.get(`/groups/${id}/competitions`);
export const addGroupCompetition = (id, competitionId) =>
  api.post(`/groups/${id}/competitions`, { competition_id: competitionId });
export const removeGroupCompetition = (id, competitionId) =>
  api.delete(`/groups/${id}/competitions/${competitionId}`);
export const getGroupLeaderboard = (groupId, competitionId, limit = 50) =>
  api.get(`/groups/${groupId}/leaderboard/${competitionId}`, { params: { limit } });
export const transferGroupOwnership = (id, newOwnerId) =>
  api.put(`/groups/${id}/transfer-owner`, { new_owner_id: newOwnerId });
export const getGroupMatchPredictions = (groupId, competitionId) =>
  api.get(`/groups/${groupId}/competitions/${competitionId}/predictions`);
