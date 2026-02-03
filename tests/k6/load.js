import http from "k6/http";
import { check, sleep } from "k6";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

export const options = {
  vus: 50,
  duration: "30s",
};

export function setup() {
  const email = `loadtest-${uuidv4()}@example.com`;
  const password = "password123";

  // Register
  const signupRes = http.post(
    "http://localhost:8080/signup",
    JSON.stringify({
      email: email,
      password: password,
    }),
    {
      headers: { "Content-Type": "application/json" },
    },
  );

  check(signupRes, {
    "signup status is 201": (r) => r.status === 201,
  });

  // Login
  const loginRes = http.post(
    "http://localhost:8080/login",
    JSON.stringify({
      email: email,
      password: password,
    }),
    {
      headers: { "Content-Type": "application/json" },
    },
  );

  check(loginRes, {
    "login status is 200": (r) => r.status === 200,
  });

  const token = loginRes.json("token");
  return { token: token };
}

export default function (data) {
  const params = {
    headers: {
      Authorization: `Bearer ${data.token}`,
      "Content-Type": "application/json",
    },
  };

  const res = http.get("http://localhost:8080/favorites", params);
  check(res, {
    "status is 200": (r) => r.status === 200,
    "protocol is HTTP/1.1": (r) => r.proto === "HTTP/1.1",
  });
  sleep(1);
}
