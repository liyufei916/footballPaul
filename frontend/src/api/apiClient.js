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
