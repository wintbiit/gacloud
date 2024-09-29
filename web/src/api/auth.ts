import { axioser } from "./index.ts";

export interface LoginResponse {
  token: string;
}

const login = async (
  username: string,
  password: string,
): Promise<LoginResponse> => {
  const resp = await axioser.post<LoginResponse>("/login", {
    username,
    password,
  });

  return resp.data;
};

export { login };
