// FROM USE - "k6 run test_k6.js"

import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '5s', target: 50 },  // Постепенно увеличиваем нагрузку до 50 пользователей
    { duration: '5s', target: 100 },  // Держим 100 пользователей
    { duration: '5s', target: 0 },   // Плавно завершаем тест
  ],
};

const BASE_URL = 'http://localhost:8080';

export default function () {
  // 1. POST /users — создание пользователя
  const createUserPayload = JSON.stringify({
    name: "Alexey",
    age: "20",
    is_student: false,
  });

  const createUserRes = http.post(`${BASE_URL}/users`, createUserPayload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(createUserRes, {
    'POST /users status is 201': (r) => r.status === 201,
  });

  // 2. GET /users — получение списка пользователей
  const listUsersRes = http.get(`${BASE_URL}/users`);
  check(listUsersRes, {
    'GET /users status is 200': (r) => r.status === 200,
  });

  // 3. GET /users/:id — получение пользователя по случайному ID (1-500)
  const randomId = Math.floor(Math.random() * 500) + 1;
  const getUserRes = http.get(`${BASE_URL}/users/${randomId}`);
  check(getUserRes, {
    'GET /users/:id status is 200': (r) => r.status === 200,
  });

  // 4. PUT /users/:id — обновление возраста пользователя
  const updateUserPayload = JSON.stringify({
    age: 21,
  });

  const updateUserRes = http.put(`${BASE_URL}/users/${randomId}`, updateUserPayload, {
    headers: { 'Content-Type': 'application/json' },
  });

  check(updateUserRes, {
    'PUT /users/:id status is 200': (r) => r.status === 200,
  });


// 5. DELETE

const randomIdDel = Math.floor(Math.random() * 500) + 1;
const getUserResDel = http.del(`${BASE_URL}/users/${randomId}`);
check(getUserResDel, {
  'DELETE /users/:id status is 200': (r) => r.status === 200,
});


  sleep(1); // Задержка между запросами (1 секунда)
}