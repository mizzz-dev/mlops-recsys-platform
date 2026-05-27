import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<300'],
  },
  scenarios: {
    smoke: {
      executor: 'constant-vus',
      vus: 5,
      duration: '30s',
    },
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const ID_TOKEN = __ENV.ID_TOKEN || '';

function requestParams() {
  if (!ID_TOKEN) {
    return {};
  }
  return {
    headers: {
      Authorization: `Bearer ${ID_TOKEN}`,
    },
  };
}

export default function () {
  const userId = `user_${String(Math.floor(Math.random() * 100) + 1).padStart(3, '0')}`;
  const response = http.get(`${BASE_URL}/v1/recommendations?user_id=${userId}&limit=5`, requestParams());

  check(response, {
    'status is 200': (res) => res.status === 200,
    'strategy is returned': (res) => Boolean(res.json('strategy')),
    'recommendations are returned': (res) => Array.isArray(res.json('recommendations')),
  });

  sleep(1);
}
